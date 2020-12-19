package vendors

import (
	"mime/multipart"

	"github.com/friends/internal/pkg/models"
)

//go:generate mockgen -destination=./usecase_mock.go -package=vendors github.com/friends/internal/pkg/vendors Usecase
type Usecase interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	Create(partnerID string, vendor models.Vendor) (int, error)
	Update(vendor models.Vendor) error
	CheckVendorOwner(userID, vendorID string) error
	AddProduct(product models.Product) (int, error)
	UpdateProduct(product models.Product) error
	DeleteProduct(productID string) error
	UpdateVendorPicture(vendorID string, file multipart.File, imgType string) (string, error)
	UpdateProductPicture(productID string, file multipart.File, imgType string) (string, error)
	GetVendorIDFromProduct(productID string) (string, error)
	GetPartnerShops(partnerID string) ([]models.Vendor, error)
	GetVendorOwner(vendorID int) (string, error)
	GetNearest(longitude, latitude float64) ([]models.Vendor, error)
	GetSimilar(vendorID string, longitude, latitude float64) ([]models.Vendor, error)
}
