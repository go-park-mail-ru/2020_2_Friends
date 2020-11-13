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

	order, productIDs, err := o.orderRepository.GetOrder(orderID)
	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("error with postgres, couldn't get order: %w", err)
	}

	products, err := o.vendorRepository.GetAllProductsWithIDs(productIDs)
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
