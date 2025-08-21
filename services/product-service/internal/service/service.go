package service

import (
	"github.com/lucas/gokafka/product-service/internal/models"
	"github.com/lucas/gokafka/product-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		repo: repository.NewProductRepository(),
	}
}

func (s *ProductService) CreateProduct() (*models.Product, error) {
	// Implement product creation logic here
	// This is a placeholder implementation
	product := &models.Product{
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       19.99,
	}

	// Save the product to the repository
	if err := s.repo.SaveProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}
