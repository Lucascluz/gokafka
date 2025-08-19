package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedModels "github.com/lucas/gokafka/shared/models"
)

func (h *Handler) RegisterUser(c *gin.Context) {

	// Initialize helper services
	validator := NewValidator(c)
	respHandler := NewResponseHandler(c)
	messaging := NewMessagingService(h)

	// Parse and validate request
	var registerReq sharedModels.RegisterRequest
	if err := validator.BindJSON(&registerReq); err != nil {
		return
	}

	// Validate required fields
	if err := validator.ValidateRequired(map[string]interface{}{
		"Email":     registerReq.Email,
		"Password":  registerReq.Password,
		"FirstName": registerReq.FirstName,
		"LastName":  registerReq.LastName,
	}); err != nil {
		return
	}

	// Send request to user service
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "register",
		Payload: registerReq,
		Key:     "user-register",
		ReplyTo: "user-service-topic",
	})
	if err != nil {
		respHandler.HandleError(http.StatusInternalServerError, "Failed to send message to user service", err.Error())
		return
	}

	// Handle service response
	respHandler.HandleServiceResponse(resp, "User registration processed")
}

func (h *Handler) LoginUser(c *gin.Context) {
	// Initialize helper services
	validator := NewValidator(c)
	respHandler := NewResponseHandler(c)
	messaging := NewMessagingService(h)

	// Parse and validate request
	var loginReq sharedModels.LoginRequest
	if err := validator.BindJSON(&loginReq); err != nil {
		return
	}

	// Validate required fields
	if err := validator.ValidateRequired(map[string]interface{}{
		"Email":    loginReq.Email,
		"Password": loginReq.Password,
	}); err != nil {
		return
	}

	// Send request to user service
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "login",
		Payload: loginReq,
		Key:     "user-login",
		ReplyTo: "user-service-topic",
	})
	if err != nil {
		respHandler.HandleError(http.StatusInternalServerError, "Failed to send message to user service", err.Error())
		return
	}

	// Handle login-specific response
	if resp.Success {
		// Parse the login response which should contain the token
		var loginResponse sharedModels.LoginResponse
		if err := json.Unmarshal([]byte(resp.Data), &loginResponse); err != nil {
			respHandler.HandleError(http.StatusInternalServerError, "Invalid login response format", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   loginResponse.Token,
			"data": gin.H{
				"id":         loginResponse.Data.ID,
				"email":      loginResponse.Data.Email,
				"first_name": loginResponse.Data.FirstName,
				"last_name":  loginResponse.Data.LastName,
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": resp.Error,
		})
	}
}

func (h *Handler) LogoutUser(c *gin.Context) {
	// Get token ID from context (set by middleware)
	tokenID, exists := c.Get("token_id")
	log.Printf("Token ID from context: %v, exists: %v", tokenID, exists)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token ID not found in context"})
		return
	}

	tokenIDStr, ok := tokenID.(string)
	log.Printf("Token ID string conversion: %v, ok: %v", tokenIDStr, ok)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID format"})
		return
	}

	// Get token expiration from context
	exp, exists := c.Get("token_exp")
	log.Printf("Token exp from context: %v, exists: %v", exp, exists)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expiration not found"})
		return
	}

	expTime, ok := exp.(int64)
	log.Printf("Token exp conversion: %v, ok: %v", expTime, ok)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token expiration format"})
		return
	}

	// Calculate remaining time until token expires
	expiration := time.Until(time.Unix(expTime, 0))
	log.Printf("Token expiration duration: %v", expiration)
	if expiration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token already expired"})
		return
	}

	// Blacklist the token
	if err := h.blacklist.BlacklistToken(tokenIDStr, expiration); err != nil {
		log.Printf("Failed to blacklist token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	log.Printf("Successfully blacklisted token: %s", tokenIDStr)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
