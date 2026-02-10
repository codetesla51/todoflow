package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/codetesla51/limitz/algorithms"
	"github.com/codetesla51/limitz/store"
	"github.com/gin-gonic/gin"
)

var limiter *algorithms.TokenBucket

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

	// Initialize the limiter once (100 capacity, 10 tokens per second)
	limiter = algorithms.NewTokenBucket(100, 10, s)
}

// RateLimitByIP is good for /auth routes or general server protection
func RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
		if !checkLimit(c, key) {
			return
		}
		c.Next()
	}
}

func RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// Fallback to IP if for some reason user_id isn't there
			key := fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
			if !checkLimit(c, key) {
				return
			}
		} else {
			key := fmt.Sprintf("rate_limit:user:%v", userID)
			if !checkLimit(c, key) {
				return
			}
		}
		c.Next()
	}
}

func checkLimit(c *gin.Context, key string) bool {
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

	if !result.Allowed {
		c.JSON(429, gin.H{"error": "Too Many Requests"})
		c.Abort()
		return false
	}
	return true
}
