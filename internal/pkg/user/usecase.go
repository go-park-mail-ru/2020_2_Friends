package user

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Create(user models.User) error
}
