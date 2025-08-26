package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/lucas/gokafka/user-service/internal/handlers"
	"github.com/lucas/gokafka/user-service/internal/repository"
	"github.com/lucas/gokafka/user-service/internal/services"
)

const (
	DefaultPort = "8081"
)

func main() {
	log.Println("Starting user-service...")

	// Initialize dependencies
	repo := repository.NewUserRepository()
	service := services.NewUserService(repo)
	handler := handlers.NewUserServiceHandler(service)

	log.Println("User-service started, waiting for requests...")

	// Start Kafka message listener in background
	go handler.ListenMessages()

	// Start HTTP server
	startHTTPServer()
}

func startHTTPServer() {
	router := gin.Default()

	// Health endpoints
	router.GET("/health", healthHandler)
	router.GET("/ready", readyHandler)

	port := utils.GetEnvOrDefault("PORT", DefaultPort)
	log.Printf("Starting HTTP server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}

func readyHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ready"})
}
