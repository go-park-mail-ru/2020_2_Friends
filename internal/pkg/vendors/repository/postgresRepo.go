package repository

import (
	"database/sql"
	"errors"
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
	err := row.Scan(
		&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture, &vendor.Longitude, &vendor.Latitude, &vendor.Radius,
	)
	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get vendor: %w", err)
	}
	vendor.HintContent = vendor.Name

	rows, err := v.db.Query(
		"SELECT id, productName, descript, price, picture FROM products WHERE vendorID=$1",
		id,
	)

	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get products for vendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Picture)
		if err != nil {
			return models.Vendor{}, fmt.Errorf("error in receiving the product: %w", err)
		}

		vendor.Products = append(vendor.Products, product)
	}

	categoryRows, err := v.db.Query(
		"SELECT category FROM vendor_categories WHERE vendorid = $1",
		id,
	)

	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get categories for vendor: %w", err)
	}
	defer categoryRows.Close()

	var category string
	for categoryRows.Next() {
		err = categoryRows.Scan(&category)
		if err != nil {
			return models.Vendor{}, fmt.Errorf("error in receiving category: %w", err)
		}

		vendor.Categories = append(vendor.Categories, category)
	}

	return vendor, nil
}

func (v VendorRepository) GetVendorInfo(id string) (models.Vendor, error) {
	row := v.db.QueryRow(
		"SELECT id, vendorName, descript, picture FROM vendors WHERE id=$1",
		id,
	)

	vendor := models.Vendor{}
	err := row.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture)
	if err != nil {
		return models.Vendor{}, fmt.Errorf("couldn't get vendor: %w", err)
	}

	return vendor, nil
}

func (v VendorRepository) GetAll() ([]models.Vendor, error) {
	rows, err := v.db.Query(
		`SELECT id, vendorName, descript, picture, ST_X(coordinates::geometry),
		ST_Y(coordinates::geometry), service_radius FROM vendors`,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendors from db")
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		vendor := models.NewEmptyVendor()
		err = rows.Scan(
			&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture, &vendor.Longitude, &vendor.Latitude, &vendor.Radius,
		)
		if err != nil {
			return nil, fmt.Errorf("error in receiving the vendor: %w", err)
		}
		vendor.HintContent = vendor.Name

		vendors = append(vendors, vendor)
	}

	for idx, vendor := range vendors {
		vendors[idx].Categories = make([]string, 0)
		categoryRows, err := v.db.Query(
			"SELECT category FROM vendor_categories WHERE vendorid = $1",
			vendor.ID,
		)

		if err != nil {
			return nil, fmt.Errorf("couldn't get categories for vendor: %w", err)
		}
		defer categoryRows.Close()

		var category string
		for categoryRows.Next() {
			err = categoryRows.Scan(&category)
			if err != nil {
				return nil, fmt.Errorf("error in receiving category: %w", err)
			}

			vendors[idx].Categories = append(vendors[idx].Categories, category)
		}
	}

	return vendors, nil
}

func (v VendorRepository) GetAllProductsWithIDsFromSameVendor(ids []int) ([]models.Product, error) {
	rows, err := v.db.Query(
		"SELECT id, vendorID, productName, descript, price, picture FROM products WHERE id = ANY ($1)",
		pq.Array(ids),
	)

	if err != nil {
		return nil, ownErr.NewServerError(fmt.Errorf("couldn't get products from products: %w", err))
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	vendorID := -1
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.VendorID, &product.Name, &product.Description, &product.Price, &product.Picture)
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

	err := row.Scan(&vendorName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err == nil {
		return fmt.Errorf("vendor exists")
	}

	return fmt.Errorf("couldn't check vendor, error with db: %w", err)
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
		vendor.Name, vendor.Description, vendor.Longitude, vendor.Latitude, vendor.Radius,
	).Scan(&vendorID)

	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("couldn't insert vendor: %w", err)
	}

	for _, category := range vendor.Categories {
		_, err = tx.Exec(
			"INSERT INTO vendor_categories (vendorID, category) VALUES($1, $2)",
			vendorID, category,
		)

		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("couldn't insert category: %w", err)
		}
	}

	_, err = tx.Exec(
		"INSERT INTO vendor_partner (partnerID, vendorID) VALUES($1, $2)",
		partnerID, vendorID,
	)

	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("couldn't insert partner and vendor: %w", err)
	}

	for _, product := range vendor.Products {
		_, err = tx.Exec(
			"INSERT INTO products (vendorID, productName, price, picture) VALUES($1, $2, $3, $4)",
			vendorID, product.Name, product.Price, product.Picture,
		)

		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("couldn't insert product: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
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
		"INSERT INTO products (vendorID, productName, descript, price) VALUES($1, $2, $3, $4) RETURNING id",
		product.VendorID, product.Name, product.Description, product.Price,
	).Scan(&productID)

	if err != nil {
		return 0, fmt.Errorf("couldn't insert product: %w", err)
	}

	return productID, nil
}

