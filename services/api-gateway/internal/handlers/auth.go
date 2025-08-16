package handlers

import (
	"encoding/json"
	"net/http"

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
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}
