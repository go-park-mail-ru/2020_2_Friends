package vendors

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	GetAllProductsWithIDs(ids []string) ([]models.Product, error)
	GetVendorIDFromProduct(productID string) (string, error)
}
