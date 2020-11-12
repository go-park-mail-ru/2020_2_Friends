package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/friends/internal/pkg/csrf"
	"github.com/go-redis/redis/v8"
)

type CSRFRepositoryRedis struct {
	redis *redis.Client
}

func New(redis *redis.Client) (csrf.Repository, error) {
	repo := CSRFRepositoryRedis{
		redis: redis,
	}

	ctx := context.Background()
	err := redis.Ping(ctx).Err()
	if err != nil {
		return CSRFRepositoryRedis{}, fmt.Errorf("redis doesn't not available: %w", err)
	}

	return repo, nil
}

func (c CSRFRepositoryRedis) Add(token string, session string, expires time.Duration) error {
	ctx := context.Background()
	err := c.redis.Set(ctx, token, session, expires).Err()
	if err != nil {
		return fmt.Errorf("couldn't set value csrf token in redis: %w", err)
	}

	return nil
}

func (c CSRFRepositoryRedis) Get(token string) (string, error) {
	ctx := context.Background()
	token, err := c.redis.Get(ctx, token).Result()
	if err != nil {
		return "", fmt.Errorf("couldn't get csrf token from redis: %w", err)
	}

	return token, nil
}
