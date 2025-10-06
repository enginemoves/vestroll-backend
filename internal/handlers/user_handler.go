package handlers

import (
	"net/http"
	"strings"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// UserHandler manages user authentication endpoints
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes registers user endpoints under /auth
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}

// Register handles POST /api/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.RegisterUser(req)
	if err != nil {
		status := http.StatusBadRequest
		code := "validation_error"
		
		// Adjust status codes for known error types
		if strings.Contains(err.Error(), "email already exists") {
			status = http.StatusConflict
			code = "duplicate_email"
		} else if strings.Contains(err.Error(), "internal") || strings.Contains(err.Error(), "failed to") {
			status = http.StatusInternalServerError
			code = "internal_error"
		}
		
		c.JSON(status, models.ErrorResponse{
			Error:   code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		status := http.StatusUnauthorized
		code := "authentication_failed"
		
		// Adjust status codes for known error types
		if strings.Contains(err.Error(), "validation") {
			status = http.StatusBadRequest
			code = "validation_error"
		} else if strings.Contains(err.Error(), "internal") || strings.Contains(err.Error(), "failed to") {
			status = http.StatusInternalServerError
			code = "internal_error"
		}
		
		c.JSON(status, models.ErrorResponse{
			Error:   code,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
