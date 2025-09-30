package handlers

import (
	"net/http"
	"strings"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// PINHandler manages PIN setup and authentication endpoints
type PINHandler struct {
	service *services.PINService
}

func NewPINHandler(service *services.PINService) *PINHandler {
	return &PINHandler{service: service}
}

// RegisterRoutes registers PIN endpoints under /auth
func (h *PINHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/setup-pin", h.SetupPIN)
	router.POST("/login-pin", h.LoginPIN)
}

// SetupPIN handles POST /api/auth/setup-pin
func (h *PINHandler) SetupPIN(c *gin.Context) {
	var req models.SetupPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if err := h.service.SetupPIN(c.Request.Context(), req); err != nil {
		status := http.StatusBadRequest
		code := "validation_error"
		// Adjust status codes for known error types
		if strings.Contains(err.Error(), "internal") {
			status = http.StatusInternalServerError
			code = "internal_error"
		}
		c.JSON(status, models.ErrorResponse{Error: code, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SetupPINResponse{Success: true, Message: "PIN setup successful"})
}

// LoginPIN handles POST /api/auth/login-pin
func (h *PINHandler) LoginPIN(c *gin.Context) {
	var req models.LoginPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if err := h.service.LoginPIN(c.Request.Context(), req); err != nil {
		status := http.StatusUnauthorized
		code := "invalid_pin"
		if err.Error() == "pin not set" {
			status = http.StatusNotFound
			code = "not_found"
		} else if strings.HasPrefix(err.Error(), "user_") || strings.HasPrefix(err.Error(), "pin") || strings.HasPrefix(err.Error(), "invalid pin") {
			status = http.StatusBadRequest
			code = "validation_error"
		} else if strings.Contains(err.Error(), "internal") {
			status = http.StatusInternalServerError
			code = "internal_error"
		}
		c.JSON(status, models.ErrorResponse{Error: code, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.LoginPINResponse{Success: true, Message: "PIN authentication successful"})
}