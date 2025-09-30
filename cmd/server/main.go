package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/database"
	"github.com/codeZe-us/vestroll-backend/internal/handlers"
	handlers_auth "github.com/codeZe-us/vestroll-backend/internal/handlers/auth"
	"github.com/codeZe-us/vestroll-backend/internal/middleware"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services/sms_service"
	"github.com/codeZe-us/vestroll-backend/internal/services/email_service"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	
	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "VestRoll API is running",
			"version": "1.0.0",
		})
	})

	// Initialize Redis client
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. OTP functionality will not be available.", err)
	}

	// Initialize services only if Redis is available
	var otpHandler *handlers.OTPHandler
	if redisClient != nil {
		// Initialize repositories
		otpRepo := repository.NewOTPRepository(redisClient, cfg.OTP.TTL)

		// Initialize services
	smsService := sms_service.NewSMSService(cfg.Twilio)
	emailService := email_service.NewEmailService(cfg.SMTP)
	otpService := services.NewOTPService(otpRepo, smsService, emailService, cfg.OTP)

		// Initialize handlers
		otpHandler = handlers.NewOTPHandler(otpService)
	}

		// API routes
		v1 := r.Group("/api/v1")
		{
			auth := v1.Group("/auth")
			{
				// Apply rate limiting to OTP endpoints
				auth.Use(middleware.OTPRateLimitMiddleware())
                
				// OTP endpoints (only if Redis is available)
				if otpHandler != nil {
					otpHandler.RegisterRoutes(auth)
				}

				// Password reset endpoints (only if Redis is available)
				if redisClient != nil {
					emailService := email_service.NewEmailService(cfg.SMTP)
					smsService := sms_service.NewSMSService(cfg.Twilio)
					passwordResetHandler := &handlers_auth.PasswordResetHandler{
						EmailService: emailService,
						SMSService: smsService,
						RedisClient: redisClient,
					}
					handlers_auth.RegisterPasswordResetRoutes(auth, passwordResetHandler)
				}

				// Existing auth endpoints
				auth.POST("/login", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - Coming soon"})
				})
				auth.POST("/register", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - Coming soon"})
				})
			}

			employees := v1.Group("/employees")
			{
				employees.GET("/", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Get employees - Coming soon"})
				})
			}

			payroll := v1.Group("/payroll")
			{
				payroll.GET("/", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Payroll management - Coming soon"})
				})
			}
		}

	fmt.Println(" VestRoll Backend starting on :8080")
	fmt.Println(" Health check: http://localhost:8080/health")
	fmt.Println(" API Base: http://localhost:8080/api/v1")
	
	if redisClient != nil {
		fmt.Println(" OTP Endpoints:")
		fmt.Println("   POST /api/v1/auth/send-otp")
		fmt.Println("   POST /api/v1/auth/verify-otp")
	} else {
		fmt.Println(" OTP endpoints disabled (Redis not available)")
	}

	log.Fatal(r.Run(":8080"))
}
