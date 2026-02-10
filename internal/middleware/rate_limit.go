package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codetesla51/limitz/algorithms"
	"github.com/codetesla51/limitz/store"
	"github.com/gin-gonic/gin"
)

var (
	ipLimiter   algorithms.RateLimiter
	userLimiter algorithms.RateLimiter
)

func InitRateLimiter() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Initialize the store once
	s, err := store.NewRedisStore(redisHost+":"+redisPort, "", redisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to rate limit store: %v", err)
	}

	ipLimiter = algorithms.NewTokenBucket(100, 10, s)

	userLimiter = algorithms.NewSlidingWindowCounter(1000, 1*time.Hour, s)
}

func RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
		if !checkLimit(c, ipLimiter, key) {
			return
		}
		c.Next()
	}
}

func RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		var key string
		var limiter algorithms.RateLimiter

		if !exists {
			key = fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
			limiter = ipLimiter
		} else {
			key = fmt.Sprintf("rate_limit:user:%v", userID)
			limiter = userLimiter
		}

		if !checkLimit(c, limiter, key) {
			return
		}
		c.Next()
	}
}

func checkLimit(c *gin.Context, limiter algorithms.RateLimiter, key string) bool {
	if limiter == nil {
		c.JSON(500, gin.H{"error": "Rate limiter not initialized"})
		c.Abort()
		return false
	}

	result, err := limiter.Allow(key)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check rate limit"})
		c.Abort()
		return false
	}

	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", result.Limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))

	if !result.Allowed {
		c.Header("Retry-After", fmt.Sprintf("%d", int(result.RetryAfter.Seconds())))
		c.JSON(429, gin.H{
			"error":       "Too Many Requests",
			"retry_after": result.RetryAfter.String(),
		})
		c.Abort()
		return false
	}
	return true
}
