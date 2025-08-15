package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	correlationID := uuid.NewString()
	replyChan := make(chan []byte, 1)

	h.mu.Lock()
	h.responseChans[correlationID] = replyChan
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.responseChans, correlationID)
		h.mu.Unlock()
	}()

	req := map[string]interface{}{
		"correlation_id": correlationID,
		"reply_to":       "user-service-topic",
		"payload":        "hello world",
	}
	reqBytes, _ := json.Marshal(req)

	err := h.writer.WriteMessages(context.Background(),
		kafka.Message{Key: []byte("key"), Value: reqBytes},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	select {
	case resp := <-replyChan:
		var respObj map[string]interface{}
		if err := json.Unmarshal(resp, &respObj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response format"})
			return
		}
		c.JSON(http.StatusOK, respObj)
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "timeout waiting for response"})
	}
}


