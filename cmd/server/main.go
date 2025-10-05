package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codeZe-us/vestroll-backend/intern			// Apply specific rate limiting for login
			loginGroup := auth.Group("/")
			loginGroup.Use(middleware.LoginRateLimitMiddleware())
			auth.RegisterLoginRoutes(loginGroup, loginHandler)onfig"
	"github.com/codeZe-us/vestroll-backend/internal/database"
	"github.com/codeZe-us/vestroll-backend/internal/handlers"
	"github.com/codeZe-us/vestroll-backend/internal/handlers/auth"
	"github.com/codeZe-us/vestroll-backend/internal/middleware"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
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

	// Initialize PostgreSQL client
	postgresClient, err := database.NewPostgresClient(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresClient.Close()

	// Initialize Redis client
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. OTP functionality will not be available.", err)
	}

	// Initialize user services (always available with PostgreSQL)
	userRepo := repository.NewUserRepository(postgresClient)
	userService := services.NewUserService(userRepo, cfg.JWT)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize PostgreSQL connection
	db, err := database.NewPostgreSQLConnection(cfg.Database)
	if err != nil {
		log.Printf("Warning: PostgreSQL connection failed: %v. Login functionality will not be available.", err)
	}


	// Initialize services only if Redis is available
	var otpHandler *handlers.OTPHandler
	var businessProfileHandler *handlers.BusinessProfileHandler
	var pinHandler *handlers.PINHandler
	if redisClient != nil {
		// Initialize repositories
		otpRepo := repository.NewOTPRepository(redisClient, cfg.OTP.TTL)
		businessRepo := repository.NewBusinessProfileRepository(redisClient, 0)
		pinRepo := repository.NewPinRepository(redisClient, 0)

		// Initialize services
		smsService := services.NewSMSService(cfg.Twilio)
		emailService := services.NewEmailService(cfg.SMTP)
		otpService := services.NewOTPService(otpRepo, smsService, emailService, cfg.OTP)
		businessService := services.NewBusinessProfileService(businessRepo)
		pinService := services.NewPINService(pinRepo)

		// Initialize handlers
		otpHandler = handlers.NewOTPHandler(otpService)
		businessProfileHandler = handlers.NewBusinessProfileHandler(businessService)
		pinHandler = handlers.NewPINHandler(pinService)
	}

	// Initialize authentication services (if PostgreSQL is available)
	var loginHandler *auth.LoginHandler
	if db != nil {
		// Initialize user repository and auth services
		userRepo := repository.NewUserRepository(db)
		jwtService := services.NewJWTService(cfg.JWT)
		authService := services.NewAuthService(userRepo, jwtService)

		// Initialize login handler
		loginHandler = auth.NewLoginHandler(authService)
	}

	// API routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			// Apply rate limiting to OTP endpoints
			auth.Use(middleware.OTPRateLimitMiddleware())
			

			// User authentication endpoints (always available)
			userHandler.RegisterRoutes(auth)

			// Login endpoints (only if PostgreSQL is available)
			if loginHandler != nil {
				// Apply specific rate limiting for login
				loginGroup := auth.Group("/")
				loginGroup.Use(middleware.LoginRateLimitMiddleware())
				auth_routes.RegisterLoginRoutes(loginGroup, loginHandler)
			}

			
			// OTP endpoints (only if Redis is available)
			if otpHandler != nil {
				otpHandler.RegisterRoutes(auth)
			}

			// PIN endpoints (only if Redis is available)
			if pinHandler != nil {
				pinHandler.RegisterRoutes(auth)
			}


			// Register endpoint (placeholder for future implementation)
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

		profile := v1.Group("/profile")
		{
			if businessProfileHandler != nil {
				businessProfileHandler.RegisterRoutes(profile)
			}
		}
	}

	// Also expose non-versioned /api/profile routes for compatibility
	api := r.Group("/api")
	{
		profile := api.Group("/profile")
		{
			if businessProfileHandler != nil {
				businessProfileHandler.RegisterRoutes(profile)
			}
		}
	}

	fmt.Println(" VestRoll Backend starting on :8080")
	fmt.Println(" Health check: http://localhost:8080/health")
	fmt.Println(" API Base: http://localhost:8080/api/v1")
	fmt.Println(" Auth Endpoints:")
	fmt.Println("   POST /api/v1/auth/register")
	fmt.Println("   POST /api/v1/auth/login")
	
	if db != nil {
		fmt.Println(" Authentication Endpoints:")
		fmt.Println("   POST /api/v1/auth/login")
	} else {
		fmt.Println(" Login endpoint disabled (PostgreSQL not available)")
	}

	if redisClient != nil {
		fmt.Println(" OTP Endpoints:")
		fmt.Println("   POST /api/v1/auth/send-otp")
		fmt.Println("   POST /api/v1/auth/verify-otp")
		fmt.Println(" Profile Endpoints:")
		fmt.Println("   POST /api/v1/profile/business-details")
		fmt.Println("   POST /api/profile/business-details")
		fmt.Println(" PIN Endpoints:")
		fmt.Println("   POST /api/v1/auth/setup-pin")
		fmt.Println("   POST /api/v1/auth/login-pin")
	} else {
		fmt.Println(" OTP endpoints disabled (Redis not available)")
		fmt.Println(" Profile endpoints disabled (Redis not available)")
	}

	log.Fatal(r.Run(":8080"))
}
