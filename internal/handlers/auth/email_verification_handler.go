package auth

import (
	"net/http"
	"regexp"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// EmailVerificationHandler verifies a 6-digit email code and marks the user verified
// It uses the OTPService storage/validation and, on success, updates the user repository.
//
type EmailVerificationHandler struct {
	otpService *services.OTPService
	userRepo   *repository.UserRepository
}

func NewEmailVerificationHandler(otp *services.OTPService, userRepo *repository.UserRepository) *EmailVerificationHandler {
	return &EmailVerificationHandler{otpService: otp, userRepo: userRepo}
}

// VerifyEmail handles POST /api/auth/verify-email with a numeric 6-digit code
// Body: { "user_id": "...", "email": "user@example.com", "code": "123456" }
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
		Code   string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	// Ensure code has only digits
	if !regexp.MustCompile(`^[0-9]{6}$`).MatchString(req.Code) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "code must be 6 digits"})
		return
	}

	// Verify via OTP service
	verifyReq := models.OTPVerificationRequest{Identifier: req.Email, Code: req.Code, Type: models.OTPTypeEmail}
	if err := h.otpService.VerifyOTP(c.Request.Context(), verifyReq); err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "verification_failed"
		msg := err.Error()
		if msg == "OTP not found or expired" || msg == "OTP has expired" {
			errorCode = "otp_expired"
		} else if msg == "invalid OTP code" {
			errorCode = "invalid_otp"
		} else if msg == "maximum verification attempts exceeded" {
			statusCode = http.StatusTooManyRequests
			errorCode = "max_attempts_exceeded"
		} else if len(msg) >= 7 && msg[:7] == "invalid" {
			errorCode = "validation_error"
		}
		c.JSON(statusCode, models.ErrorResponse{Error: errorCode, Message: msg})
		return
	}

	// Mark user as verified
	if err := h.userRepo.SetEmailVerified(c.Request.Context(), req.UserID, req.Email, time.Now()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.VerifyEmailResponse{Success: true, Message: "Email verified successfully"})
}
