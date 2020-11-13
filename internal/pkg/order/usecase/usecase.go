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

	products, err := o.vendorRepository.GetAllProductsWithIDs(order.ProductIDs)
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("error with postgres, couldn't get products: %w", err)
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

	return order, nil
}

func (o OrderUsecase) GetVendorOrders(vendorID string) ([]models.OrderResponse, error) {
	orders, err := o.orderRepository.GetVendorOrders(vendorID)
	if err != nil {
		return nil, fmt.Errorf("error with postgres: %w", err)
	}

	for _, order := range orders {
		products, err := o.vendorRepository.GetAllProductsWithIDs(order.ProductIDs)
		if err != nil {
			return nil, fmt.Errorf("error with postgres, couldn't get products: %w", err)
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
	}

	return orders, nil
}

func (o OrderUsecase) UpdateOrderStatus(orderID string, status string) error {
	return o.orderRepository.UpdateOrderStatus(orderID, status)
}
