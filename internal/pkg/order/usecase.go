package order

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./usecase_mock.go -package=order github.com/friends/internal/pkg/order Usecase
type Usecase interface {
	AddOrder(userID string, order models.OrderRequest) (int, error)
	GetOrder(userID string, orderID string) (models.OrderResponse, error)
	GetUserOrders(userID string) ([]models.OrderResponse, error)
	GetVendorOrders(vendorID string) (models.VendorOrdersResponse, error)
	UpdateOrderStatus(orderID string, status string) error
	GetVendorIDFromOrder(orderID int) (int, error)
	GetUserIDFromOrder(orderID int) (string, error)
}
