package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title VestRoll Payroll API
// @version 1.0
// @description Enterprise-grade payroll management system
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Set Gin mode
	gin.SetMode(gin.DebugMode)

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "VestRoll API is running",
			"version": "1.0.0",
		})
	})

	// API version 1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (placeholder)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - Coming soon"})
			})
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - Coming soon"})
			})
		}

		// Employee routes (placeholder)
		employees := v1.Group("/employees")
		{
			employees.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Get employees - Coming soon"})
			})
		}

		// Payroll routes (placeholder)
		payroll := v1.Group("/payroll")
		{
			payroll.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Payroll management - Coming soon"})
			})
		}
	}

	// Start server
	fmt.Println("🚀 VestRoll Backend starting on :8080")
	fmt.Println("📊 Health check: http://localhost:8080/health")
	fmt.Println("📘 API Base: http://localhost:8080/api/v1")

	log.Fatal(r.Run(":8080"))
}
