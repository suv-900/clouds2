package routers

import (
	"github.com/gorilla/mux"
	"github.com/suv-900/blog/controllers"
)

func HandleRoutes(router *mux.Router) {
	router.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/checkusername", controllers.SearchUsername).Methods("POST")
	//	router.HandleFunc("/profile", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/serverstatus", controllers.CheckServerHealth).Methods("GET")
	router.HandleFunc("/createpost", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/authtoken", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/deleteuser", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/addcomment/{id:[0-9]+}", controllers.AddComment).Methods("POST")
	router.HandleFunc("/likecomment/{id:[0-9]+}", controllers.LikeComment).Methods("POST")
	router.HandleFunc("/dislikecomment/{id:[0-9]+}", controllers.DislikeComment).Methods("POST")

	router.HandleFunc("/posts/getposts/{offset:[0-9]+}", controllers.GetPostsMetaData).Methods("GET")

	router.HandleFunc("/posts/getall", controllers.GetAllPostsMetaData).Methods("GET")
	router.HandleFunc("/posts/get-featured-posts/{offset:[0-9]+}", controllers.GetFeaturedPosts).Methods("GET")

	router.HandleFunc("/viewpost/{id:[0-9]+}", controllers.GetPostByID).Methods("GET")
	router.HandleFunc("/viewPostToken/{id:[0-9]+}", controllers.GetPostByID_WithUserPreferences).Methods("GET")

	router.HandleFunc("/likepost/{postid:[0-9]+}", controllers.LikePost).Methods("POST")
	router.HandleFunc("/removelike/{postid:[0-9]+}", controllers.RemoveLikeFromPost).Methods("POST")

	router.HandleFunc("/dislikepost/{postid:[0-9]+}", controllers.DislikePost).Methods("POST")
	router.HandleFunc("/removedislike/{postid:[0-9]+}", controllers.RemoveDislikeFromPost).Methods("POST")
}
