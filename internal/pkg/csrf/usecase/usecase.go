package usecase

import (
	"fmt"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/csrf"
	"github.com/lithammer/shortuuid/v3"
)

type CSRFUsecase struct {
	repository csrf.Repository
}

func New(repo csrf.Repository) CSRFUsecase {
	return CSRFUsecase{
		repository: repo,
	}
}

func (c CSRFUsecase) Add(session string) (string, error) {
	token := shortuuid.New()
	err := c.repository.Add(token, session, configs.CSRFTokenExpireTime)
	if err != nil {
		return "", fmt.Errorf("error with db, token not added: %w", err)
	}

	return token, nil
}

func (c CSRFUsecase) Check(token string, session string) bool {
	userSession, err := c.repository.Get(token)
	if err != nil {
		return false
	}

	return userSession == session
}
