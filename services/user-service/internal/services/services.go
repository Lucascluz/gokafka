package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas/gokafka/shared/auth"
	sharedModels "github.com/lucas/gokafka/shared/models"
	userAuth "github.com/lucas/gokafka/user-service/internal/auth"
	userModels "github.com/lucas/gokafka/user-service/internal/models"
	"github.com/lucas/gokafka/user-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: *repo,
	}
}

func (s *UserService) RegisterUser(req sharedModels.RegisterRequest) (*userModels.User, error) {
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
	hashedPassword, err := userAuth.HashPassword(req.Password)
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

func (s *UserService) LoginUser(req sharedModels.LoginRequest) (*sharedModels.LoginResponse, error) {
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
	if !userAuth.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Prepare response
	loginResponse := &sharedModels.LoginResponse{
		Token: token,
		Data: sharedModels.UserData{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	return loginResponse, nil
}

func (s *UserService) GetUserProfile(userID string) (*sharedModels.UserData, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get user by ID
	users, err := s.repo.GetUserByID(userID)
	if err != nil || len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	user := users[0]

	// Create a copy of user without password for security
	userResponse := &sharedModels.UserData{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
}

func (s *UserService) GetAllUserProfile() ([]*sharedModels.UserData, error) {

	// Get user by ID
	users, err := s.repo.GetAllUsers()
	if err != nil || len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}

	var userListResponse []*sharedModels.UserData
	for _, user := range users {
		// Create a copy of user without password for security
		userResponse := &sharedModels.UserData{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		userListResponse = append(userListResponse, userResponse)
	}

	return userListResponse, nil
}
