package user

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Create(user models.User) (userID string, err error)
	CheckIfUserExists(user models.User) error
	Verify(user models.User) (userID string, err error)
	Delete(userID string) error
}
