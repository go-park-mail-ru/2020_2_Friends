package usecase

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
	"github.com/lithammer/shortuuid"
)

type VendorUsecase struct {
	repository vendors.Repository
}

func NewVendorUsecase(repo vendors.Repository) vendors.Usecase {
	return VendorUsecase{
		repository: repo,
	}
}

func (v VendorUsecase) Get(id int) (models.Vendor, error) {
	return v.repository.Get(id)
}

func (v VendorUsecase) GetAll() ([]models.Vendor, error) {
	return v.repository.GetAll()
}

func (v VendorUsecase) Create(partnerID string, vendor models.Vendor) (int, error) {
	err := v.repository.IsVendorExists(vendor.Name)
	if err != nil {
		return 0, err
	}

	return v.repository.Create(partnerID, vendor)
}

func (v VendorUsecase) Update(vendor models.Vendor) error {
	return v.repository.Update(vendor)
}

func (v VendorUsecase) CheckVendorOwner(userID, vendorID string) error {
	return v.repository.CheckVendorOwner(userID, vendorID)
}

func (v VendorUsecase) AddProduct(product models.Product) (int, error) {
	return v.repository.AddProduct(product)
}

func (v VendorUsecase) UpdateProduct(product models.Product) error {
	return v.repository.UpdateProduct(product)
}

func (v VendorUsecase) DeleteProduct(productID string) error {
	return v.repository.DeleteProduct(productID)
}

func (v VendorUsecase) UpdateVendorPicture(vendorID string, file multipart.File) (string, error) {
	img, imgType, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("unsupporter img type: %w", err)
	}

	imgName := shortuuid.New()
	imgFullName := imgName + "." + imgType
	avatarFile, err := os.Create(filepath.Join(configs.FileServerPath+"/img", filepath.Base(imgFullName)))
	if err != nil {
		return "", fmt.Errorf("couldn't create file: %w", err)
	}

	switch imgType {
	case "png":
		err = png.Encode(avatarFile, img)
	case "jpg":
		fallthrough
	case "jpeg":
		err = jpeg.Encode(avatarFile, img, nil)
	default:
		return "", fmt.Errorf("unsupporter img type: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("couldn't encode: %w", err)
	}

	err = v.repository.UpdateVendorImage(vendorID, imgFullName)
	if err != nil {
		return "", fmt.Errorf("couldn't save link vendors picture: %w", err)
	}

	return imgFullName, nil
}

func (v VendorUsecase) UpdateProductPicture(productID string, file multipart.File) (string, error) {
	img, imgType, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("unsupporter img type: %w", err)
	}

	imgName := shortuuid.New()
	imgFullName := imgName + "." + imgType
	avatarFile, err := os.Create(filepath.Join(configs.FileServerPath+"/img", filepath.Base(imgFullName)))
	if err != nil {
		return "", fmt.Errorf("couldn't create file: %w", err)
	}

	switch imgType {
	case "png":
		err = png.Encode(avatarFile, img)
	case "jpg":
		fallthrough
	case "jpeg":
		err = jpeg.Encode(avatarFile, img, nil)
	default:
		return "", fmt.Errorf("unsupporter img type: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("couldn't encode: %w", err)
	}

	err = v.repository.UpdateProductImage(productID, imgFullName)
	if err != nil {
		return "", fmt.Errorf("couldn't save link products picture: %w", err)
	}

	return imgFullName, nil
}

func (v VendorUsecase) GetVendorIDFromProduct(productID string) (string, error) {
	return v.repository.GetVendorIDFromProduct(productID)
}
