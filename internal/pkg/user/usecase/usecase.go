package usecase

import (
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
)

type UserUsecase struct {
	Repository user.Repository
}

func (uu UserUsecase) Create(user models.User) error {
	return uu.Repository.Create(user)
}
