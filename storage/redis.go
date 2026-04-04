package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() error {
	url := os.Getenv("REDIS_URL")

	if url == "" {
		url = "redis:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: url,
	})

	// test connection
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	fmt.Println("✅ Redis connected")
	return nil
}