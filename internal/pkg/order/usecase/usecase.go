package usecase

import (
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
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
	vendor, err := o.vendorRepository.GetVendorFromProduct(order.Products[0])
	if err != nil {
		return 0, fmt.Errorf("error with db: %w", err)
	}

	order.VendorID = vendor.ID
	order.VendorName = vendor.Name

	return o.orderRepository.AddOrder(userID, order)
}

func (o OrderUsecase) GetOrder(userID string, orderID string) (models.OrderResponse, error) {
	if !o.orderRepository.CheckOrderByUser(userID, orderID) {
		return models.OrderResponse{}, fmt.Errorf("user is not order owner")
	}

	order, err := o.orderRepository.GetOrder(orderID)
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("error with postgres, couldn't get order: %w", err)
	}

	err = o.getProductsFromOrder(&order)
	if err != nil {
		return models.OrderResponse{}, err
	}

	return order, nil
}

func (o OrderUsecase) GetUserOrders(userID string) ([]models.OrderResponse, error) {
	orders, err := o.orderRepository.GetUserOrders(userID)
	if err != nil {
		return nil, fmt.Errorf("error with postgres: %w", err)
	}

	for idx := range orders {
		err = o.getProductsFromOrder(&orders[idx])
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (o OrderUsecase) GetVendorOrders(vendorID string) ([]models.OrderResponse, error) {
	orders, err := o.orderRepository.GetVendorOrders(vendorID)
	if err != nil {
		return nil, fmt.Errorf("error with postgres: %w", err)
	}

	for idx := range orders {
		err = o.getProductsFromOrder(&orders[idx])
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (o OrderUsecase) UpdateOrderStatus(orderID string, status string) error {
	return o.orderRepository.UpdateOrderStatus(orderID, status)
}

func (o OrderUsecase) getProductsFromOrder(order *models.OrderResponse) error {
	products, err := o.vendorRepository.GetAllProductsWithIDs(order.ProductIDs)
	if err != nil {
		return fmt.Errorf("error with postgres, couldn't get products: %w", err)
	}

	for _, product := range products {
		orderProduct := models.OrderProduct{
			ID:      product.ID,
			Name:    product.Name,
			Price:   product.Price,
			Picture: product.Picture,
		}

		order.Products = append(order.Products, orderProduct)
	}

	return nil
}
