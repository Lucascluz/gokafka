package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/user-service/internal/services"
	"github.com/segmentio/kafka-go"
)

type UserServiceHandler struct {
	service *services.UserService
	writer  *kafka.Writer
	reader  *kafka.Reader
}

func NewUserServiceHandler(service *services.UserService) *UserServiceHandler {
	return &UserServiceHandler{
		service: service,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{"localhost:9092"},
			// Remove Topic from here to allow per-message topic specification
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
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
		case "logout":
			// Handle logout
			resp = models.Response{
				CorrelationID: req.CorrelationID,
				Data:          "User logged out: " + req.Payload,
			}
		case "get_all":
			// Handle get all users
			resp = models.Response{
				CorrelationID: req.CorrelationID,
				Data:          "All users: ...", // Replace with actual data
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
