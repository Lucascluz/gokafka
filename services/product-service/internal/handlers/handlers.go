package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lucas/gokafka/product-service/internal/models"
	"github.com/lucas/gokafka/product-service/internal/service"
	sharedModels "github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/segmentio/kafka-go"
)

type ProductHandler struct {
	writer  *kafka.Writer
	reader  *kafka.Reader
	service *service.ProductService
}

func NewProductHandler() *ProductHandler {
	broker := utils.GetEnvOrDefault("KAFKA_BROKERS", "localhost:9092")

	return &ProductHandler{
		service: service.NewProductService(),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{broker},
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   "api-gateway-topic",
			GroupID: "product-service-group",
		}),
	}
}

func (h *ProductHandler) ListenMessages() {
	for {
		m, err := h.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		var req sharedModels.Request
		if err := json.Unmarshal(m.Value, &req); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		var resp sharedModels.Response
		switch req.Type {
		case "new-product":
			// Parse the registration request from payload
			var registerReq models.CreateProductRequest
			if err := json.Unmarshal([]byte(req.Payload), &registerReq); err != nil {
				log.Printf("Failed to parse register request: %v", err)
				resp = sharedModels.Response{
					CorrelationID: req.CorrelationID,
					Success:       false,
					Error:         "Invalid create product request format",
				}
			} else {
				// Handle registration logic here
				result, err := h.service.CreateProduct()
				if err != nil {
					log.Printf("Registration failed: %v", err)
					resp = sharedModels.Response{
						CorrelationID: req.CorrelationID,
						Success:       false,
						Error:         err.Error(),
					}
				} else {
					// Return success response with user data
					resultBytes, _ := json.Marshal(result)
					resp = sharedModels.Response{
						CorrelationID: req.CorrelationID,
						Success:       true,
						Data:          string(resultBytes),
					}
				}
			}
		default:
			resp = sharedModels.Response{
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
