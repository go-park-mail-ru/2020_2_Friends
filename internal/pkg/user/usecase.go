package user

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Create(user models.User) error
	Verify(user models.User) (userID string, err error)
}