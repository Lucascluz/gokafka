package service

import (
	"fmt"

	"github.com/lucas/gokafka/product-service/internal/models"
	"github.com/lucas/gokafka/product-service/internal/repository"
	sharedModels "github.com/lucas/gokafka/shared/models"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		repo: repository.NewProductRepository(),
	}
}

// Helper method to convert Product to ProductData
func (s *ProductService) productToProductData(product *models.Product) *sharedModels.ProductData {
	return s.repo.ProductToProductData(product)
}

func (s *ProductService) CreateProduct(req sharedModels.CreateProductRequest) (*sharedModels.ProductData, error) {
	// Validate input
	if req.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("product price must be greater than 0")
	}

	// Create new product
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	// Save product to repository
	if err := s.repo.CreateProduct(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.productToProductData(product), nil
}

func (s *ProductService) GetProductByID(id int) (*sharedModels.ProductData, error) {
	// Validate input
	if id <= 0 {
		return nil, fmt.Errorf("invalid product ID")
	}

	// Get product from repository
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	return s.productToProductData(product), nil
}

func (s *ProductService) GetAllProducts() ([]*sharedModels.ProductData, error) {
	// Get all products from repository
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// Convert to ProductData
	productDataList := make([]*sharedModels.ProductData, 0, len(products))
	for _, product := range products {
		if product != nil {
			productDataList = append(productDataList, s.productToProductData(product))
		}
	}

	return productDataList, nil
}

func (s *ProductService) UpdateProduct(req sharedModels.UpdateProductRequest) (*sharedModels.ProductData, error) {
	// Validate input
	if req.ID <= 0 {
		return nil, fmt.Errorf("invalid product ID")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if req.Price <= 0 {
		return nil, fmt.Errorf("product price must be greater than 0")
	}

	// Check if product exists
	existingProduct, err := s.repo.GetProductByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Update product fields
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = req.Price

	// Update in repository
	if err := s.repo.UpdateProduct(existingProduct); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return s.productToProductData(existingProduct), nil
}

func (s *ProductService) DeleteProduct(id int) error {
	// Validate input
	if id <= 0 {
		return fmt.Errorf("invalid product ID")
	}

	// Delete from repository
	if err := s.repo.DeleteProduct(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
