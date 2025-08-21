package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/product-service/internal/handlers"
	"github.com/lucas/gokafka/shared/utils"
)

func main() {

	log.Printf("Starting product-service")

	handlers := handlers.NewProductHandler()

	go handlers.ListenMessages()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	port := utils.GetEnvOrDefault("PORT", "8082")
	router.Run(port)
}
