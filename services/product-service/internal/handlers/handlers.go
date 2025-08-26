package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/lucas/gokafka/product-service/internal/service"
	sharedModels "github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/segmentio/kafka-go"
)

// Request type constants
const (
	RequestTypeHealth            = "health"
	RequestTypeCreateProduct     = "create-product"
	RequestTypeGetProduct        = "get-product"
	RequestTypeGetProductByID    = "get-product-by-id"
	RequestTypeListProducts      = "list-products"
	RequestTypeUpdateProduct     = "update-product"
	RequestTypeDeleteProduct     = "delete-product"
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

// Helper method to create error responses
func (h *ProductHandler) createErrorResponse(correlationID, errorMsg string) sharedModels.Response {
	return sharedModels.Response{
		CorrelationID: correlationID,
		Success:       false,
		Error:         errorMsg,
	}
}

// Helper method to create success responses
func (h *ProductHandler) createSuccessResponse(correlationID string, data interface{}) sharedModels.Response {
	dataBytes, _ := json.Marshal(data)
	return sharedModels.Response{
		CorrelationID: correlationID,
		Success:       true,
		Data:          string(dataBytes),
	}
}

// Helper method to unmarshal request payload
func (h *ProductHandler) unmarshalPayload(payload string, target interface{}) error {
	return json.Unmarshal([]byte(payload), target)
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
func (h *ProductHandler) handleRequest(req sharedModels.Request) sharedModels.Response {
	switch req.Type {
	case RequestTypeHealth:
		return h.handleHealth(req.CorrelationID)
	case RequestTypeCreateProduct:
		return h.handleCreateProduct(req)
	case RequestTypeGetProduct, RequestTypeGetProductByID:
		return h.handleGetProduct(req)
	case RequestTypeListProducts:
		return h.handleListProducts(req.CorrelationID)
	case RequestTypeUpdateProduct:
		return h.handleUpdateProduct(req)
	case RequestTypeDeleteProduct:
		return h.handleDeleteProduct(req)
	default:
		return sharedModels.Response{
			CorrelationID: req.CorrelationID,
			Data:          "Unknown request type: " + req.Type,
		}
	}
}

// handleHealth returns health status
func (h *ProductHandler) handleHealth(correlationID string) sharedModels.Response {
	healthResponse := map[string]interface{}{
		"service":   "product-service",
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	}
	return h.createSuccessResponse(correlationID, healthResponse)
}

// handleCreateProduct processes product creation
func (h *ProductHandler) handleCreateProduct(req sharedModels.Request) sharedModels.Response {
	var createReq sharedModels.CreateProductRequest
	if err := h.unmarshalPayload(req.Payload, &createReq); err != nil {
		log.Printf("Failed to parse create product request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid create product request format")
	}

	result, err := h.service.CreateProduct(createReq)
	if err != nil {
		log.Printf("Product creation failed: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	response := sharedModels.ProductResponse{
		Status:  "success",
		Message: "Product created successfully",
		Data:    *result,
	}
	return h.createSuccessResponse(req.CorrelationID, response)
}

// handleGetProduct processes get product by ID request
func (h *ProductHandler) handleGetProduct(req sharedModels.Request) sharedModels.Response {
	var getReq sharedModels.GetProductRequest
	if err := h.unmarshalPayload(req.Payload, &getReq); err != nil {
		log.Printf("Failed to parse get product request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid get product request format")
	}

	result, err := h.service.GetProductByID(getReq.ID)
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	response := sharedModels.GetProductResponse{
		Status: "success",
		Data:   *result,
	}
	return h.createSuccessResponse(req.CorrelationID, response)
}

// handleListProducts processes list all products request
func (h *ProductHandler) handleListProducts(correlationID string) sharedModels.Response {
	result, err := h.service.GetAllProducts()
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		return h.createErrorResponse(correlationID, err.Error())
	}

	// Convert []*ProductData to []ProductData
	productDataVals := make([]sharedModels.ProductData, len(result))
	for i, p := range result {
		if p != nil {
			productDataVals[i] = *p
		}
	}

	response := sharedModels.ListProductResponse{
		Status: "success",
		Data:   productDataVals,
	}
	return h.createSuccessResponse(correlationID, response)
}

// handleUpdateProduct processes product update
func (h *ProductHandler) handleUpdateProduct(req sharedModels.Request) sharedModels.Response {
	var updateReq sharedModels.UpdateProductRequest
	if err := h.unmarshalPayload(req.Payload, &updateReq); err != nil {
		log.Printf("Failed to parse update product request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid update product request format")
	}

	result, err := h.service.UpdateProduct(updateReq)
	if err != nil {
		log.Printf("Product update failed: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	response := sharedModels.ProductResponse{
		Status:  "success",
		Message: "Product updated successfully",
		Data:    *result,
	}
	return h.createSuccessResponse(req.CorrelationID, response)
}

// handleDeleteProduct processes product deletion
func (h *ProductHandler) handleDeleteProduct(req sharedModels.Request) sharedModels.Response {
	var deleteReq sharedModels.DeleteProductRequest
	if err := h.unmarshalPayload(req.Payload, &deleteReq); err != nil {
		log.Printf("Failed to parse delete product request: %v", err)
		return h.createErrorResponse(req.CorrelationID, "Invalid delete product request format")
	}

	err := h.service.DeleteProduct(deleteReq.ID)
	if err != nil {
		log.Printf("Product deletion failed: %v", err)
		return h.createErrorResponse(req.CorrelationID, err.Error())
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Product deleted successfully",
		"id":      deleteReq.ID,
	}
	return h.createSuccessResponse(req.CorrelationID, response)
}
