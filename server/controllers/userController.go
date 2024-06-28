package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	// "strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	// "github.com/gorilla/mux"
	"github.com/suv-900/blog/models"
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

// returns 200
func CheckServerHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
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
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userid uint64
	var tokenExpired bool
	var tokenInvalid bool
	a := make(chan int, 1)
	go func() {
		tokenExpired, userid, tokenInvalid = AuthenticateTokenAndSendUserID(r)
		a <- 1

	}()

	if tokenExpired {
		w.WriteHeader(401)
		return
	}

	if tokenInvalid {
		w.WriteHeader(400)
		return
	}
	/*
		var userExists bool
		c := make(chan int, 1)
		go func() {
			userExists = models.CheckUserExists(userid)
			c <- 1
		}()
		<-c
		if !userExists {
			w.WriteHeader(400)
			return
		}
	*/
	var err error
	err = nil
	b := make(chan int, 1)
	go func() {

		err = models.DeleteUser(userid)
		if err != nil {
			b <- 1
			return
		}
		//err = models.DeleteUser(userid)
		/*
			if err != nil {
				b <- 1
				return
			}*/
		b <- 1
	}()
	<-b
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
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

// func GetUserInfo(w http.ResponseWriter, r *http.Request) {
// 	var userid uint64
// 	var username string
// 	var err error

// 	vars := mux.Vars(r)
// 	useridstr := vars["userid"]
// 	userid, err = strconv.ParseUint(useridstr, 10, 64)

// 	var userInfo UserInfo

// }

/*
	func AddProfilePicture(w http.ResponseWriter, r *http.Request) {
		//TODO add the pic to a folder and add the address to the DB
		var err error
		var tokenExpired bool
		var tokenInvalid bool
		var userid uint64
		a := make(chan int, 1)
		go func() {
			tokenExpired, userid, tokenInvalid = AuthenticateTokenAndSendUserID(r)
			a <- 1

		}()
		<-a
		if tokenInvalid {
			w.WriteHeader(401)
			return
		}
		if tokenExpired {
			w.WriteHeader(400)
			return
		}

		var imageString string
		rbody, err := io.ReadAll(r.Body)
		if err != nil {
			serverError(&w, err)
			return
		}
		err = json.Unmarshal(rbody, &imageString)
		if err != nil {
			serverError(&w, err)
			return
		}

}
*/
func Getass(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(400)
		return
	}

	var wg sync.WaitGroup

	//1

	channel1 := make(chan models.Users)
	wg.Add(1)
	go func() {
		defer wg.Done()
		channel1 <- models.GetUserDetails(username)
		//userDetails.Username = username
	}()
	userDetails := <-channel1
	fmt.Println(userDetails)

	//2
	channel2 := make(chan []models.Posts)
	wg.Add(1)
	go func() {
		defer wg.Done()
		channel2 <- models.GetPostsByUserId(userDetails.UserID)
	}()
	posts := <-channel2

	wg.Wait()

	userPost := models.UserAndPost{User: userDetails, Posts: posts}
	parsedRes, err := json.Marshal(userPost)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
	w.Write(parsedRes)

}

// func parseReply(data any)
func serverError(w *http.ResponseWriter, err error) {
	if err != nil {
		fmt.Println(err)
	}
	(*w).WriteHeader(500)
}

/*
func GenerateSessionId() int64 {

}
*/

func CreateToken(w http.ResponseWriter, r *http.Request) {
	p := CustomPayload{
		1,
		jwt.StandardClaims{
			ExpiresAt: Tokenexpirytime.Unix(),
			Issuer:    "createUser handler",
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	token, err := rawToken.SignedString(JWTKEY)
	if err != nil {
		serverError(&w, err)
		return
	}
	res, _ := json.Marshal(token)
	w.Write(res)
}
