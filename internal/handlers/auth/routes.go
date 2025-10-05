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
	g := r.Group("/google")
	g.POST("/login", handler.GoogleHandleLogin)
	g.POST("/callback", handler.GoogleHandleCallback)
}
