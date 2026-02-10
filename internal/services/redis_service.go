package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis successfully")
}
func SetCache(key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, json, expiration).Err()
}
func GetCache(key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}
