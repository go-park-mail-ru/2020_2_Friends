package carts

import (
	"fmt"

	"github.com/friends/internal/pkg/models"
)

var ErrCartIsEmpty = fmt.Errorf("cart is empty")

type Repository interface {
	Add(userID, productID, vendorID string) error
	Remove(userID, productID string) error
	Get(userID string) ([]models.Product, error)
	GetVendorIDFromCart(userID string) (int, error)
}
