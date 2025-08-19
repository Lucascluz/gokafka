package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sharedModels "github.com/lucas/gokafka/shared/models"
	"github.com/segmentio/kafka-go"
)

// MessagingService handles common Kafka messaging operations
type MessagingService struct {
	handler *Handler
}

// NewMessagingService creates a new messaging service
func NewMessagingService(handler *Handler) *MessagingService {
	return &MessagingService{
		handler: handler,
	}
}

// SendRequest represents a request to send via Kafka
type SendRequest struct {
	Type    string      // e.g., "register", "login", "get-user-profile"
	Payload interface{} // The request payload to be marshaled
	Key     string      // Kafka message key
	ReplyTo string      // Topic to reply to
	Timeout time.Duration
}

// SendResponse represents the response from a Kafka request
type SendResponse struct {
	CorrelationID string
	Success       bool
	Data          string
	Error         string
}

// SendAndWait sends a Kafka message and waits for a response
func (ms *MessagingService) SendAndWait(req SendRequest) (*SendResponse, error) {
	// Generate correlation ID for tracking the request
	correlationID := uuid.NewString()
	replyChan := make(chan []byte, 1)

	// Store the reply channel
	ms.handler.mu.Lock()
	ms.handler.responseChans[correlationID] = replyChan
	ms.handler.mu.Unlock()

	// Cleanup reply channel
	defer func() {
		ms.handler.mu.Lock()
		delete(ms.handler.responseChans, correlationID)
		ms.handler.mu.Unlock()
	}()

	// Convert payload to JSON string
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	// Create Kafka request message
	kafkaReq := sharedModels.Request{
		Type:          req.Type,
		CorrelationID: correlationID,
		ReplyTo:       req.ReplyTo,
		Payload:       string(payloadBytes),
	}

	log.Printf("Sending message with correlationID: %s and type: %s", correlationID, req.Type)

	// Send message to service via Kafka
	reqBytes, _ := json.Marshal(kafkaReq)
	err = ms.handler.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(req.Key),
			Value: reqBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Set default timeout if not provided
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Wait for response
	select {
	case resp := <-replyChan:
		var respObj sharedModels.Response
		if err := json.Unmarshal(resp, &respObj); err != nil {
			return nil, fmt.Errorf("invalid response format: %w", err)
		}

		return &SendResponse{
			CorrelationID: respObj.CorrelationID,
			Success:       respObj.Success,
			Data:          respObj.Data,
			Error:         respObj.Error,
		}, nil

	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for response from service")
	}
}

// ResponseHandler handles common response scenarios
type ResponseHandler struct {
	c *gin.Context
}

// NewResponseHandler creates a new response handler
func NewResponseHandler(c *gin.Context) *ResponseHandler {
	return &ResponseHandler{c: c}
}

// HandleError sends an error response
func (rh *ResponseHandler) HandleError(statusCode int, message string, details ...string) {
	response := gin.H{"error": message}
	if len(details) > 0 {
		response["details"] = details[0]
	}
	rh.c.JSON(statusCode, response)
}

// HandleSuccess sends a success response
func (rh *ResponseHandler) HandleSuccess(statusCode int, message string, data interface{}) {
	rh.c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
	})
}

// HandleServiceResponse handles a response from the messaging service
func (rh *ResponseHandler) HandleServiceResponse(resp *SendResponse, successMessage string) {
	if resp.Success {
		// Try to parse the response data as JSON, otherwise return as string
		var responseData interface{}
		if err := json.Unmarshal([]byte(resp.Data), &responseData); err != nil {
			// If it's not valid JSON, return as string
			responseData = resp.Data
		}

		rh.c.JSON(200, gin.H{
			"message":        successMessage,
			"correlation_id": resp.CorrelationID,
			"data":           responseData,
		})
	} else {
		rh.c.JSON(400, gin.H{
			"error": resp.Error,
		})
	}
}

// Validation helpers
type Validator struct {
	c *gin.Context
}

// NewValidator creates a new validator
func NewValidator(c *gin.Context) *Validator {
	return &Validator{c: c}
}

// ValidateRequired checks if required fields are present
func (v *Validator) ValidateRequired(fields map[string]interface{}) error {
	for fieldName, value := range fields {
		if str, ok := value.(string); ok && str == "" {
			v.c.JSON(400, gin.H{"error": fmt.Sprintf("%s is required", fieldName)})
			return fmt.Errorf("%s is required", fieldName)
		}
	}
	return nil
}

// BindJSON binds JSON request and handles errors
func (v *Validator) BindJSON(obj interface{}) error {
	if err := v.c.ShouldBindJSON(obj); err != nil {
		v.c.JSON(400, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return err
	}
	return nil
}
