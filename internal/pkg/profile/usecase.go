package profile

import (
	"mime/multipart"

	"github.com/friends/internal/pkg/models"
)

//go:generate mockgen -destination=./usecase_mock.go -package=profile github.com/friends/internal/pkg/profile Usecase
type Usecase interface {
	Create(userID string) error
	Get(userID string) (models.Profile, error)
	Update(models.Profile) error
	UpdateAvatar(userID string, file multipart.File) (string, error)
	UpdateAddresses(userID string, addresses []string) error
	Delete(userID string) error
}
