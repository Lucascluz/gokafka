package main

import (
	"log"

	"github.com/lucas/gokafka/user-service/internal/handlers"
	"github.com/lucas/gokafka/user-service/internal/repository"
	"github.com/lucas/gokafka/user-service/internal/services"
)

func main() {

	log.Println("Starting user-service...")

	// Initialize the user repository
	repo := repository.NewUserRepository()

	// Initialize the user service
	service := services.NewUserService(repo)

	// Initialize the user service handler
	handler := handlers.NewUserServiceHandler(service)

	log.Println("User-service started, waiting for requests...")
	go handler.ListenMessages()

	select {}
}
