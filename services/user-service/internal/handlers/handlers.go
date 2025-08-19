package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/lucas/gokafka/user-service/internal/services"
	"github.com/segmentio/kafka-go"
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

		var resp models.Response
		switch req.Type {
		case "register":
			// Parse the registration request from payload
			var registerReq models.RegisterRequest
			if err := json.Unmarshal([]byte(req.Payload), &registerReq); err != nil {
				log.Printf("Failed to parse register request: %v", err)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         "Invalid registration request format",
				}
			} else {
				// Handle registration logic here
				result, err := h.service.RegisterUser(registerReq)
				if err != nil {
					log.Printf("Registration failed: %v", err)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       false,
						Error:         err.Error(),
					}
				} else {
					// Return success response with user data
					resultBytes, _ := json.Marshal(result)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       true,
						Data:          string(resultBytes),
					}
				}
			}
		case "login":
			// Parse the login request from payload
			var loginReq models.LoginRequest
			if err := json.Unmarshal([]byte(req.Payload), &loginReq); err != nil {
				log.Printf("Failed to parse login request: %v", err)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         "Invalid login request format",
				}
			} else {
				// Handle login logic here
				result, err := h.service.LoginUser(loginReq)
				if err != nil {
					log.Printf("Login failed: %v", err)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       false,
						Error:         err.Error(),
					}
				} else {
					// Return success response with user data
					resultBytes, _ := json.Marshal(result)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       true,
						Data:          string(resultBytes),
					}
				}
			}

		case "get-user-profile":
			// Parse the request from payload
			var getProfileReq models.GetProfileRequest
			if err := json.Unmarshal([]byte(req.Payload), &getProfileReq); err != nil {
				log.Printf("Failed to parse get profile request: %v", err)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         "Invalid get profile request format",
				}
			} else {
				// Handle get profile logic here
				result, err := h.service.GetUserProfile(getProfileReq.ID)
				if err != nil {
					log.Printf("Failed to get user profile: %v", err)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       false,
						Error:         err.Error(),
					}
				} else {
					// Create response structure
					profileResponse := models.GetProfileResponse{
						Status: "success",
						Data:   *result,
					}
					// Return success response with user data
					resultBytes, _ := json.Marshal(profileResponse)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       true,
						Data:          string(resultBytes),
					}
				}
			}
		case "logout":
			// Handle logout
			resp = models.Response{
				CorrelationID: req.CorrelationID,
				Data:          "User logged out: " + req.Payload,
			}
		case "get-by-id":
			// Parse the request from payload
			var getProfileReq models.GetProfileRequest
			if err := json.Unmarshal([]byte(req.Payload), &getProfileReq); err != nil {
				log.Printf("Failed to parse request: %v", err)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         "Invalid request format",
				}
			} else {
				// Handle logic here
				result, err := h.service.GetUserProfile(getProfileReq.ID)
				if err != nil {
					log.Printf("Failed to get user profile: %v", err)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       false,
						Error:         err.Error(),
					}
				} else {
					// Return success response with user data
					resultBytes, _ := json.Marshal(result)
					resp = models.Response{
						CorrelationID: req.CorrelationID,
						Success:       true,
						Data:          string(resultBytes),
					}
				}
			}
		case "list-user-profiles":
			// Handle get profile logic here
			result, err := h.service.GetAllUserProfile()
			if err != nil {
				log.Printf("Failed to list user profiles: %v", err)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         err.Error(),
				}
			} else {
				// Create response structure
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

				// Return success response with user list
				resultBytes, _ := json.Marshal(profileListResponse)
				resp = models.Response{
					CorrelationID: req.CorrelationID,
					Success:       true,
					Data:          string(resultBytes),
				}
			}
		default:
			resp = models.Response{
				CorrelationID: req.CorrelationID,
				Data:          "Unknown request type: " + req.Type,
			}
		}

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
