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
	ipLimiter     algorithms.RateLimiter
	userLimiter   algorithms.RateLimiter
	dbIpLimiter   algorithms.RateLimiter
	dbUserLimiter algorithms.RateLimiter
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

	// Initialize Redis store
	redisStore, err := store.NewRedisStore(redisHost+":"+redisPort, "", redisPassword)
	if err == nil {
		ipLimiter = algorithms.NewTokenBucket(100, 10, redisStore)
		userLimiter = algorithms.NewSlidingWindowCounter(1000, 1*time.Hour, redisStore)
		log.Println("Redis rate limiter initialized")
	} else {
		log.Printf("Failed to connect to Redis for rate limiting: %v", err)
	}

	// Initialize Database store as fallback
	dbAddr := os.Getenv("DATABASE_URL")
	if dbAddr != "" {
		dbStore, err := store.NewDatabaseStore(dbAddr)
		if err == nil {
			dbIpLimiter = algorithms.NewTokenBucket(100, 10, dbStore)
			dbUserLimiter = algorithms.NewSlidingWindowCounter(1000, 1*time.Hour, dbStore)
			log.Println("Database fallback rate limiter initialized")

			go func() {
				ticker := time.NewTicker(10 * time.Minute)
				for range ticker.C {
					if err := dbStore.CleanupExpired(); err != nil {
						log.Printf("DB rate limit cleanup error: %v", err)
					}
				}
			}()
		} else {
			log.Printf("Failed to connect to Database for rate limiting: %v", err)
		}
	}
}

func RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
		if !checkLimit(c, ipLimiter, dbIpLimiter, key) {
			return
		}
		c.Next()
	}
}

func RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		var key string
		var primary, fallback algorithms.RateLimiter

		if !exists {
			key = fmt.Sprintf("rate_limit:ip:%s", c.ClientIP())
			primary = ipLimiter
			fallback = dbIpLimiter
		} else {
			key = fmt.Sprintf("rate_limit:user:%v", userID)
			primary = userLimiter
			fallback = dbUserLimiter
		}

		if !checkLimit(c, primary, fallback, key) {
			return
		}
		c.Next()
	}
}

func checkLimit(c *gin.Context, primary algorithms.RateLimiter, fallback algorithms.RateLimiter, key string) bool {
	var result algorithms.Result
	var err error

	// Try primary (Redis)
	if primary != nil {
		result, err = primary.Allow(key)
		if err == nil {
			return handleResult(c, result)
		}
		log.Printf("Primary rate limiter error: %v, attempting fallback", err)
	}

	// Fallback to DB
	if fallback != nil {
		result, err = fallback.Allow(key)
		if err == nil {
			return handleResult(c, result)
		}
		log.Printf("Fallback rate limiter error: %v", err)
	}

	// Both failed, fail-open
	return true
}

func handleResult(c *gin.Context, result algorithms.Result) bool {
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
