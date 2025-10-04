package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
		cleanup:  time.Hour, // Clean up old limiters every hour
	}

	// Start cleanup goroutine
	go rl.cleanupRoutine()
	
	return rl
}

// getLimiter gets or creates a limiter for the given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double check in case another goroutine created it
		if limiter, exists = rl.limiters[key]; !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[key] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// Allow checks if the request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	return rl.getLimiter(key).Allow()
}

// cleanupRoutine periodically removes unused limiters
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for key, limiter := range rl.limiters {
			// Remove limiters that haven't been used recently
			if limiter.Tokens() == float64(rl.burst) {
				delete(rl.limiters, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		// Use IP address as the key for rate limiting
		key := c.ClientIP()
		
		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OTPRateLimitMiddleware creates a rate limiting middleware specifically for OTP endpoints
func OTPRateLimitMiddleware() gin.HandlerFunc {
	// Allow 10 requests per minute for OTP endpoints
	return RateLimitMiddleware(10.0/60.0, 10)
}

// LoginRateLimitMiddleware creates a rate limiting middleware specifically for login endpoints
func LoginRateLimitMiddleware() gin.HandlerFunc {
	// Allow 10 login attempts per 15 minutes
	return RateLimitMiddleware(10.0/(15.0*60.0), 10)
}