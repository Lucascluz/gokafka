package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/product-service/internal/handlers"
	"github.com/lucas/gokafka/shared/utils"
)

const (
	DefaultPort = "8082"
)

func main() {
	log.Println("Starting product-service...")

	// Initialize handler
	handler := handlers.NewProductHandler()

	log.Println("Product-service started, waiting for requests...")
	
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
