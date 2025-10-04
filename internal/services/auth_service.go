package services

import (
	"fmt"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/utils"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user with email and password
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	
	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	// Generate JWT token
	token, err := s.jwtService.GenerateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	
	// Return successful login response
	return &models.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    *user,
	}, nil
}