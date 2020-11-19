package vendors

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	GetAllProductsWithIDsFromSameVendor(ids []int) ([]models.Product, error)
	GetVendorIDFromProduct(productID string) (string, error)
	GetVendorFromProduct(productID int) (models.Vendor, error)
	IsVendorExists(vendorName string) error
	Create(partnerID string, vendor models.Vendor) (int, error)
	Update(models.Vendor) error
	CheckVendorOwner(userID, vendorID string) error
	AddProduct(product models.Product) (int, error)
	UpdateProduct(product models.Product) error
	DeleteProduct(productID string) error
	UpdateVendorImage(vendorID string, link string) error
	UpdateProductImage(vendorID string, link string) error
	GetPartnerShops(partnerID string) ([]models.Vendor, error)
}
