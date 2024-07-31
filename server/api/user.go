package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"sync"
	"time"

	"cloud/models"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

// TODO token blacklist
var bycryptCost = 3
var JWTKEY = []byte(os.Getenv("JWT_KEY"))
var Tokenexpirytime = time.Now().Add(60 * time.Minute)

type CustomPayload struct {
	ID uint64 `json:"id"`
	jwt.StandardClaims
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (app *application) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)

	username := m["username"]
	if len(username) == 0 {
		app.badRequest(w, r)
		return
	}
}

func SearchUsername(w http.ResponseWriter, r *http.Request) {
	rbyte, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}

	var username string

	err = json.Unmarshal(rbyte, &username)
	if err != nil {
		serverError(&w, err)
		return
	}

	userFound := models.FindUser(username)

	if userFound {
		w.WriteHeader(409)
	} else {
		w.WriteHeader(200)
	}

}

// creates a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var err error

	rbyte, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}
	var user models.Users
	err = json.Unmarshal(rbyte, &user)
	if err != nil {
		serverError(&w, err)
		return
	}

	userFound := models.FindUser(user.Username)
	if userFound {
		w.WriteHeader(409)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bycryptCost)
	if err != nil {
		serverError(&w, err)
		return
	}

	user.Password = string(hashedPassword)

	var userid uint64
	userid, err = models.CreateUser(user)
	if err != nil {
		serverError(&w, err)
		return
	}

	//TODO add user with same name
	//TODO i dont think this call gets awaited/the goroutine waits for this call

	//var err error

	var token string
	payload := CustomPayload{
		ID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: Tokenexpirytime.Unix(),
			Issuer:    "createUser handler",
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err = rawToken.SignedString(JWTKEY)
	if err != nil {
		serverError(&w, err)
		return
	}

	t, err := json.Marshal(token)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(t)
}

// completed
func LoginUser(w http.ResponseWriter, r *http.Request) {
	rbytes, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
	}
	var user models.Users
	json.Unmarshal(rbytes, &user)

	dbpassword, exists, id := models.LoginUser(user.Username)
	if err != nil {
		serverError(&w, err)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	passValid := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(user.Password))
	if passValid != nil {
		w.WriteHeader(401)
		return
	}

	payload := CustomPayload{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: Tokenexpirytime.Unix(),
			Issuer:    "loginHandler",
		},
	}
	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := Token.SignedString(JWTKEY)
	if err != nil {
		serverError(&w, err)
		return
	}
	ts, err := json.Marshal(t)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
	w.Write(ts)
}
func UpdateUserAbout(w http.ResponseWriter, r *http.Request) {
	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired || tokenInvalid {
		w.WriteHeader(401)
		return
	}

	rbyte, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}
	var userAbout string
	err = json.Unmarshal(rbyte, &userAbout)
	if err != nil {
		serverError(&w, err)
		return
	}

	err = models.UpdateUserAbout(userAbout, userid)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired || tokenInvalid {
		w.WriteHeader(401)
		return
	}

	err := models.DeleteUser(userid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
}

func UpdateUserPass(w http.ResponseWriter, r *http.Request) {
	rbyte, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}

	var newPass string
	err = json.Unmarshal(rbyte, &newPass)
	if err != nil {
		serverError(&w, err)
		return
	}

	var workers sync.WaitGroup

	//1
	channel1 := make(chan uint64, 1)
	errorChannel := make(chan error)
	workers.Add(1)
	go func() {
		defer workers.Done()
		tokenstr := GetCookieByName(r.Cookies(), "userToken")
		token, err := jwt.ParseWithClaims(tokenstr, &CustomPayload{}, nil)
		if err != nil {
			errorChannel <- err
			return
		}
		if payload, ok := token.Claims.(*CustomPayload); ok && token.Valid {
			channel1 <- payload.ID
		}
	}()
	if <-errorChannel != nil {
		//fmt.Println(<-errorChannel)
		//w.WriteHeader(500)
		//return
		serverError(&w, <-errorChannel)
		return
	}
	userid := <-channel1

	//2

	channel2 := make(chan []byte)
	workers.Add(1)
	go func() {
		defer workers.Done()
		pass, err := bcrypt.GenerateFromPassword([]byte(newPass), bycryptCost)
		if err != nil {
			errorChannel <- err
			return
		}
		channel2 <- pass
	}()
	if <-errorChannel != nil {
		serverError(&w, <-errorChannel)
		return
	}
	hashpass := <-channel2

	//3

	workers.Add(1)
	go func() {
		defer workers.Done()
		errorChannel <- models.UpdatePass(string(hashpass), userid)
	}()
	if <-errorChannel != nil {
		serverError(&w, <-errorChannel)
		return
	}

	//wait
	workers.Wait()

	w.WriteHeader(200)

}
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(400)
		return
	}

	u, err := models.GetUserInfo(username)
	if u.UserID == -1 {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		serverError(&w, err)
		return
	}

	user := models.UserInfo{
		UserID:   u.UserID,
		Username: u.Username,
		About:    u.About,
		JoinDate: u.Createdat.Local().Format(time.RFC822),
	}

	response, err := json.Marshal(user)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}
