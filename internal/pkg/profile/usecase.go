package profile

import (
	"mime/multipart"

	"github.com/friends/internal/pkg/models"
)

type Usecase interface {
	Create(userID string) error
	Get(userID string) (models.Profile, error)
	Update(models.Profile) error
	UpdateAvatar(userID string, file multipart.File) (string, error)
	Delete(userID string) error
}
