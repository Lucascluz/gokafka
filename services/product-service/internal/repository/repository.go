package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/lucas/gokafka/product-service/internal/models"
	sharedModels "github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
)

// Database constants
const (
	DefaultPostgresHost     = "localhost"
	DefaultPostgresPort     = "5432"
	DefaultPostgresDB       = "gokafka"
	DefaultPostgresUser     = "postgres"
	DefaultPostgresPassword = "postgres"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository() *ProductRepository {
	db := initDatabase()
	repo := &ProductRepository{db: db}
	repo.createTableIfNotExists()
	return repo
}

// initDatabase initializes the database connection
func initDatabase() *sql.DB {
	host := utils.GetEnvOrDefault("POSTGRES_HOST", DefaultPostgresHost)
	port := utils.GetEnvOrDefault("POSTGRES_PORT", DefaultPostgresPort)
	dbname := utils.GetEnvOrDefault("POSTGRES_DB", DefaultPostgresDB)
	user := utils.GetEnvOrDefault("POSTGRES_USER", DefaultPostgresUser)
	password := utils.GetEnvOrDefault("POSTGRES_PASSWORD", DefaultPostgresPassword)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	log.Printf("Connecting to PostgreSQL at %s:%s", host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
	return db
}

func (r *ProductRepository) createTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	if _, err := r.db.Exec(query); err != nil {
		log.Fatalf("failed to create products table: %v", err)
	}
	
	log.Println("Products table is ready")
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, 
		time.Now(), time.Now()).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	
	return nil
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	query := `
	SELECT id, name, description, price, created_at, updated_at 
	FROM products WHERE id = $1`
	
	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, 
		&product.Price, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]*models.Product, error) {
	query := `
	SELECT id, name, description, price, created_at, updated_at 
	FROM products ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description,
			&product.Price, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	
	return products, nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := `
	UPDATE products 
	SET name = $1, description = $2, price = $3, updated_at = $4
	WHERE id = $5
	RETURNING updated_at`
	
	err := r.db.QueryRow(query, product.Name, product.Description, 
		product.Price, time.Now(), product.ID).Scan(&product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	
	return nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	
	return nil
}

// Helper method to convert Product to ProductData
func (r *ProductRepository) ProductToProductData(product *models.Product) *sharedModels.ProductData {
	return &sharedModels.ProductData{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}
