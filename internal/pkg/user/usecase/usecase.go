package usecase

import (
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
)

type UserUsecase struct {
	repository user.Repository
}

func NewUserUsecase(repo user.Repository) user.Usecase {
	return UserUsecase{
		repository: repo,
	}
}

func (u UserUsecase) Create(user models.User) (userID string, err error) {
	return u.repository.Create(user)
}

func (u UserUsecase) CheckIfUserExists(user models.User) error {
	return u.repository.CheckIfUserExists(user)
}

func (u UserUsecase) Verify(user models.User) (userID string, err error) {
	return u.repository.CheckLoginAndPassword(user)
}

func (u UserUsecase) Delete(userID string) error {
	return u.repository.Delete(userID)
}

func (u UserUsecase) CheckUsersRole(userID string) (int, error) {
	return u.repository.CheckUsersRole(userID)
}
