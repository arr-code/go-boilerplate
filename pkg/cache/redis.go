package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"xnoia-go-boilerplate/internal/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client connection
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("invalid redis DB number: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       db,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// CloseRedis closes the Redis client connection
func CloseRedis(client *redis.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// TODO: Implement session storage functions
// Example functions to implement:
// - SetSession(ctx context.Context, key string, value interface{}, expiration time.Duration) error
// - GetSession(ctx context.Context, key string) (string, error)
// - DeleteSession(ctx context.Context, key string) error
// - SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
// - GetCache(ctx context.Context, key string) (string, error)
// - DeleteCache(ctx context.Context, key string) error
