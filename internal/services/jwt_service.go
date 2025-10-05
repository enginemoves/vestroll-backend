package services

import (
	"fmt"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config config.JWTConfig
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService(config config.JWTConfig) *JWTService {
	return &JWTService{config: config}
}

// GenerateToken generates a JWT token for the given user
func (j *JWTService) GenerateToken(user models.User) (string, error) {
	// Set token expiration
	expirationTime := time.Now().Add(j.config.TTL)
	
	// Create claims
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "vestroll-backend",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}
	
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token with secret
	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	
	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token claims")
}