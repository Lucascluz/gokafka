package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/api-gateway/internal/handlers"
	"github.com/lucas/gokafka/api-gateway/internal/middleware"
)

func main() {
	router := gin.Default()
	handlers := handlers.NewHandler()

	// Public routes
	router.GET("/test", handlers.Test)

	auth := router.Group("api/v1/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)

		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/logout", handlers.LogoutUser)
		}
	}

	// Protected routes for authenticated users
	api := router.Group("api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", handlers.GetUserProfile)
		api.PUT("/profile", handlers.UpdateUserProfile)
	}

	// Admin-only routes
	admin := router.Group("api/v1/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/users", handlers.ListUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)
	}

	router.Run()
}
