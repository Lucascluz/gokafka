package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	writer        *kafka.Writer
	readers       []*kafka.Reader
	responseChans map[string]chan []byte
	mu            sync.Mutex
}

func NewHandler() *Handler {
	readers := []*kafka.Reader{
		// user-service
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "user-service-topic",
			GroupID: "api-gateway-group",
		}),
		// add new readers here
	}
	h := &Handler{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{"localhost:9092"},
			Topic:        "api-gateway-topic",
			RequiredAcks: int(kafka.RequireOne),
		}),
		readers:       readers,
		responseChans: make(map[string]chan []byte),
	}
	for _, reader := range h.readers {
		go h.listenResponses(reader)
	}
	return h
}

func (h *Handler) listenResponses(reader *kafka.Reader) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		var resp struct {
			CorrelationID string `json:"correlation_id"`
			Data          string `json:"data"`
		}
		if err := json.Unmarshal(m.Value, &resp); err != nil {
			continue
		}
		h.mu.Lock()
		ch, ok := h.responseChans[resp.CorrelationID]
		h.mu.Unlock()
		if ok {
			ch <- m.Value
		}
	}
}

func (h *Handler) Test(c *gin.Context) {
	// Initialize helper services
	respHandler := NewResponseHandler(c)
	messaging := NewMessagingService(h)

	// Send test message
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "test",
		Payload: "hello world",
		Key:     "key",
		ReplyTo: "user-service-topic",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		respHandler.HandleError(http.StatusInternalServerError, err.Error())
		return
	}

	// Return raw response for test endpoint
	var respObj map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Data), &respObj); err != nil {
		respHandler.HandleError(http.StatusInternalServerError, "invalid response format")
		return
	}
	c.JSON(http.StatusOK, respObj)
}
