package order

import "github.com/friends/internal/pkg/models"

type Repository interface {
	AddOrder(userID string, order models.OrderRequest) (int, error)
	GetOrder(orderID string) (models.OrderResponse, error)
	GetUserOrders(userID string) ([]models.OrderResponse, error)
	CheckOrderByUser(userID string, orderID string) bool
	GetVendorOrders(vendorID string) ([]models.OrderResponse, error)
	GetVendorOrdersIDs(vendorID string) ([]int, error)
	UpdateOrderStatus(orderID string, status string) error
	GetProductsFromOrder(order *models.OrderResponse) error
	GetVendorIDFromOrder(orderID int) (int, error)
	SetOrderReviewStatus(orderID int, status bool) error
	GetUserIDFromOrder(orderID int) (string, error)
}
