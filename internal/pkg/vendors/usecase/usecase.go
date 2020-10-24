package usecase

import (
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
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
