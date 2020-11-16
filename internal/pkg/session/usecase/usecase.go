package usecase

import (
	"fmt"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/lithammer/shortuuid/v3"
)

type SessionUsecase struct {
	repository session.Repository
}

func NewSessionUsecase(repo session.Repository) session.Usecase {
	return SessionUsecase{
		repository: repo,
	}
}

func (su SessionUsecase) Create(userID string) (string, error) {
	sessionName := "session:" + shortuuid.New()
	session := models.Session{
		Name:       sessionName,
		UserID:     userID,
		ExpireTime: configs.ExpireTime,
	}

	err := su.repository.Create(session)
	if err != nil {
		return "", fmt.Errorf("session not created: %w", err)
	}

	return sessionName, nil
}

func (su SessionUsecase) Check(sessionName string) (userID string, err error) {
	return su.repository.Check(sessionName)
}

func (su SessionUsecase) Delete(sessionName string) error {
	return su.repository.Delete(sessionName)
}
