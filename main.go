package main

import (
	"capstone-backend/DAO"
	"capstone-backend/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// Register HTTP Handler
func main() {
	initialize()
	registerHandler()
	startServer()
}

func initialize() {
	DAO.Init()
}

func registerHandler() {
	r := gin.Default()
	// Allow all cross site access
	r.Use(cors.Default())

	// register all handlers
	r.GET("/", handler.IndexHandler)
	r.GET("/course/list", handler.GetCourseListHandler)
	r.GET("/course/info", handler.GetCourseInfoHandler)
	r.POST("/course/info/add", handler.AddCourseInfoHandler)
	r.GET("/course/professor/info", handler.GetCourseProfessorInfoHandler)

	// Retrieve port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// start server on port
	log.Printf("Listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
