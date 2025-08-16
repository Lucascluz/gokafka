package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	shared "github.com/lucas/gokafka/shared/models"
)

func (h *Handler) GetUserProfile(c *gin.Context) {
	// Initialize helper services
	respHandler := NewResponseHandler(c)
	messaging := NewMessagingService(h)

	// Get user ID from context (set by middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		respHandler.HandleError(http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		respHandler.HandleError(http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	// Create request with just the user ID
	profileReq := shared.GetProfileRequest{
		ID: userIDStr,
	}

	// Send request to user service
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "get-user-profile",
		Payload: profileReq,
		Key:     "get-user-profile",
		ReplyTo: "user-service-topic",
	})
	if err != nil {
		respHandler.HandleError(http.StatusInternalServerError, "Failed to send message to user service", err.Error())
		return
	}

	// Handle profile-specific response
	if resp.Success {
		// Parse the profile response
		var profileRes shared.GetProfileResponse
		if err := json.Unmarshal([]byte(resp.Data), &profileRes); err != nil {
			respHandler.HandleError(http.StatusInternalServerError, "Invalid response format", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User profile retrieved successfully",
			"data": gin.H{
				"id":         profileRes.Data.ID,
				"email":      profileRes.Data.Email,
				"first_name": profileRes.Data.FirstName,
				"last_name":  profileRes.Data.LastName,
				"created_at": profileRes.Data.CreatedAt,
				"updated_at": profileRes.Data.UpdatedAt,
			},
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resp.Error,
		})
	}
}

func (h *Handler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}
