package api

import (
	"fmt"
	"net/http"

	"cloud/internal/data"
	"cloud/internal/logger"
	"cloud/models"
	"cloud/routers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var log = logger.GetLogger()

type application struct {
	models *data.Models
}

// implement cors
func main() {

	router := mux.NewRouter()
	err := models.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:5173"},
	})

	handler := c.Handler(router)
	routers.HandleRoutes(router)
	fmt.Println("Server started at port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
