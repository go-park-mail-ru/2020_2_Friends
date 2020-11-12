package cart

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./usecase_mock.go -package=cart github.com/friends/internal/pkg/cart Usecase
type Usecase interface {
	Add(userID, productID string) error
	Remove(userID, productID string) error
	Get(userID string) ([]models.Product, error)
}
