package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	shared "github.com/lucas/gokafka/shared/models"
	"github.com/segmentio/kafka-go"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	// Parse request body
	var registerReq shared.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if registerReq.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if registerReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}
	if registerReq.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name is required"})
		return
	}
	if registerReq.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Last name is required"})
		return
	}

	// Generate correlation ID for tracking the request
	correlationID := uuid.NewString()
	replyChan := make(chan []byte, 1)

	// Store the reply channel
	h.mu.Lock()
	h.responseChans[correlationID] = replyChan
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.responseChans, correlationID)
		h.mu.Unlock()
	}()

	// Convert register request to JSON string for payload
	payloadBytes, err := json.Marshal(registerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to serialize request",
			"details": err.Error(),
		})
		return
	}

	// Create Kafka request message
	kafkaReq := shared.Request{
		Type:          "register",
		CorrelationID: correlationID,
		ReplyTo:       "user-service-topic",
		Payload:       string(payloadBytes),
	}

	// Send message to user-service via Kafka
	reqBytes, _ := json.Marshal(kafkaReq)
	err = h.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("user-register"),
			Value: reqBytes,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send message to user service",
			"details": err.Error(),
		})
		return
	}

	// Wait for response from user-service
	select {
	case resp := <-replyChan:
		var respObj shared.Response
		if err := json.Unmarshal(resp, &respObj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invalid response format from user service",
				"details": err.Error(),
			})
			return
		}

		// Try to parse the response data as JSON, otherwise return as string
		var responseData any
		if err := json.Unmarshal([]byte(respObj.Data), &responseData); err != nil {
			// If it's not valid JSON, return as string
			responseData = respObj.Data
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "User registration processed",
			"correlation_id": respObj.CorrelationID,
			"data":           responseData,
		})

	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Timeout waiting for response from user service",
		})
	}
}

func (h *Handler) LoginUser(c *gin.Context) {
	// Parse request body
	var loginReq shared.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if loginReq.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if loginReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	// Generate correlation ID for tracking the request
	correlationID := uuid.NewString()
	replyChan := make(chan []byte, 1)

	// Store the reply channel
	h.mu.Lock()
	h.responseChans[correlationID] = replyChan
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.responseChans, correlationID)
		h.mu.Unlock()
	}()

	// Convert login request to JSON string for payload
	payloadBytes, err := json.Marshal(loginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to serialize request",
			"details": err.Error(),
		})
		return
	}

	// Create Kafka request message
	kafkaReq := shared.Request{
		Type:          "login",
		CorrelationID: correlationID,
		ReplyTo:       "user-service-topic",
		Payload:       string(payloadBytes),
	}

	// Send message to user-service via Kafka
	reqBytes, _ := json.Marshal(kafkaReq)
	err = h.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("user-login"),
			Value: reqBytes,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send message to user service",
			"details": err.Error(),
		})
		return
	}

	// Wait for response from user-service
	select {
	case resp := <-replyChan:
		var respObj shared.Response
		if err := json.Unmarshal(resp, &respObj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invalid response format from user service",
				"details": err.Error(),
			})
			return
		}

		// Check if login was successful
		if respObj.Success {
			// Parse the login response which should contain the token
			var loginResponse shared.LoginResponse
			if err := json.Unmarshal([]byte(respObj.Data), &loginResponse); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Invalid login response format",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   loginResponse.Token,
				"user": gin.H{
					"id":    loginResponse.User.ID,
					"email": loginResponse.User.Email,
					"role":  loginResponse.User.Role,
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": respObj.Error,
			})
		}

	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Timeout waiting for response from user service",
		})
	}
}

func (h *Handler) LogoutUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}
