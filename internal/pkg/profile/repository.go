package profile

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Create(userID string) error
	Get(userID string) (models.Profile, error)
	Update(models.Profile) error
	UpdateAvatar(userID string, link string) error
	Delete(userID string) error
}
