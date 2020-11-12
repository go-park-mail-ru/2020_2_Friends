package user

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./repo_mock.go -package=user github.com/friends/internal/pkg/user Repository
type Repository interface {
	Create(user models.User) (userID string, err error)
	CheckIfUserExists(user models.User) error
	CheckLoginAndPassword(user models.User) (userID string, err error)
	Delete(userID string) error
	CheckUsersRole(userID string) (int, error)
}
