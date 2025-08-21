package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/lucas/gokafka/product-service/internal/models"
	"github.com/lucas/gokafka/shared/utils"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository() *ProductRepository {

	// Build connection string from environment variables
	host := utils.GetEnvOrDefault("POSTGRES_HOST", "localhost")
	port := utils.GetEnvOrDefault("POSTGRES_PORT", "5432")
	dbname := utils.GetEnvOrDefault("POSTGRES_DB", "gokafka")
	user := utils.GetEnvOrDefault("POSTGRES_USER", "postgres")
	password := utils.GetEnvOrDefault("POSTGRES_PASSWORD", "postgres")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	log.Printf("Connecting to PostgreSQL at %s:%s", host, port)
	db, err := sql.Open("postgres", connStr)

	// Check for connection errors
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	// Create users table if it doesn't exist
	repo := &ProductRepository{db: db}
	repo.createTableIfNotExists()

	return repo
}

func (repo *ProductRepository) createTableIfNotExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := repo.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	log.Println("Products table is ready")
	return nil
}

func (repo *ProductRepository) SaveProduct(product *models.Product) error {
	query := `
	INSERT INTO products (name, description, price)
	VALUES ($1, $2, $3)
	RETURNING id`
	err := repo.db.QueryRow(query, product.Name, product.Description, product.Price).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}
	log.Printf("Product saved with ID: %d", product.ID)
	return nil
}
