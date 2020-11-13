package order

import "github.com/friends/internal/pkg/models"

type Repository interface {
	AddOrder(userID string, order models.OrderRequest) (int, error)
	GetOrder(orderID string) (models.OrderResponse, []string, error)
	CheckOrderByUser(userID string, orderID string) bool
}
