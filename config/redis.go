package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// InitRedis to initialize Redis client
func InitRedis(cfg Config, ctx context.Context) *redis.Client {
	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisUri,
	})

	_, err := client.Ping(ctx).Result()

	// Test the connection
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis successfully connected.")

	return client
}
