package user

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./usecase_mock.go -package=user github.com/friends/internal/pkg/user Usecase
type Usecase interface {
	Create(user models.User) (userID string, err error)
	CheckIfUserExists(user models.User) error
	Verify(user models.User) (userID string, err error)
	Delete(userID string) error
}
