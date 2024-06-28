package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/suv-900/blog/models"
)

func AddComment(w http.ResponseWriter, r *http.Request) {

	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(401)
		return
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	//login,register
	//home,viewpost,comments,likepost,likecomment etc
	//addcomment
	//createpost
	username := models.GetUsername(userid)

	var postid uint64
	var comment_content string
	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}
	json.Unmarshal(rbody, &comment_content)

	vars := mux.Vars(r)
	postidstr := vars["id"]
	postid, err = strconv.ParseUint(postidstr, 10, 64)
	if err != nil {
		serverError(&w, err)
		return
	}

	commentID, err := models.AddComment(postid, userid, username, comment_content)
	if err != nil {
		serverError(&w, err)
		return
	}
	comment := models.Comment{
		Comment_id:      commentID,
		User_id:         userid,
		Comment_content: comment_content,
		Username:        username,
		CreatedAt:       time.Now(),
		Comment_likes:   0,
	}
	reply, err := json.Marshal(comment)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(reply)
}

func FetchComments(w http.ResponseWriter, r *http.Request) {
	//session token
	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}
	var offset uint16
	json.Unmarshal(rbody, &offset)

	var wg sync.WaitGroup

	//1
	channel1 := make(chan uint64)
	errChannel := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ok, p := TokenVerifier("postToken", r)
		if !ok {
			fmt.Println("Error while parsing token!")
			errChannel <- true
			return
		}
		channel1 <- p.ID
		errChannel <- false
	}()

	if <-errChannel {
		w.WriteHeader(401)
		return
	}
	postid := <-channel1

	//2
	channel2 := make(chan []byte)
	err2 := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		comments := models.GetAllCommentsByPostID(postid)
		parsedRes, err := json.Marshal(comments)
		if err != nil {
			fmt.Println(err)
			err2 <- true
			return
		}
		channel2 <- parsedRes
		err2 <- false
	}()
	if <-err2 {
		w.WriteHeader(500)
		return
	}

	wg.Wait()

	w.WriteHeader(200)
	w.Write(<-channel2)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	var commentId uint64
	ok, _ := TokenVerifier("userToken", r)
	if ok {
		commentId, ok = ParseToken(GetCookieByName(r.Cookies(), "commentToken"))
		if !ok {
			w.WriteHeader(500)
			fmt.Println("Token not ok OR not Valid")
			return
		}

	} else {
		w.WriteHeader(401)
		return
	}

	rbody, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}

	var comment string
	json.Unmarshal(rbody, &comment)

	models.EditComment(commentId, comment)

	w.WriteHeader(200)

}

func RemoveComment(w http.ResponseWriter, r *http.Request) {
	var commentId uint64
	ok, _ := TokenVerifier("userToken", r)
	if ok {
		commentId, ok = ParseToken(GetCookieByName(r.Cookies(), "commentToken"))
		if !ok {
			w.WriteHeader(500)
			fmt.Println("Token not ok OR not Valid")
			return
		}

	} else {
		w.WriteHeader(401)
		return
	}

	models.RemoveComment(commentId)

	w.WriteHeader(200)

}

func ParseToken(token string) (uint64, bool) {
	t, err := jwt.ParseWithClaims(token, &CustomPayload{}, nil)
	if err != nil {
		fmt.Println("Error while parsing token.", err)
	}
	if p, ok := t.Claims.(*CustomPayload); ok && t.Valid {
		return p.ID, true
	}
	return 0, false
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	var tokenInvalid bool
	var tokenExpired bool
	a := make(chan int, 1)
	go func() {
		tokenExpired, _, tokenInvalid = AuthenticateTokenAndSendUserID(r)
		a <- 1
	}()
	<-a
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	if tokenExpired {
		w.WriteHeader(401)
		return
	}
	var commentID uint64
	var err error
	b := make(chan int, 1)
	go func() {
		rbody, err := io.ReadAll(r.Body)
		if err != nil {
			b <- 1
			return
		}
		json.Unmarshal(rbody, &commentID)
		b <- 1
	}()
	<-b
	if err != nil {
		serverError(&w, err)
		return
	}

	c := make(chan int, 1)
	go func() {
		models.LikeAComment(commentID)
		c <- 1
	}()
	<-c

	w.WriteHeader(200)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	var tokenExpired bool
	var tokenInvalid bool
	a := make(chan int, 1)
	go func() {
		tokenExpired, _, tokenInvalid = AuthenticateTokenAndSendUserID(r)
		a <- 1
	}()
	<-a

	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	if tokenExpired {
		w.WriteHeader(401)
		return
	}
	var commentID uint64
	var err error
	b := make(chan int, 1)
	go func() {
		rbody, err := io.ReadAll(r.Body)
		if err != nil {
			b <- 1
			return
		}
		json.Unmarshal(rbody, &commentID)
		b <- 1
	}()
	<-b
	if err != nil {
		serverError(&w, err)
		return
	}

	c := make(chan int, 1)
	go func() {
		models.DislikeAComment(commentID)
		c <- 1
	}()
	<-c

	w.WriteHeader(200)
}
