package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
	sharedModels "github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/shared/utils"
	"github.com/lucas/gokafka/user-service/internal/auth"
	"github.com/lucas/gokafka/user-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
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
	repo := &UserRepository{db: db}
	repo.createTableIfNotExists()

	return repo
}

func (r *UserRepository) createTableIfNotExists() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user'
	)`

	if _, err := r.db.Exec(query); err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}

	// Ensure admin user exists
	if err := r.InsertAdminUser(); err != nil {
		log.Fatalf("failed to ensure admin user exists: %v", err)
	}

	// Ensure users table is ready
	log.Println("Users table ready")
}

func (r *UserRepository) InsertAdminUser() error {

	// Create base admin user
	adminID := uuid.New().String()
	adminEmail := utils.GetEnvOrDefault("ADMIN_EMAIL", "admin@example.com")

	hashedAdminPassword, err := auth.HashPassword("admin123")
	if err != nil {
		log.Fatalf("failed to hash default admin password: %v", err)
	}

	adminPassword := utils.GetEnvOrDefault("ADMIN_PASSWORD", hashedAdminPassword)
	adminFirstName := utils.GetEnvOrDefault("ADMIN_FIRST_NAME", "Admin")
	adminLastName := utils.GetEnvOrDefault("ADMIN_LAST_NAME", "User")
	adminCreatedAt := utils.GetEnvOrDefault("ADMIN_CREATED_AT", time.Now().Format(time.RFC3339))
	adminUpdatedAt := utils.GetEnvOrDefault("ADMIN_UPDATED_AT", time.Now().Format(time.RFC3339))

	query := `
	INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at, role)
	VALUES ($1, $2, $3, $4, $5, $6, $7, 'admin')
	`

	if _, err := r.db.Exec(query, adminID, adminEmail, adminPassword, adminFirstName, adminLastName, adminCreatedAt, adminUpdatedAt); err != nil {
		if err.Error() != "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			log.Fatalf("failed to insert admin user: %v", err)
			return err
		} else {
			log.Println("Admin user already exists, skipping insert")
			return nil

		}
	} else {
		log.Println("Admin user created successfully")
		return nil
	}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, created_at, updated_at, role 
		FROM users WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Password, user.FirstName,
		user.LastName, user.CreatedAt, user.UpdatedAt, user.Role,
	)

	return err
}

func (r *UserRepository) GetUserByID(id string) ([]*sharedModels.UserData, error) {
	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at 
		FROM users WHERE id = $1
	`

	var user sharedModels.UserData
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName,
		&user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return []*sharedModels.UserData{&user}, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, created_at, updated_at, role 
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.FirstName,
			&user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.Role,
		)
		if err != nil {
			return nil, err
		}
		// Don't return passwords
		user.Password = ""
		users = append(users, &user)
	}

	return users, nil
}
