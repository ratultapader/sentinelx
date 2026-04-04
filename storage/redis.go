package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	url := os.Getenv("REDIS_URL")

	// Default for Docker local
	if url == "" {
		url = "redis:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: url,
	})

	// 🔥 Try connection
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		// ✅ DO NOT CRASH — make Redis optional
		fmt.Println("⚠️ Redis not available, running without Redis:", err)
		RDB = nil
		return
	}

	fmt.Println("✅ Redis connected")
}