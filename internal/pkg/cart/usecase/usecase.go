package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/cart"
	"github.com/friends/internal/pkg/models"
)

type CartUsecase struct {
	repository cart.Repository
}

func NewCartUsecase(repo cart.Repository) cart.Usecase {
	return CartUsecase{
		repository: repo,
	}
}

func (c CartUsecase) Add(userID, productID, vendorID string) error {
	cartVendorID, err := c.repository.GetVendorIDFromCart(userID)
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

	err = c.repository.Add(userID, productID, vendorID)
	if err != nil {
		return fmt.Errorf("couldn't add product to cart: %w", err)
	}

	return nil
}

func (c CartUsecase) Remove(userID, productID string) error {
	return c.repository.Remove(userID, productID)
}

func (c CartUsecase) Get(userID string) ([]models.Product, error) {
	return c.repository.Get(userID)
}
