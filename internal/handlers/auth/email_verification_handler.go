package auth

import (
	"net/http"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// EmailVerificationHandler manages the verify-email endpoint
//
type EmailVerificationHandler struct {
	service *services.EmailVerificationService
}

func NewEmailVerificationHandler(s *services.EmailVerificationService) *EmailVerificationHandler {
	return &EmailVerificationHandler{service: s}
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "token is required"})
		return
	}
	if err := h.service.Verify(c.Request.Context(), req.Token); err != nil {
		status := http.StatusBadRequest
		code := "invalid_or_expired_token"
		c.JSON(status, models.ErrorResponse{Error: code, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.VerifyEmailResponse{Success: true, Message: "Email verified successfully"})
}
