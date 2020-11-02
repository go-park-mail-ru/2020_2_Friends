package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/cart"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) cart.Repository {
	return CartRepository{
		db: db,
	}
}

func (c CartRepository) Add(userID, productID, vendorID string) error {
	_, err := c.db.Exec(
		"INSERT INTO carts (userID, productID, vendorID) VALUES($1, $2, $3)",
		userID, productID, vendorID,
	)

	if err != nil {
		return fmt.Errorf("couldn't add product to cart: %w", err)
	}

	return nil
}

func (c CartRepository) Remove(userID, productID string) error {
	_, err := c.db.Exec(
		"DELETE FROM carts WHERE userID=$1 and productID=$2",
		userID, productID,
	)

	if err != nil {
		return fmt.Errorf("couldn't remove product from cart: %w", err)
	}

	return nil
}

func (c CartRepository) GetProductIDs(userID string) ([]string, error) {
	rows, err := c.db.Query(
		"SELECT productID from carts where userID=$1",
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get product ids from cart: %w", err)
	}

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error in receiving id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (c CartRepository) GetVendorIDFromCart(userID string) (string, error) {
	row := c.db.QueryRow(
		"SELECT vendorID from carts WHERE userID=$1 limit 1",
		userID,
	)

	var vendorID string
	switch err := row.Scan(&vendorID); err {
	case nil:
		return vendorID, nil
	case sql.ErrNoRows:
		return "", cart.ErrCartIsEmpty
	default:
		return "", err
	}
}
