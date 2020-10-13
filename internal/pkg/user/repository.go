package user

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Create(user models.User) error
}
