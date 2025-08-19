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

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}

func (h *Handler) ListUserProfiles(c *gin.Context) {

	// Initialize helper services
	respHandler := NewResponseHandler(c)
	messaging := NewMessagingService(h)

	// Get user role from the token in the auth header
	userRole, exists := c.Get("user_role")
	if !exists || userRole != "admin" {
		respHandler.HandleError(http.StatusForbidden, "Only admins can list user profiles")
		return
	}

	// Send request to user service
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "list-user-profiles",
		Payload: "",
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
		var profileRes shared.ListProfileResponse
		if err := json.Unmarshal([]byte(resp.Data), &profileRes); err != nil {
			respHandler.HandleError(http.StatusInternalServerError, "Invalid response format", err.Error())
			return
		}

		// Return the list of user profiles
		var profiles []gin.H
		for _, userData := range profileRes.Data {
			profiles = append(profiles, gin.H{
				"id":         userData.ID,
				"email":      userData.Email,
				"first_name": userData.FirstName,
				"last_name":  userData.LastName,
				"created_at": userData.CreatedAt,
				"updated_at": userData.UpdatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User profiles retrieved successfully",
			"data":    profiles,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resp.Error,
		})
	}
}

func (h *Handler) DeleteUserProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "method not implemented"})
}
