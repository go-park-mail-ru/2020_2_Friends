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

func (uu UserUsecase) Create(user models.User) error {
	return uu.repository.Create(user)
}

func (uu UserUsecase) Verify(user models.User) (userID string, err error) {
	return uu.repository.CheckLoginAndPassword(user)
}
