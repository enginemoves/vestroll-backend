package auth

import (
	"net/http"
	"strings"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	authService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

// Login handles POST /api/v1/auth/login
func (h *LoginHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}
	
	// Authenticate user
	response, err := h.authService.Login(req)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusUnauthorized
		errorCode := "authentication_failed"
		
		// Handle specific error cases
		if strings.Contains(err.Error(), "invalid credentials") {
			statusCode = http.StatusUnauthorized
			errorCode = "invalid_credentials"
		} else if strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		} else {
			statusCode = http.StatusInternalServerError
			errorCode = "internal_error"
		}
		
		c.JSON(statusCode, models.ErrorResponse{
			Error:   errorCode,
			Message: "Authentication failed",
		})
		return
	}
	
	// Return successful login response
	c.JSON(http.StatusOK, response)
}