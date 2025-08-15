package main

import (
	"log"

	"github.com/lucas/gokafka/user-service/internal/handlers"
	"github.com/lucas/gokafka/user-service/internal/repository"
	"github.com/lucas/gokafka/user-service/internal/services"
	"github.com/lucas/gokafka/user-service/internal/session"
)

func main() {

	log.Println("Starting user-service...")

	// Initialize the user repository
	repo := repository.NewUserRepository()

	// Initialize Redis session store
	sessionStore := session.NewRedisSessionStore("localhost:6379", "", 0)

	// Initialize the user service
	service := services.NewUserService(repo, sessionStore)

	// Initialize the user service handler
	handler := handlers.NewUserServiceHandler(service)

	log.Println("User-service started, waiting for requests...")
	go handler.ListenMessages()

	select {}
}
