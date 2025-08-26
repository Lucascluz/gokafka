package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/lucas/gokafka/user-service/internal/services"
	"github.com/segmentio/kafka-go"
)

// Request type constants
const (
	RequestTypeHealth           = "health"
	RequestTypeRegister         = "register"
	RequestTypeLogin            = "login"
	RequestTypeGetUserProfile   = "get-user-profile"
	RequestTypeLogout           = "logout"
	RequestTypeGetByID          = "get-by-id"
	RequestTypeListUserProfiles = "list-user-profiles"
)

type UserServiceHandler struct {
	service *services.UserService
	writer  *kafka.Writer
	reader  *kafka.Reader
}

func NewUserServiceHandler(service *services.UserService) *UserServiceHandler {
	broker := utils.GetEnvOrDefault("KAFKA_BROKERS", "localhost:9092")

	return &UserServiceHandler{
		service: service,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{broker},
			// Remove Topic from here to allow per-message topic specification
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   "api-gateway-topic",
			GroupID: "user-service-group",
		}),
	}
}

// Helper method to create error responses
func (h *UserServiceHandler) createErrorResponse(correlationID, errorMsg string) models.Response {
	return models.Response{
		CorrelationID: correlationID,
		Success:       false,
		Error:         errorMsg,
	}
}

// Helper method to create success responses
func (h *UserServiceHandler) createSuccessResponse(correlationID string, data interface{}) models.Response {
	dataBytes, _ := json.Marshal(data)
	return models.Response{
		CorrelationID: correlationID,
		Success:       true,
		Data:          string(dataBytes),
	}
}

// Helper method to unmarshal request payload
func (h *UserServiceHandler) unmarshalPayload(payload string, target interface{}) error {
	return json.Unmarshal([]byte(payload), target)
}

func (h *UserServiceHandler) ListenMessages() {
	for {
		m, err := h.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		var req models.Request
		if err := json.Unmarshal(m.Value, &req); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		resp := h.handleRequest(req)

		respBytes, _ := json.Marshal(resp)
		err = h.writer.WriteMessages(context.Background(),
			kafka.Message{
				Topic: req.ReplyTo,
				Value: respBytes,
			},
		)
		if err != nil {
			log.Println("write error:", err)
		} else {
			log.Printf("responded to %s with correlation_id %s", req.ReplyTo, req.CorrelationID)
		}
	}
}

// handleRequest processes different request types
func (h *UserServiceHandler) handleRequest(req models.Request) models.Response {
	switch req.Type {
	case RequestTypeHealth:
		return h.handleHealth(req.CorrelationID)
	case RequestTypeRegister:
		return h.handleRegister(req)
	case RequestTypeLogin:
		return h.handleLogin(req)
	case RequestTypeGetUserProfile:
		return h.handleGetUserProfile(req)
	case RequestTypeLogout:
		return h.handleLogout(req)
	case RequestTypeGetByID:
		return h.handleGetById(req)
	case RequestTypeListUserProfiles:
		return h.handleListUserProfiles(req.CorrelationID)
	default:
		return models.Response{
			CorrelationID: req.CorrelationID,
			Data:          "Unknown request type: " + req.Type,
		}
	}
}

// handleHealth returns health status
func (h *UserServiceHandler) handleHealth(correlationID string) models.Response {
	healthResponse := map[string]interface{}{
		"service":   "user-service",
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}
	return h.createSuccessResponse(correlationID, healthResponse)
}

// handleRegister processes user registration
func (h *UserServiceHandler) handleRegister(req models.Request) models.Response {
	var registerReq models.RegisterRequest
	if err := h.unmarshalPayload(req.Payload, &registerReq); err != nil {
		log.Printf("Failed to parse register request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid registration request format")
	}

	result, err := h.service.RegisterUser(registerReq)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	return h.createSuccessResponse(req.CorrelationID, result)
}

// handleLogin processes user login
func (h *UserServiceHandler) handleLogin(req models.Request) models.Response {
	var loginReq models.LoginRequest
	if err := h.unmarshalPayload(req.Payload, &loginReq); err != nil {
		log.Printf("Failed to parse login request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid login request format")
	}

	result, err := h.service.LoginUser(loginReq)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	return h.createSuccessResponse(req.CorrelationID, result)
}

// handleGetUserProfile processes get user profile request
func (h *UserServiceHandler) handleGetUserProfile(req models.Request) models.Response {
	var getProfileReq models.GetProfileRequest
	if err := h.unmarshalPayload(req.Payload, &getProfileReq); err != nil {
		log.Printf("Failed to parse get profile request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid get profile request format")
	}

	result, err := h.service.GetUserProfile(getProfileReq.ID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	profileResponse := models.GetProfileResponse{
		Status: "success",
		Data:   *result,
	}
	return h.createSuccessResponse(req.CorrelationID, profileResponse)
}

// handleLogout processes logout request
func (h *UserServiceHandler) handleLogout(req models.Request) models.Response {
	return models.Response{
		CorrelationID: req.CorrelationID,
		Data:          "User logged out: " + req.Payload,
	}
}

// handleGetById processes get user by ID request
func (h *UserServiceHandler) handleGetById(req models.Request) models.Response {
	var getProfileReq models.GetProfileRequest
	if err := h.unmarshalPayload(req.Payload, &getProfileReq); err != nil {
		log.Printf("Failed to parse request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid request format")
	}

	result, err := h.service.GetUserProfile(getProfileReq.ID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	return h.createSuccessResponse(req.CorrelationID, result)
}

// handleListUserProfiles processes list all user profiles request
func (h *UserServiceHandler) handleListUserProfiles(correlationID string) models.Response {
	result, err := h.service.GetAllUserProfile()
	if err != nil {
		log.Printf("Failed to list user profiles: %v", err)
		return h.createErrorResponse(correlationID, err.Error())
	}

	// Convert []*models.UserData to []models.UserData
	userDataVals := make([]models.UserData, len(result))
	for i, u := range result {
		if u != nil {
			userDataVals[i] = *u
		}
	}
	profileListResponse := models.ListProfileResponse{
		Status: "success",
		Data:   userDataVals,
	}

	return h.createSuccessResponse(correlationID, profileListResponse)
}
