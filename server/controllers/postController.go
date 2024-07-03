package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	//"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/suv-900/blog/models"
)

//TODO no jwt
/*
->CRUD
->getbyid
*/
func GetallpostsbyUser(w http.ResponseWriter, r *http.Request) {
	var userId uint64
	ok, p := TokenVerifier("userToken", r)
	if ok {
		userId = p.ID
	} else {
		w.WriteHeader(401)
		return
	}

	posts := models.GetPostsByUserId(userId)

	parsedRes, err := json.Marshal(posts)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
	w.Write(parsedRes)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {

	tokenExpired, authorid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired || tokenInvalid {
		w.WriteHeader(401)
		return
	}

	rbyte, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(&w, err)
		return
	}
	var post models.Posts
	err = json.Unmarshal(rbyte, &post)
	if err != nil {
		serverError(&w, err)
		return
	}
	if post.Post_content == "" || post.Post_title == "" {
		w.WriteHeader(400)
		return
	}
	author_name := models.GetUsername(authorid)
	post.Author_id = authorid
	post.Author_name = author_name

	fmt.Println(post)

	postid, err := models.CreatePost(post)
	if err != nil {
		serverError(&w, err)
		return
	}

	parsedres, err := json.Marshal(postid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(parsedres)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	var postid uint64
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
	pipe2 := make(chan bool, 1)
	go func() {
		err := models.DeletePost(postid)
		if err != nil {
			serverError(&w, err)
			pipe2 <- false
			return
		}
		pipe2 <- true
	}()
	if !<-pipe2 {
		fmt.Println("error while deleting post")
		return
	}
	fmt.Println("post deleted succesfully")
	w.WriteHeader(200)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var postid uint64
	var tokenExpired bool
	var tokenInvalid bool
	a := make(chan int, 1)
	go func() {
		tokenExpired, _, tokenInvalid = AuthenticateTokenAndSendUserID(r)
		a <- 1
	}()
	<-a

	if tokenExpired {
		w.WriteHeader(401)
		return
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	pipe2 := make(chan bool, 1)
	go func() {

		rbyte, err := io.ReadAll(r.Body)
		if err != nil {
			pipe2 <- false
			serverError(&w, err)
			return
		}

		var post models.Posts
		err = json.Unmarshal(rbyte, &post)
		if err != nil {
			pipe2 <- false
			serverError(&w, err)
			return
		}

		err = models.UpdatePost(postid, post)
		if err != nil {
			pipe2 <- false
			serverError(&w, err)
			return
		}
		pipe2 <- true

	}()
	if !<-pipe2 {
		fmt.Println("error while updating post")
		return
	}
	w.WriteHeader(200)
}

func GetPostsMetaData(w http.ResponseWriter, r *http.Request) {
	offsetString := r.URL.Query().Get("offset")
	limitString := r.URL.Query().Get("limit")

	offset, err := strconv.ParseUint(offsetString, 10, 16)
	if err != nil {
		serverError(&w, err)
		return
	}
	limit, err := strconv.ParseUint(limitString, 10, 16)
	if err != nil {
		serverError(&w, err)
		return
	}
	posts := models.GetPostsMetaData(offset, limit)

	response, err := json.Marshal(posts)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

func GetAllPostsMetaData(w http.ResponseWriter, r *http.Request) {

	posts := models.GetAllPostsMetaData()

	response, err := json.Marshal(posts)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(response)

}
func GetFeaturedPosts(w http.ResponseWriter, r *http.Request) {
	var offset uint64
	offsetString := r.URL.Query().Get("offset")
	offset, err := strconv.ParseUint(offsetString, 10, 16)
	if err != nil {
		serverError(&w, err)
		return
	}

	postMetaData := models.GetFeaturedPosts(offset)
	response, err := json.Marshal(postMetaData)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

// sends post username and top 5 comments
func GetPostByID(w http.ResponseWriter, r *http.Request) {
	postidString := r.URL.Query().Get("postid")
	postid, err := strconv.ParseUint(postidString, 10, 16)
	if err != nil {
		serverError(&w, err)
		return
	}

	post, err := models.PostById(postid)
	if err != nil {
		serverError(&w, err)
		return
	}
	year, month, day := post.Createdat.Date()
	post.Createdat_str = fmt.Sprintf("%d-%d-%d", day, month, year)

	comments, err := models.GetComments(postid, 5, 0)
	if err != nil {
		serverError(&w, err)
		return
	}
	for i := 0; i < len(comments); i++ {
		comments[i].Createdat_str = comments[i].Createdat.Local().Format(time.RFC822)
	}
	finalRes := models.PostandComments{Post: post, Comments: comments}
	parsedRes, err := json.Marshal(finalRes)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
	w.Write(parsedRes)

}
func GetPostByID_WithUserPreferences(w http.ResponseWriter, r *http.Request) {

	var userid uint64
	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(401)
		return
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	postidString := r.URL.Query().Get("postid")
	postid, err := strconv.ParseUint(postidString, 10, 16)
	if err != nil {
		serverError(&w, err)
		return
	}

	post, err := models.PostById(postid)
	if err != nil {
		serverError(&w, err)
		return
	}
	year, month, day := post.Createdat.Date()
	post.Createdat_str = fmt.Sprintf("%d-%d-%d", day, month, year)

	userLikedPost, userDislikedPost, err := models.CheckUserReaction(userid, postid)
	if err != nil {
		serverError(&w, err)
		return
	}

	comments := models.GetUserCommentReaction(postid, userid)
	for i := 0; i < len(comments); i++ {
		comments[i].Createdat_str = comments[i].Createdat.Format(time.RFC1123)
	}
	// if err != nil {
	// 	serverError(&w, err)
	// 	return
	// }

	finalResult := &models.PostComments_WithUserPreference{
		Post:               post,
		PostLikedByUser:    userLikedPost,
		PostDislikedByUser: userDislikedPost,
		Comments:           comments,
	}

	var jsonReply []byte
	jsonReply, err = json.Marshal(finalResult)
	if err != nil {
		serverError(&w, err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonReply)

}
func LikePost(w http.ResponseWriter, r *http.Request) {
	var postid uint64

	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(400)
		return
	}
	if tokenInvalid {
		w.WriteHeader(401)
		return
	}

	vars := mux.Vars(r)
	postidstr := vars["postid"]
	postid, err := strconv.ParseUint(postidstr, 10, 64)
	if err != nil {
		serverError(&w, err)
		return
	}

	//save user prefrence
	err = models.LikePost(postid, userid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(401)
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	var postid uint64
	vars := mux.Vars(r)
	postidstr := vars["postid"]
	postid, err := strconv.ParseUint(postidstr, 10, 64)
	if err != nil {
		serverError(&w, err)
		return
	}

	//save user prefrence
	models.DislikePost(postid, userid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
}

func RemoveLikeFromPost(w http.ResponseWriter, r *http.Request) {
	var postid uint64
	var err error

	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(401)
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}
	vars := mux.Vars(r)
	postidstr := vars["postid"]
	postid, err = strconv.ParseUint(postidstr, 10, 64)
	if err != nil {
		serverError(&w, err)
		return
	}

	err = models.RemoveLikeFromPost(postid, userid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
}
func RemoveDislikeFromPost(w http.ResponseWriter, r *http.Request) {
	var postid uint64

	tokenExpired, userid, tokenInvalid := AuthenticateTokenAndSendUserID(r)
	if tokenExpired {
		w.WriteHeader(401)
	}
	if tokenInvalid {
		w.WriteHeader(400)
		return
	}

	vars := mux.Vars(r)
	postidstr := vars["postid"]
	postid, err := strconv.ParseUint(postidstr, 10, 64)
	if err != nil {
		serverError(&w, err)
		return
	}

	err = models.RemoveDislikeFromPost(postid, userid)
	if err != nil {
		serverError(&w, err)
		return
	}
	w.WriteHeader(200)
}
func TokenVerifier(s string, r *http.Request) (bool, *CustomPayload) {
	t := GetCookieByName(r.Cookies(), s)
	if t == "" {
		fmt.Println("no cookie got.")
		return false, nil
	}
	token, err := jwt.ParseWithClaims(t, &CustomPayload{}, nil)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if p, ok := token.Claims.(*CustomPayload); ok && token.Valid {
		return true, p
	} else {
		fmt.Println("Token not ok!")
		return false, nil
	}
}

func GetCookieByName(cookies []*http.Cookie, cookiename string) string {
	result := ""
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == cookiename {
			result += cookies[i].Value
			break
		}
	}
	return result
}
