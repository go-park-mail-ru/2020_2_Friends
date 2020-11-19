package usecase

import (
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
	ownErr "github.com/friends/pkg/error"
)

type OrderUsecase struct {
	orderRepository  order.Repository
	vendorRepository vendors.Repository
}

func New(orderRepository order.Repository, vendorRepository vendors.Repository) order.Usecase {
	return OrderUsecase{
		orderRepository:  orderRepository,
		vendorRepository: vendorRepository,
	}
}

func (o OrderUsecase) AddOrder(userID string, order models.OrderRequest) (int, error) {
	vendor, err := o.vendorRepository.GetVendorFromProduct(order.ProductIDs[0])
	if err != nil {
		return 0, fmt.Errorf("error with db: %w", err)
	}

	order.VendorID = vendor.ID
	order.VendorName = vendor.Name

	products, err := o.vendorRepository.GetAllProductsWithIDsFromSameVendor(order.ProductIDs)
	if err != nil {
		return 0, err
	}

	for _, product := range products {
		order.Products = append(order.Products, models.OrderProduct{
			Name:    product.Name,
			Price:   product.Price,
			Picture: product.Picture,
		})

		order.Price += product.Price
	}

	return o.orderRepository.AddOrder(userID, order)
}

func (o OrderUsecase) GetOrder(userID string, orderID string) (models.OrderResponse, error) {
	if !o.orderRepository.CheckOrderByUser(userID, orderID) {
		return models.OrderResponse{}, ownErr.NewClientError(fmt.Errorf("user is not order owner"))
	}

	order, err := o.orderRepository.GetOrder(orderID)
	if err != nil {
		return models.OrderResponse{}, err
	}

	return order, nil
}

func (o OrderUsecase) GetUserOrders(userID string) ([]models.OrderResponse, error) {
	orders, err := o.orderRepository.GetUserOrders(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o OrderUsecase) GetVendorOrders(vendorID string) ([]models.OrderResponse, error) {
	orders, err := o.orderRepository.GetVendorOrders(vendorID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o OrderUsecase) UpdateOrderStatus(orderID string, status string) error {
	return o.orderRepository.UpdateOrderStatus(orderID, status)
}
