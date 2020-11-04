package vendors

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	Create(models.Vendor) error
}
