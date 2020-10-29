package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/cart"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
)

type CartUsecase struct {
	cartsRepository  cart.Repository
	vendorRepository vendors.Repository
}

func NewCartUsecase(cartsRepo cart.Repository, vendorsRepo vendors.Repository) cart.Usecase {
	return CartUsecase{
		cartsRepository:  cartsRepo,
		vendorRepository: vendorsRepo,
	}
}

func (c CartUsecase) Add(userID, productID, vendorID string) error {
	cartVendorID, err := c.cartsRepository.GetVendorIDFromCart(userID)
	if err != nil && !errors.Is(err, cart.ErrCartIsEmpty) {
		return fmt.Errorf("error with db: %w", err)
	}

	if err == nil {
		intVendorID, err := strconv.Atoi(vendorID)
		if err != nil {
			return fmt.Errorf("err with converting vendorID to int: %w", err)
		}

		if intVendorID != cartVendorID {
			return fmt.Errorf("wrong vendor id")
		}
	}

	err = c.cartsRepository.Add(userID, productID, vendorID)
	if err != nil {
		return fmt.Errorf("couldn't add product to cart: %w", err)
	}

	return nil
}

func (c CartUsecase) Remove(userID, productID string) error {
	return c.cartsRepository.Remove(userID, productID)
}

func (c CartUsecase) Get(userID string) ([]models.Product, error) {
	ids, err := c.cartsRepository.GetProductIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get product ids: %w", err)
	}
	return c.vendorRepository.GetAllProductsWithIDs(ids)
}
