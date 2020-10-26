package carts

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Add(userID, productID string) error
	Remove(userID, productID string) error
	Get(userID string) ([]models.Product, error)
}
