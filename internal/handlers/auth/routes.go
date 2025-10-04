package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterPasswordResetRoutes(r *gin.RouterGroup, handler *PasswordResetHandler) {
	r.POST("/forgot-password", handler.ForgotPassword)
	r.POST("/verify-reset-code", handler.VerifyResetCode)
	r.POST("/reset-password", handler.ResetPassword)
}

func RegisterLoginRoutes(r *gin.RouterGroup, handler *LoginHandler) {
	r.POST("/login", handler.Login)
}
