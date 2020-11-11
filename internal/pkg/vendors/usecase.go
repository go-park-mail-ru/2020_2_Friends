package vendors

import (
	"mime/multipart"

	"github.com/friends/internal/pkg/models"
)

type Usecase interface {
	Get(id int) (models.Vendor, error)
	GetAll() ([]models.Vendor, error)
	Create(partnerID string, vendor models.Vendor) (int, error)
	Update(vendor models.Vendor) error
	CheckVendorOwner(userID, vendorID string) error
	AddProduct(product models.Product) (int, error)
	UpdateProduct(product models.Product) error
	DeleteProduct(productID int) error
	UpdateVendorPicture(vendorID string, file multipart.File) (string, error)
	UpdateProductPicture(productID string, file multipart.File) (string, error)
	GetVendorIDFromProduct(productID string) (string, error)
}
