package repository

import (
	"context"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/go-redis/redis/v8"
)

type SessionRedisRepo struct {
	redis *redis.Client
}

func NewSessionRedisRepo(redis *redis.Client) (session.Repository, error) {
	repo := SessionRedisRepo{
		redis: redis,
	}

	ctx := context.Background()
	err := redis.Ping(ctx).Err()
	if err != nil {
		return SessionRedisRepo{}, err
	}

	return repo, nil
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
