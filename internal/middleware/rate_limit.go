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

	// Initialize the limiter once (100 capacity, 2 tokens per second)
	limiter = algorithms.NewTokenBucket(100, 2, s)
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.JSON(500, gin.H{"error": "Rate limiter not initialized"})
			c.Abort()
			return
		}

		userIP := c.ClientIP()
		cacheKey := fmt.Sprintf("rate_limit:%s", userIP)

		result, err := limiter.Allow(cacheKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to check rate limit"})
			c.Abort()
			return
		}

		if !result.Allowed {
			c.JSON(429, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
