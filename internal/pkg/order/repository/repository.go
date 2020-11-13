package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}

func New(db *sql.DB) order.Repository {
	return OrderRepository{
		db: db,
	}
}

func (o OrderRepository) AddOrder(userID string, order models.OrderRequest) (int, error) {
	var orderID int

	err := o.db.QueryRow(
		"INSERT INTO orders (userID, vendorID, vendorName, products, createdAt, clientAddress) VALUES($1, $2, $3, $4, $5) RETURNING id",
		userID, order.VendorID, order.VendorName, pq.Array(order.Products), order.CreatedAt, order.Address,
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("couldn't insert order: %w", err)
	}

	return orderID, nil
}

func (o OrderRepository) GetOrder(orderID string) (models.OrderResponse, error) {
	var dbProducts pq.Int64Array
	var order models.OrderResponse
	err := o.db.QueryRow(
		`SELECT id, userID, vendorName, products, createdAt, clientAddress, orderStatus
		FROM orders WHERE id = $1`,
		orderID,
	).Scan(&order.ID, &order.UserID, &order.VendorName, &dbProducts, &order.CreatedAt, &order.Address, &order.Status)

	if err != nil {
		return models.OrderResponse{}, fmt.Errorf("couldn't get order from db: %w", err)
	}

	for _, product := range dbProducts {
		order.ProductIDs = append(order.ProductIDs, strconv.Itoa(int(product)))
	}

	return order, nil
}

func (o OrderRepository) CheckOrderByUser(userID string, orderID string) bool {
	var dbUserID string
	err := o.db.QueryRow(
		"SELECT userID FROM orders WHERE id = $1",
		orderID,
	).Scan(&dbUserID)

	if err != nil {
		return false
	}

	return userID == dbUserID
}

func (o OrderRepository) GetVendorOrders(vendorID string) ([]models.OrderResponse, error) {
	rows, err := o.db.Query(
		`SELECT id, userID, vendorName, products, createdAt, clientAddress, orderStatus
		FROM orders WHERE vendorID = $1`,
		vendorID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get orders from db: %w", err)
	}

	var orders []models.OrderResponse
	for rows.Next() {
		var order models.OrderResponse
		var dbProducts pq.Int64Array
		err = rows.Scan(&order.ID, &order.UserID, &order.VendorName, &dbProducts, &order.CreatedAt, &order.Address, &order.Status)
		if err != nil {
			return nil, fmt.Errorf("couldn't get order from db: %w", err)
		}

		for _, product := range dbProducts {
			order.ProductIDs = append(order.ProductIDs, strconv.Itoa(int(product)))
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o OrderRepository) UpdateOrderStatus(orderID string, status string) error {
	_, err := o.db.Exec(
		"UPDATE orders SET orderStatus = $1 WHERE id = $2",
		status, orderID,
	)

	if err != nil {
		return fmt.Errorf("couldn't update status on orderID: %w", err)
	}

	return nil
}
