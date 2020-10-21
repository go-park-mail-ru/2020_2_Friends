package vendors

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Get(id int) (models.Vendor, error)
}
