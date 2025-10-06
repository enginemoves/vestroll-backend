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

func RegisterAppleAuthRoutes(r *gin.RouterGroup, handler *AppleOAuth) {
	a := r.Group("/apple")
	a.POST("/login", handler.AppleHandleLogin)
	a.POST("/callback", handler.AppleHandleCallback)
func RegisterLoginRoutes(r *gin.RouterGroup, handler *LoginHandler) {
	r.POST("/login", handler.Login)
}
