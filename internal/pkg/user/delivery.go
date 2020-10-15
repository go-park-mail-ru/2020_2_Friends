package user

import "github.com/friends/internal/pkg/models"

type Delivery interface {
	Verify(user models.User) (userID string, err error)
}
