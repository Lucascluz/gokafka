package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/api-gateway/internal/cache"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	writer        *kafka.Writer
	readers       []*kafka.Reader
	responseChans map[string]chan []byte
	mu            sync.Mutex
	blacklist     *cache.TokenBlacklist
}

func NewHandler() *Handler {

	broker := utils.GetEnvOrDefault("KAFKA_BROKERS", "localhost:9092")

	readers := []*kafka.Reader{
		// user-service
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   "user-service-topic",
			GroupID: "api-gateway-group",
		}),
		// product-service
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   "product-service-topic",
			GroupID: "api-gateway-group",
		}),
		// add new readers here
	}
	h := &Handler{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{broker},
			Topic:        "api-gateway-topic",
			RequiredAcks: int(kafka.RequireOne),
		}),
		readers:       readers,
		responseChans: make(map[string]chan []byte),
		blacklist:     cache.NewTokenBlacklist(),
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

func (h *Handler) Health(c *gin.Context) {
	// Send health check to all the services
	services := []string{"user-service", "product-service"} // add new services here
	var wg sync.WaitGroup
	responses := make([]map[string]interface{}, len(services))
	errors := make([]string, len(services))
	for i, service := range services {
		wg.Add(1)
		go func(i int, service string) {
			defer wg.Done()
			messaging := NewMessagingService(h)
			resp, err := messaging.SendAndWait(SendRequest{
				Type:    "health",
				Payload: "",
				Key:     "key",
				ReplyTo: service + "-topic",
				Timeout: 5 * time.Second,
			})
			if err != nil {
				errors[i] = service + ": " + err.Error()
				return
			}
			var respObj map[string]interface{}
			if err := json.Unmarshal([]byte(resp.Data), &respObj); err != nil {
				errors[i] = service + ": invalid response format"
				return
			}
			responses[i] = respObj
		}(i, service)
	}
	wg.Wait()

	// Compile final response
	finalResp := map[string]interface{}{
		"status":    "ok",
		"services":  responses,
		"errors":    errors,
		"timestamp": time.Now().UTC(),
	}
	c.JSON(http.StatusOK, finalResp)
}
