package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/shared/utils"
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

	// Start http server
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	port := utils.GetEnvOrDefault("PORT", "8081")
	router.Run(":" + port)
}
