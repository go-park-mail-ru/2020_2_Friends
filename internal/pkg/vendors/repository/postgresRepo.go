package vendors

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
)

type VendorRepository struct {
	db *sql.DB
}

func NewVendorRepository(db *sql.DB) vendors.Repository {
	return VendorRepository{
		db: db,
	}
}

func (v VendorRepository) Get(id int) (models.Vendor, error) {
	row := v.db.QueryRow(
		"SELECT id, vendorName FROM vendors WHERE id=$1",
		id,
	)

	vendor := models.Vendor{}
	err := row.Scan(&vendor.ID, &vendor.Name)
	if err != nil {
		return models.Vendor{}, fmt.Errorf("no such vendor")
	}

	rows, err := v.db.Query(
		"SELECT id, productName, price, picture FROM products WHERE vendorID=$1",
		id,
	)

	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get products for vendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Picture)
		if err != nil {
			continue
		}

		vendor.Products = append(vendor.Products, product)
	}

	return vendor, nil
}

func (v VendorRepository) GetAll() ([]models.Vendor, error) {
	rows, err := v.db.Query(
		"SELECT id, vendorName FROM vendors",
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendors from db")
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		vendor := models.Vendor{}
		err = rows.Scan(&vendor.ID, &vendor.Name)
		if err != nil {
			continue
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}