func (v VendorRepository) UpdateProduct(product models.Product) error {
	_, err := v.db.Exec(
		"UPDATE products SET productName = $1, descript = $2, price = $3 WHERE id = $4",
		product.Name, product.Description, product.Price, product.ID,
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

	vendors := make([]models.Vendor, 0)
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

func (v VendorRepository) GetVendorOwner(vendorID int) (string, error) {
	row := v.db.QueryRow(
		"SELECT partnerID FROM vendor_partner WHERE vendorID = $1",
		vendorID,
	)

	var partnerID string
	err := row.Scan(&partnerID)
	if err != nil {
		return "", fmt.Errorf("couldn't get partnerID from vendor with id = %v. Error: %w", vendorID, err)
	}

	return partnerID, nil
}

func (v VendorRepository) GetNearest(longitude, latitude float64) ([]models.Vendor, error) {
	rows, err := v.db.Query(
		`SELECT id, vendorName, ST_X(coordinates::geometry), ST_Y(coordinates::geometry), service_radius
		FROM vendors WHERE ST_DWithin(coordinates, ST_SetSRID(ST_Point($1, $2), 4326), 5 * 1000)`,
		longitude, latitude,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendors from db")
	}
	defer rows.Close()

	vendors := make([]models.Vendor, 0)
	for rows.Next() {
		vendor := models.Vendor{}
		err = rows.Scan(&vendor.ID, &vendor.HintContent, &vendor.Longitude, &vendor.Latitude, &vendor.Radius)
		if err != nil {
			return nil, fmt.Errorf("error in receiving the vendor: %w", err)
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

func (v VendorRepository) GetSimilar(vendorID string, longitude, latitude float64) ([]models.Vendor, error) {
	var geoCondition string
	if latitude != 0 && longitude != 0 {
		geoCondition = fmt.Sprintf(
			"ST_DWithin(coordinates, ST_SetSRID(ST_Point(%v, %v), 4326), 5 * 1000) AND", longitude, latitude,
		)
	}

	rows, err := v.db.Query(
		fmt.Sprintf(`SELECT v.id, v.vendorName, v.descript, v.picture FROM vendors AS v
		JOIN vendor_categories AS vc ON v.id = vc.vendorid
		WHERE %v
		category IN (SELECT category FROM vendor_categories WHERE vendorid = $1) AND vendorid != $1
		GROUP BY v.id
		ORDER BY COUNT(category) DESC`, geoCondition),
		vendorID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get recomendations: %w", err)
	}
	defer rows.Close()

	vendors := make([]models.Vendor, 0)
	vendor := models.Vendor{}
	for rows.Next() {
		err = rows.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture)
		if err != nil {
			return nil, fmt.Errorf("couldn't get recomendations: %w", err)
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

func (v VendorRepository) Get3RandomVendors() ([]models.Vendor, error) {
	rows, err := v.db.Query(
		"SELECT id, vendorName, descript, picture FROM vendors ORDER BY RANDOM() LIMIT 3",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vendors := make([]models.Vendor, 0, 3)
	vendor := models.Vendor{}
	for rows.Next() {
		err = rows.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.Picture)
		if err != nil {
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	return vendors, nil
}

func (v VendorRepository) GetAllCategories() ([]string, error) {
	rows, err := v.db.Query(
		"SELECT category FROM categories",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]string, 0)
	var category string
	for rows.Next() {
		err = rows.Scan(&category)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}
