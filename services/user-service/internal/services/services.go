package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas/gokafka/shared/models"
	"github.com/lucas/gokafka/user-service/internal/auth"
	userModels "github.com/lucas/gokafka/user-service/internal/models"
	"github.com/lucas/gokafka/user-service/internal/repository"
	"github.com/lucas/gokafka/user-service/internal/session"
)

type UserService struct {
	repo         repository.UserRepository
	sessionStore *session.RedisSessionStore
}

func NewUserService(repo *repository.UserRepository, sessionStore *session.RedisSessionStore) *UserService {
	return &UserService{
		repo:         *repo,
		sessionStore: sessionStore,
	}
}

func (s *UserService) RegisterUser(req models.RegisterRequest) (*userModels.User, error) {
	// Validate input
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &userModels.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Role:      "user", // Default role
	}

	// Save user to repository
	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create a copy of user without password for security
	userResponse := &userModels.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  "", // Don't return password
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Role:      user.Role,
	}
	return userResponse, nil
}

func (s *UserService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session in Redis
	sessionID, err := s.sessionStore.CreateSession(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Prepare response
	loginResponse := &models.LoginResponse{
		Token:     token,
		SessionID: sessionID,
		User: models.User{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return loginResponse, nil
}
