package main

import (
	"log"

	"github.com/barelyhuman/tasks/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverRef := server.NewServer()
	serverRef.Router.GET("/", serverRef.CreateHandler(HomeHandler))
	serverRef.Router.GET("/edit", serverRef.CreateHandler(EditHandler))
	serverRef.Router.POST("/edit", serverRef.CreateHandler(UpdateTasksHandler))

	serverRef.StartServer()
}
