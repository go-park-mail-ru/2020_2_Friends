package cart

import (
	"fmt"
)

var ErrCartIsEmpty = fmt.Errorf("cart is empty")

type Repository interface {
	Add(userID, productID, vendorID string) error
	Remove(userID, productID string) error
	GetProductIDs(userID string) ([]string, error)
	GetVendorIDFromCart(userID string) (string, error)
}
