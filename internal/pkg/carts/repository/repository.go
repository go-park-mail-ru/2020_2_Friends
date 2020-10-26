package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/carts"
	"github.com/friends/internal/pkg/models"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) carts.Repository {
	return CartRepository{
		db: db,
	}
}

func (c CartRepository) Add(userID, productID string) error {
	_, err := c.db.Exec(
		"INSERT INTO carts (userID, productID) VALUES($1, $2)",
		userID, productID,
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

func (c CartRepository) Get(userID string) ([]models.Product, error) {
	rows, err := c.db.Query(
		`SELECT p.ID, p.vendorID, p.productName, p.Price, p.picture FROM products AS p
		INNER JOIN carts AS c
		ON p.ID = c.productID WHERE userID=$1`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get products from cart")
	}

	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.ID, &product.VendorID, &product.Name, &product.Price, &product.Picture)
		if err != nil {
			return nil, fmt.Errorf("error in receiving the product: %w", err)
		}

		products = append(products, product)
	}

	return products, nil
}
