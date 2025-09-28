package handlers

import (
	"net/http"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	otpService *services.OTPService
}

func NewOTPHandler(otpService *services.OTPService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

// SendOTP handles POST /api/auth/send-otp
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req models.OTPRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.otpService.SendOTP(c.Request.Context(), req); err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"
		
		// Handle specific error cases
		if err.Error() == "rate limit exceeded" || 
		   err.Error()[:17] == "rate limit exceeded" {
			statusCode = http.StatusTooManyRequests
			errorCode = "rate_limit_exceeded"
		} else if err.Error()[:7] == "invalid" {
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		} else if err.Error()[:13] == "SMS service is" || 
				  err.Error()[:15] == "email service is" {
			statusCode = http.StatusServiceUnavailable
			errorCode = "service_unavailable"
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.OTPResponse{
		Success: true,
		Message: "OTP sent successfully",
	})
}

// VerifyOTP handles POST /api/auth/verify-otp
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req models.OTPVerificationRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.otpService.VerifyOTP(c.Request.Context(), req); err != nil {
		// Determine status code based on error type
		statusCode := http.StatusBadRequest
		errorCode := "verification_failed"
		
		// Handle specific error cases
		if err.Error() == "OTP not found or expired" ||
		   err.Error() == "OTP has expired" {
			errorCode = "otp_expired"
		} else if err.Error() == "invalid OTP code" {
			errorCode = "invalid_otp"
		} else if err.Error() == "maximum verification attempts exceeded" {
			statusCode = http.StatusTooManyRequests
			errorCode = "max_attempts_exceeded"
		} else if err.Error()[:7] == "invalid" {
			errorCode = "validation_error"
		}

		c.JSON(statusCode, models.ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.OTPVerificationResponse{
		Success: true,
		Message: "OTP verified successfully",
	})
}

// RegisterRoutes registers the OTP routes
func (h *OTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/send-otp", h.SendOTP)
	router.POST("/verify-otp", h.VerifyOTP)
}