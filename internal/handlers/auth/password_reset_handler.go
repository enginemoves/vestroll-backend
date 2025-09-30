package auth

import (
	"net/http"
	"time"
	"github.com/go-redis/redis/v8"
	"github.com/gin-gonic/gin"
	"github.com/codeZe-us/vestroll-backend/internal/services/email_service"
	"github.com/codeZe-us/vestroll-backend/internal/services/sms_service"
	"github.com/codeZe-us/vestroll-backend/internal/services/password_reset_service"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/utils"
)

// Handler struct for password reset
// Add dependencies as needed
// EmailService, SMSService, RedisClient

type PasswordResetHandler struct {
	EmailService *email_service.EmailService
	SMSService   *sms_service.SMSService
	RedisClient  *redis.Client
}

// POST /api/auth/forgot-password
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Channel    string `json:"channel"` // "email" or "sms"
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Identifier == "" || (req.Channel != "email" && req.Channel != "sms") {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	resetService := getPasswordResetService(h)
	_, err := resetService.GenerateAndSendResetCode(c.Request.Context(), req.Identifier, req.Channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reset code sent"})
}

// POST /api/auth/verify-reset-code
func (h *PasswordResetHandler) VerifyResetCode(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Code       string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Identifier == "" || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	resetService := getPasswordResetService(h)
	ok, err := resetService.VerifyResetCode(c.Request.Context(), req.Identifier, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid or expired code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Code verified"})
}

// POST /api/auth/reset-password
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Identifier  string `json:"identifier"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Identifier == "" || req.Code == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}
	resetService := getPasswordResetService(h)
	ok, err := resetService.VerifyResetCode(c.Request.Context(), req.Identifier, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid or expired code"})
		return
	}
	if err := utils.ValidatePasswordStrength(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// TODO: Update password in user database here
	_ = resetService.DeleteResetCode(c.Request.Context(), req.Identifier)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password reset successful"})
}

// Utility: Generate and cache reset code/token

// Utility: Validate password strength
// Helper to get PasswordResetService
func getPasswordResetService(h *PasswordResetHandler) *password_reset_service.PasswordResetService {
	return password_reset_service.NewPasswordResetService(
		repository.NewPasswordResetRepository(h.RedisClient, 5*time.Minute),
		h.EmailService,
		h.SMSService,
		5*time.Minute,
	)
}
