package services

import (
	"errors"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repository.UserRepository
	jwtCfg   config.JWTConfig
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, jwtCfg config.JWTConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// RegisterUser handles user registration
func (s *UserService) RegisterUser(req models.RegisterRequest) (*models.RegisterResponse, error) {
	// Validate email format (additional validation beyond binding)
	if !isValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password strength
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, errors.New("failed to check email availability")
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user account")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}

	// Clear password from response
	user.Password = ""

	return &models.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User:    *user,
	}, nil
}

// LoginUser handles user login
func (s *UserService) LoginUser(req models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("login failed")
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}

	// Clear password from response
	user.Password = ""

	return &models.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    *user,
	}, nil
}

// generateJWTToken creates a JWT token for the user
func (s *UserService) generateJWTToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(s.jwtCfg.TTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtCfg.Secret))
}

// isValidEmail performs additional email validation
func isValidEmail(email string) bool {
	// Basic email validation - more comprehensive than just binding
	if len(email) < 5 || len(email) > 254 {
		return false
	}
	
	// Check for basic email structure
	atCount := 0
	dotAfterAt := false
	
	for i, char := range email {
		if char == '@' {
			atCount++
			if atCount > 1 {
				return false
			}
		} else if char == '.' && atCount == 1 {
			dotAfterAt = true
		}
	}
	
	return atCount == 1 && dotAfterAt && len(email) > 0
}
