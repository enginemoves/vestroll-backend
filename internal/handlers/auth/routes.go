package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterPasswordResetRoutes(r *gin.RouterGroup, handler *PasswordResetHandler) {
	r.POST("/forgot-password", handler.ForgotPassword)
	r.POST("/verify-reset-code", handler.VerifyResetCode)
	r.POST("/reset-password", handler.ResetPassword)
}

func RegisterGoogleAuthRoutes(r *gin.RouterGroup, handler *GoogleOAuth) {
	r.POST("/google-login", handler.LoginURL)
	r.POST("/google-callback", handler.HandleCallbackGin)
}
