package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
	ownErr "github.com/friends/pkg/error"
	"github.com/lib/pq"
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
		`SELECT id, vendorName, descript, picture, ST_X(coordinates::geometry), ST_Y(coordinates::geometry), service_radius
		FROM vendors WHERE id=$1`,
		id,
	)

	vendor := models.NewEmptyVendor()
	err := row.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture, &vendor.Longtitude, &vendor.Latitude, &vendor.Radius)
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
			return models.Vendor{}, fmt.Errorf("error in receiving the product: %w", err)
		}

		vendor.Products = append(vendor.Products, product)
	}

	return vendor, nil
}

func (v VendorRepository) GetAll() ([]models.Vendor, error) {
	rows, err := v.db.Query(
		"SELECT id, vendorName, descript, picture, ST_X(coordinates::geometry), ST_Y(coordinates::geometry), service_radius FROM vendors",
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendors from db")
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		vendor := models.NewEmptyVendor()
		err = rows.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture, &vendor.Longtitude, &vendor.Latitude, &vendor.Radius)
		if err != nil {
			return nil, fmt.Errorf("error in receiving the vendor: %w", err)
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

func (v VendorRepository) GetAllProductsWithIDsFromSameVendor(ids []int) ([]models.Product, error) {
	rows, err := v.db.Query(
		"SELECT id, vendorID, productName, price, picture FROM products WHERE id = ANY ($1)",
		pq.Array(ids),
	)

	if err != nil {
		return nil, ownErr.NewServerError(fmt.Errorf("couldn't get products from products: %w", err))
	}
	defer rows.Close()

	var products []models.Product
	vendorID := -1
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.VendorID, &product.Name, &product.Price, &product.Picture)
		if err != nil {
			return nil, ownErr.NewServerError(fmt.Errorf("error in receiving product: %w", err))
		}

		if product.VendorID != vendorID && vendorID != -1 {
			return nil, ownErr.NewClientError(fmt.Errorf("products from different vendors"))
		}

		products = append(products, product)
		vendorID = product.VendorID
	}

	return products, nil
}

func (v VendorRepository) GetVendorIDFromProduct(productID string) (string, error) {
	row := v.db.QueryRow(
		"SELECT vendorID FROM products WHERE id = $1",
		productID,
	)

	var vendorID string
	err := row.Scan(&vendorID)
	if err != nil {
		return "", fmt.Errorf("couldn't get vendorID from products: %w", err)
	}

	return vendorID, nil
}

func (v VendorRepository) GetVendorFromProduct(productID int) (models.Vendor, error) {
	vendor := models.Vendor{}
	err := v.db.QueryRow(
		`SELECT v.id, v.vendorName, v.descript, v.picture FROM vendors AS v
		JOIN products AS p on v.id = p.vendorID
		WHERE p.id = $1`,
		productID,
	).Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture)

	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get vendor: %w", err)
	}

	return vendor, nil
}

func (v VendorRepository) IsVendorExists(vendorName string) error {
	row := v.db.QueryRow(
		"SELECT vendorName FROM vendors WHERE vendorName=$1",
		vendorName,
	)

	var vn string
	switch err := row.Scan(&vn); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return fmt.Errorf("vendor exists")
	default:
		return fmt.Errorf("couldn't check vendor, error with db: %w", err)
	}
}

func (v VendorRepository) Create(partnerID string, vendor models.Vendor) (int, error) {
	tx, err := v.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("couldn't create transaction: %w", err)
	}
	var vendorID int

	err = tx.QueryRow(
		`INSERT INTO vendors (vendorName, descript, coordinates, service_radius)
		VALUES ($1, $2, ST_SetSRID(ST_Point($3, $4), 4326), $5) RETURNING id`,
		vendor.Name, vendor.Description, vendor.Longtitude, vendor.Latitude, vendor.Radius,
	).Scan(&vendorID)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("couldn't insert vendor: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO vendor_partner (partnerID, vendorID) VALUES($1, $2)",
		partnerID, vendorID,
	)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("couldn't insert partner and vendor: %w", err)
	}

	for _, product := range vendor.Products {
		_, err := tx.Exec(
			"INSERT INTO products (vendorID, productName, price, picture) VALUES($1, $2, $3, $4)",
			vendorID, product.Name, product.Price, product.Picture,
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

	return vendorID, nil
}

func (v VendorRepository) Update(vendor models.Vendor) error {
	_, err := v.db.Exec(
		"UPDATE vendors SET vendorName = $1, descript = $2 WHERE id = $3",
		vendor.Name, vendor.Description, vendor.ID,
	)

	if err != nil {
		return fmt.Errorf("couln't update vendor: %w", err)
	}

	return nil
}

func (v VendorRepository) CheckVendorOwner(userID, vendorID string) error {
	rows, err := v.db.Query(
		"SELECT partnerID FROM vendor_partner WHERE vendorID = $1",
		vendorID,
	)

	if err != nil {
		return fmt.Errorf("couldn't check vendors owner: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var partner string
		err = rows.Scan(&partner)
		if err != nil {
			return fmt.Errorf("error with scaning partner id: %w", err)
		}

		if partner == userID {
			return nil
		}
	}

	return fmt.Errorf("no rights for this vendor")
}

func (v VendorRepository) AddProduct(product models.Product) (int, error) {
	var productID int

	err := v.db.QueryRow(
		"INSERT INTO products (vendorID, productName, price) VALUES($1, $2, $3) RETURNING id",
		product.VendorID, product.Name, product.Price,
	).Scan(&productID)

	if err != nil {
		return 0, fmt.Errorf("couldn't insert product: %w", err)
	}

	return productID, nil
}

func (v VendorRepository) UpdateProduct(product models.Product) error {
	_, err := v.db.Exec(
		"UPDATE products SET productName = $1, price = $2 WHERE id = $3",
		product.Name, product.Price, product.ID,
	)

	if err != nil {
		return fmt.Errorf("couln't update product: %w", err)
	}

	return nil
}

func (v VendorRepository) DeleteProduct(productID string) error {
	_, err := v.db.Exec(
		"DELETE FROM products WHERE id = $1",
		productID,
	)

	if err != nil {
		return fmt.Errorf("couln't delete product: %w", err)
	}

	return nil
}

func (v VendorRepository) UpdateVendorImage(vendorID string, link string) error {
	_, err := v.db.Exec(
		"UPDATE vendors SET picture=$1 WHERE id=$2",
		link, vendorID,
	)

	if err != nil {
		return fmt.Errorf("couln't update vendors picture: %w", err)
	}

	return nil
}

func (v VendorRepository) UpdateProductImage(productID string, link string) error {
	_, err := v.db.Exec(
		"UPDATE products SET picture=$1 WHERE id=$2",
		link, productID,
	)

	if err != nil {
		return fmt.Errorf("couln't update products picture: %w", err)
	}

	return nil
}

func (v VendorRepository) GetPartnerShops(partnerID string) ([]models.Vendor, error) {
	rows, err := v.db.Query(
		`SELECT v.id, v.vendorName, v.descript, v.picture FROM vendors AS v
		JOIN vendor_partner AS vp ON v.id = vp.vendorID
		WHERE vp.partnerID = $1`,
		partnerID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendors from db")
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		vendor := models.NewEmptyVendor()
		err = rows.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture)
		if err != nil {
			return nil, fmt.Errorf("error in receiving the vendor: %w", err)
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}
