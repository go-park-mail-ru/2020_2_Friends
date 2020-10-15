package user

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Create(user models.User) (userID string, err error)
	CheckIfUserExists(user models.User) bool
	CheckLoginAndPassword(user models.User) (userID string, err error)
}
