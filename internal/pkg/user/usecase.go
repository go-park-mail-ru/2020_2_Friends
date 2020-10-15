package user

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Create(user models.User) (userID string, err error)
	CheckIfUserExists(user models.User) bool
	Verify(user models.User) (userID string, err error)
}
