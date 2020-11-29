package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"

	ownErr "github.com/friends/pkg/error"
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
	tx, err := o.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("couldn't create transaction: %w", err)
	}
	var orderID int

	err = tx.QueryRow(
		`INSERT INTO orders (userID, vendorID, vendorName, createdAt, clientAddress, price)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		userID, order.VendorID, order.VendorName, order.CreatedAt, order.Address, order.Price,
	).Scan(&orderID)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("couldn't insert order: %w", err)
	}

	for _, product := range order.Products {
		_, err = tx.Exec(
			"INSERT INTO products_in_order (orderID, productName, price, picture) VALUES($1, $2, $3, $4)",
			orderID, product.Name, product.Price, product.Picture,
		)

		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("couldn't insert product: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return orderID, nil
}

func (o OrderRepository) GetOrder(orderID string) (models.OrderResponse, error) {
	var order models.OrderResponse
	err := o.db.QueryRow(
		`SELECT id, userID, vendorID, vendorName, createdAt, clientAddress, orderStatus, price, reviewed
		FROM orders WHERE id = $1`,
		orderID,
	).Scan(
		&order.ID, &order.UserID, &order.VendorID, &order.VendorName, &order.CreatedAt,
		&order.Address, &order.Status, &order.Price, &order.Reviewed,
	)
	order.CreatedAtStr = order.CreatedAt.Format(configs.TimeFormat)

	if err != nil {
		return models.OrderResponse{}, ownErr.NewServerError(fmt.Errorf("couldn't get order from db: %w", err))
	}

	err = o.GetProductsFromOrder(&order)
	if err != nil {
		return models.OrderResponse{}, err
	}

	return order, nil
}

func (o OrderRepository) GetUserOrders(userID string) ([]models.OrderResponse, error) {
	rows, err := o.db.Query(
		`SELECT id, userID, vendorName, createdAt, clientAddress, orderStatus, price, reviewed
		FROM orders WHERE userID = $1`,
		userID,
	)

	if err != nil {
		return nil, ownErr.NewServerError(fmt.Errorf("couldn't get orders from db: %w", err))
	}

	orders := make([]models.OrderResponse, 0)
	for rows.Next() {
		var order models.OrderResponse
		err = rows.Scan(
			&order.ID, &order.UserID, &order.VendorName, &order.CreatedAt,
			&order.Address, &order.Status, &order.Price, &order.Reviewed,
		)
		if err != nil {
			return nil, ownErr.NewServerError(fmt.Errorf("couldn't get order from db: %w", err))
		}
		order.CreatedAtStr = order.CreatedAt.Format(configs.TimeFormat)

		err = o.GetProductsFromOrder(&order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
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
		`SELECT id, userID, createdAt, clientAddress, orderStatus, price, reviewed
		FROM orders WHERE vendorID = $1`,
		vendorID,
	)

	if err != nil {
		return nil, ownErr.NewServerError(fmt.Errorf("couldn't get orders from db: %w", err))
	}

	orders := make([]models.OrderResponse, 0)
	for rows.Next() {
		var order models.OrderResponse
		err = rows.Scan(
			&order.ID, &order.UserID, &order.CreatedAt,
			&order.Address, &order.Status, &order.Price, &order.Reviewed,
		)
		if err != nil {
			return nil, ownErr.NewServerError(fmt.Errorf("couldn't get order from db: %w", err))
		}
		order.CreatedAtStr = order.CreatedAt.Format(configs.TimeFormat)

		err = o.GetProductsFromOrder(&order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o OrderRepository) GetVendorOrdersIDs(vendorID string) ([]int, error) {
	rows, err := o.db.Query(
		"SELECT id FROM orders WHERE vendorID = $1",
		vendorID,
	)

	if err != nil {
		return nil, ownErr.NewServerError(fmt.Errorf("couldn't get order ids from db: %w", err))
	}

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, ownErr.NewServerError(fmt.Errorf("couldn't get orderID from db: %w", err))
		}

		ids = append(ids, id)
	}

	return ids, nil
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

func (o OrderRepository) GetProductsFromOrder(order *models.OrderResponse) error {
	rows, err := o.db.Query(
		"SELECT productName, price, picture FROM products_in_order WHERE orderID = $1",
		order.ID,
	)

	if err != nil {
		return fmt.Errorf("couldn't get products from order: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		product := models.OrderProduct{}
		err = rows.Scan(&product.Name, &product.Price, &product.Picture)
		if err != nil {
			return fmt.Errorf("couldn't get product: %w", err)
		}

		order.Products = append(order.Products, product)
	}

	return nil
}

func (o OrderRepository) GetVendorIDFromOrder(orderID int) (int, error) {
	var vendorID int
	err := o.db.QueryRow(
		"SELECT vendorID FROM orders WHERE id = $1",
		orderID,
	).Scan(&vendorID)

	if err == sql.ErrNoRows {
		return 0, ownErr.NewClientError(fmt.Errorf("no such order"))
	}

	if err != nil {
		return 0, ownErr.NewServerError(fmt.Errorf("couldn't get vendorID from db: %w", err))
	}

	return vendorID, nil
}

func (o OrderRepository) SetOrderReviewStatus(orderID int, status bool) error {
	_, err := o.db.Exec(
		"UPDATE orders SET reviewed = $1 WHERE id = $2",
		status, orderID,
	)

	if err != nil {
		return fmt.Errorf("couldn't update order (id = %v) review status. Error: %w", orderID, err)
	}

	return nil
}

func (o OrderRepository) GetUserIDFromOrder(orderID int) (string, error) {
	var userID string
	err := o.db.QueryRow(
		"SELECT userID FROM orders WHERE id = $1",
		orderID,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", ownErr.NewClientError(fmt.Errorf("no such order"))
	}

	if err != nil {
		return "", ownErr.NewServerError(fmt.Errorf("couldn't get userID from db: %w", err))
	}

	return userID, nil
}
