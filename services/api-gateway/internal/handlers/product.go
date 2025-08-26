package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedModels "github.com/lucas/gokafka/shared/models"
)

// Product CRUD handlers

// CreateProduct handles product creation
func (h *Handler) CreateProduct(c *gin.Context) {
	// Parse and validate request
	var req sharedModels.CreateProductRequest
	validator := NewValidator(c)
	if err := validator.BindJSON(&req); err != nil {
		return
	}

	// Validate required fields
	if err := validator.ValidateRequired(map[string]interface{}{
		"name": req.Name,
	}); err != nil {
		return
	}

	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	// Send request to product service
	messaging := NewMessagingService(h)
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "create-product",
		Payload: req,
		Key:     "product-create",
		ReplyTo: "product-service-topic",
		Timeout: 10 * time.Second,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseHandler := NewResponseHandler(c)
	responseHandler.HandleServiceResponse(resp, "Product created successfully")
}

// GetProduct handles getting a product by ID
func (h *Handler) GetProduct(c *gin.Context) {
	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Create request
	req := sharedModels.GetProductRequest{ID: id}

	// Send request to product service
	messaging := NewMessagingService(h)
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "get-product-by-id",
		Payload: req,
		Key:     "product-get",
		ReplyTo: "product-service-topic",
		Timeout: 10 * time.Second,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseHandler := NewResponseHandler(c)
	responseHandler.HandleServiceResponse(resp, "Product retrieved successfully")
}

// ListProducts handles listing all products
func (h *Handler) ListProducts(c *gin.Context) {
	// Send request to product service
	messaging := NewMessagingService(h)
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "list-products",
		Payload: "",
		Key:     "product-list",
		ReplyTo: "product-service-topic",
		Timeout: 10 * time.Second,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseHandler := NewResponseHandler(c)
	responseHandler.HandleServiceResponse(resp, "Products retrieved successfully")
}

// UpdateProduct handles product updates
func (h *Handler) UpdateProduct(c *gin.Context) {
	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Parse and validate request
	var updateData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	validator := NewValidator(c)
	if err := validator.BindJSON(&updateData); err != nil {
		return
	}

	// Validate required fields
	if err := validator.ValidateRequired(map[string]interface{}{
		"name": updateData.Name,
	}); err != nil {
		return
	}

	if updateData.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	// Create update request
	req := sharedModels.UpdateProductRequest{
		ID:          id,
		Name:        updateData.Name,
		Description: updateData.Description,
		Price:       updateData.Price,
	}

	// Send request to product service
	messaging := NewMessagingService(h)
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "update-product",
		Payload: req,
		Key:     "product-update",
		ReplyTo: "product-service-topic",
		Timeout: 10 * time.Second,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseHandler := NewResponseHandler(c)
	responseHandler.HandleServiceResponse(resp, "Product updated successfully")
}

// DeleteProduct handles product deletion
func (h *Handler) DeleteProduct(c *gin.Context) {
	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Create delete request
	req := sharedModels.DeleteProductRequest{ID: id}

	// Send request to product service
	messaging := NewMessagingService(h)
	resp, err := messaging.SendAndWait(SendRequest{
		Type:    "delete-product",
		Payload: req,
		Key:     "product-delete",
		ReplyTo: "product-service-topic",
		Timeout: 10 * time.Second,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseHandler := NewResponseHandler(c)
	responseHandler.HandleServiceResponse(resp, "Product deleted successfully")
}
