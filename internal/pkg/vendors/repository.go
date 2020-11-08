package vendors

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	GetAllProductsWithIDs(ids []string) ([]models.Product, error)
	GetVendorIDFromProduct(productID string) (string, error)
	IsVendorExists(vendorName string) error
	Create(models.Vendor) error
	Update(models.Vendor) error
	CheckVendorOwner(userID, vendorID string) error
	AddProduct(product models.Product) error
	UpdateProduct(product models.Product) error
	DeleteProduct(productID int) error
	UpdateVendorImage(vendorID string, link string) error
	UpdateProductImage(vendorID string, link string) error
}
