package routers

import (
	"cloud/controllers"

	"github.com/gorilla/mux"
)

func HandleRoutes(router *mux.Router) {
	router.HandleFunc("/users/register", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/deleteuser", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/users/checkusername", controllers.SearchUsername).Methods("POST")
	router.HandleFunc("/users/update-about", controllers.UpdateUserAbout).Methods("PUT")

	router.HandleFunc("/user", controllers.GetUserInfo).Methods("GET")

	//	router.HandleFunc("/profile", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/serverstatus", controllers.CheckServerHealth).Methods("GET")
	router.HandleFunc("/createpost", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/authtoken", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/deleteuser", controllers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/addcomment/{id:[0-9]+}", controllers.AddComment).Methods("POST")
	router.HandleFunc("/likecomment/{id:[0-9]+}", controllers.LikeComment).Methods("POST")
	router.HandleFunc("/dislikecomment/{id:[0-9]+}", controllers.DislikeComment).Methods("POST")
	router.HandleFunc("/comments/getcomments", controllers.GetComments).Methods("GET")
	// router.HandleFunc("/getcommentstoken", controllers.GetCommentsWithToken).Methods("POST")

	router.HandleFunc("/posts/getposts", controllers.GetPostsMetaData).Methods("GET")
	router.HandleFunc("/posts/getall", controllers.GetAllPostsMetaData).Methods("GET")
	router.HandleFunc("/posts/get-featured-posts", controllers.GetFeaturedPosts).Methods("GET")
	router.HandleFunc("/posts/get-author-posts", controllers.GetPostsByAuthorID).Methods("GET")
	router.HandleFunc("/posts/delete-posts", controllers.DeletePosts).Methods("DELETE")
	router.HandleFunc("/posts/update-post-title", controllers.UpdatePostTitle).Methods("PUT")
	router.HandleFunc("/posts/update-post-content", controllers.UpdatePostContent).Methods("PUT")

	router.HandleFunc("/viewpost", controllers.GetPostByID).Methods("GET")
	router.HandleFunc("/viewpost-token", controllers.GetPostByID_WithUserPreferences).Methods("GET")

	router.HandleFunc("/likepost/{postid:[0-9]+}", controllers.LikePost).Methods("POST")
	router.HandleFunc("/removelike/{postid:[0-9]+}", controllers.RemoveLikeFromPost).Methods("POST")

	router.HandleFunc("/dislikepost/{postid:[0-9]+}", controllers.DislikePost).Methods("POST")
	router.HandleFunc("/removedislike/{postid:[0-9]+}", controllers.RemoveDislikeFromPost).Methods("POST")

}
