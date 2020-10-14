package usecase

import (
	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
)

type SessionUsecase struct {
	repository session.Repository
}

func NewSessionUsecase(repo session.Repository) SessionUsecase {
	return SessionUsecase{
		repository: repo,
	}
}

func (su SessionUsecase) Create(userID string) (string, error) {
	sessionName := "session:" + userID
	session := models.Session{
		Name:       sessionName,
		UserID:     userID,
		ExpireTime: configs.ExpireTime,
	}

	err := su.repository.Create(session)
	if err != nil {
		return "", err
	}

	return sessionName, nil
}
