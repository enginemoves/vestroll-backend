package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "VestRoll API is running",
			"version": "1.0.0",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
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

	fmt.Println("🚀 VestRoll Backend starting on :8080")
	fmt.Println("📊 Health check: http://localhost:8080/health")
	fmt.Println("📘 API Base: http://localhost:8080/api/v1")

	log.Fatal(r.Run(":8080"))
}
