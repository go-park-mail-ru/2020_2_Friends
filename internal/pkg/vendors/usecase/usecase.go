package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/friends/internal/pkg/fileserver"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc/metadata"
)

type VendorUsecase struct {
	repository vendors.Repository
	fsClient   fileserver.UploadServiceClient
}

func NewVendorUsecase(repo vendors.Repository, fileserverClient fileserver.UploadServiceClient) vendors.Usecase {
	return VendorUsecase{
		repository: repo,
		fsClient:   fileserverClient,
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

func (v VendorUsecase) UpdatePicture(file multipart.File, imageType string) (string, error) {
	imgName := shortuuid.New()
	imgFullName := imgName + "." + imageType

	md := metadata.New(map[string]string{"fileName": imgFullName})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := v.fsClient.Upload(ctx)
	if err != nil {
		return "", err
	}

	chunk := make([]byte, 1024)
	for {
		size, err := file.Read(chunk)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		err = stream.Send(&fileserver.Chunk{Content: chunk[:size]})
		if err != nil {
			return "", err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return imgFullName, nil
}

func (v VendorUsecase) UpdateVendorPicture(vendorID string, file multipart.File, imgType string) (string, error) {
	imgName, err := v.UpdatePicture(file, imgType)
	if err != nil {
		return "", err
	}

	err = v.repository.UpdateVendorImage(vendorID, imgName)
	if err != nil {
		return "", fmt.Errorf("couldn't save link vendors picture: %w", err)
	}

	return imgName, nil
}

func (v VendorUsecase) UpdateProductPicture(productID string, file multipart.File, imgType string) (string, error) {
	imgName, err := v.UpdatePicture(file, imgType)
	if err != nil {
		return "", err
	}

	err = v.repository.UpdateProductImage(productID, imgName)
	if err != nil {
		return "", fmt.Errorf("couldn't save link products picture: %w", err)
	}

	return imgName, nil
}

func (v VendorUsecase) GetVendorIDFromProduct(productID string) (string, error) {
	return v.repository.GetVendorIDFromProduct(productID)
}

func (v VendorUsecase) GetPartnerShops(partnerID string) ([]models.Vendor, error) {
	return v.repository.GetPartnerShops(partnerID)
}

func (v VendorUsecase) GetVendorOwner(vendorID int) (string, error) {
	return v.repository.GetVendorOwner(vendorID)
}

func (v VendorUsecase) GetNearest(longitude, latitude float64) ([]models.Vendor, error) {
	return v.repository.GetNearest(longitude, latitude)
}
