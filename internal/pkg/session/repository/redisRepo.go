package repository

import (
	"context"

	"github.com/friends/internal/pkg/models"
	"github.com/go-redis/redis/v8"
)

type SessionRedisRepo struct {
	redis *redis.Client
}

func NewSessionRedisRepo(redis *redis.Client) SessionRedisRepo {
	return SessionRedisRepo{
		redis: redis,
	}
}

func (srr SessionRedisRepo) Create(session models.Session) error {
	ctx := context.Background()
	err := srr.redis.Set(ctx, session.Name, session.UserID, session.ExpireTime).Err()

	return err
}

func (srr SessionRedisRepo) Check(sessionName string) (userID string, err error) {
	ctx := context.Background()
	userID, err = srr.redis.Get(ctx, sessionName).Result()

	return userID, err
}

func (srr SessionRedisRepo) Delete(sessionName string) error {
	ctx := context.Background()
	_, err := srr.redis.Do(ctx, "DEL", sessionName).Result()

	return err
}
