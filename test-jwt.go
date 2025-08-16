package main

import (
	"fmt"
	"log"

	"github.com/lucas/gokafka/shared/auth"
)

// Simple test to verify JWT functionality
func main() {
	fmt.Println("Testing JWT functionality...")

	// Test token generation
	userID := "test-user-123"
	email := "test@example.com"
	role := "user"

	token, err := auth.GenerateToken(userID, email, role)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	fmt.Printf("Generated token: %s\n", token)

	// Test token validation
	claims, err := auth.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	fmt.Printf("Validated claims:\n")
	fmt.Printf("  UserID: %s\n", claims.UserID)
	fmt.Printf("  Email: %s\n", claims.Email)
	fmt.Printf("  Role: %s\n", claims.Role)

	// Test with invalid token
	_, err = auth.ValidateToken("invalid-token")
	if err != nil {
		fmt.Printf("Successfully rejected invalid token: %v\n", err)
	} else {
		log.Fatal("Failed to reject invalid token!")
	}

	fmt.Println("JWT functionality test completed successfully!")
}
