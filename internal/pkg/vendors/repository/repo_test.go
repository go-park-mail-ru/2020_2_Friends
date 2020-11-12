package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
	"github.com/lib/pq"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendor := models.Vendor{
		ID:          1,
		Name:        "b",
		Description: "bb",
		Picture:     "b.jpg",
		Products: []models.Product{
			{
				ID:      0,
				Name:    "aaa",
				Price:   "0",
				Picture: "aaa.png",
			},
			{
				ID:      1,
				Name:    "bbb",
				Price:   "1",
				Picture: "bbb.png",
			},
		},
	}

	rows := mock.NewRows([]string{"id", "vendorName", "descript", "picture"})
	rows.AddRow(vendor.ID, vendor.Name, vendor.Description, vendor.Picture)

	// good query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(vendor.ID).
		WillReturnRows(rows)

	productRows := mock.NewRows([]string{"id", "productName", "price", "picture"})
	for _, product := range vendor.Products {
		productRows.AddRow(product.ID, product.Name, product.Price, product.Picture)
	}

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(vendor.ID).
		WillReturnRows(productRows)

	dbVendor, err := repo.Get(vendor.ID)

	if !reflect.DeepEqual(vendor, dbVendor) {
		t.Errorf("expected: %v\n got: %v", vendor, dbVendor)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(vendor.ID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendor, err = repo.Get(vendor.ID)

	if !reflect.DeepEqual(dbVendor, models.Vendor{}) {
		t.Errorf("expected: %v\n got: %v", models.Vendor{}, dbVendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(vendor.ID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(vendor.ID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendor, err = repo.Get(vendor.ID)

	if !reflect.DeepEqual(dbVendor, models.Vendor{}) {
		t.Errorf("expected: %v\n got: %v", models.Vendor{}, dbVendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query2
	productRows = mock.NewRows([]string{"id"})
	for _, product := range vendor.Products {
		productRows.AddRow(product.ID)
	}

	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(vendor.ID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(vendor.ID).
		WillReturnRows(productRows)

	dbVendor, err = repo.Get(vendor.ID)

	if !reflect.DeepEqual(dbVendor, models.Vendor{}) {
		t.Errorf("expected: %v\n got: %v", models.Vendor{}, dbVendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendors := []models.Vendor{
		{
			ID:          0,
			Name:        "a",
			Description: "aa",
			Picture:     "a.png",
		},
		{
			ID:          1,
			Name:        "b",
			Description: "bb",
			Picture:     "b.jpg",
		},
	}

	rows := mock.NewRows([]string{"id", "vendorName", "descript", "picture"})
	for _, vendor := range vendors {
		rows.AddRow(vendor.ID, vendor.Name, vendor.Description, vendor.Picture)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	dbVendors, err := repo.GetAll()

	if !reflect.DeepEqual(vendors, dbVendors) {
		t.Errorf("expected: %v\n got: %v", vendors, dbVendors)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WillReturnError(fmt.Errorf("error with db"))

	dbVendors, err = repo.GetAll()

	if dbVendors != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbVendors)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query2
	rows = mock.NewRows([]string{"id", "vendorName"})
	for _, vendor := range vendors {
		rows.AddRow(vendor.ID, vendor.Name)
	}

	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	dbVendors, err = repo.GetAll()

	if dbVendors != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbVendors)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetAllProductsWithIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	ids := []string{"0", "1"}
	products := []models.Product{
		{
			ID:       0,
			Name:     "aaa",
			Price:    "0",
			Picture:  "aaa.png",
			VendorID: 1,
		},
		{
			ID:       1,
			Name:     "bbb",
			Price:    "1",
			Picture:  "bbb.png",
			VendorID: 1,
		},
	}

	rows := mock.NewRows([]string{"id", "vendorID", "productName", "price", "picture"})
	for _, product := range products {
		rows.AddRow(product.ID, product.VendorID, product.Name, product.Price, product.Picture)
	}

	// good query
	mock.
		ExpectQuery("SELECT id, vendorID").
		WithArgs(pq.Array(ids)).
		WillReturnRows(rows)

	dbProducts, err := repo.GetAllProductsWithIDs(ids)

	if !reflect.DeepEqual(dbProducts, products) {
		t.Errorf("expected: %v\n got: %v", dbProducts, products)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT id, vendorID").
		WithArgs(pq.Array(ids)).
		WillReturnError(fmt.Errorf("error with db"))

	dbProducts, err = repo.GetAllProductsWithIDs(ids)

	if dbProducts != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbProducts)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query2
	rows = mock.NewRows([]string{"id"})
	for _, product := range products {
		rows.AddRow(product.ID)
	}

	mock.
		ExpectQuery("SELECT id, vendorID").
		WithArgs(pq.Array(ids)).
		WillReturnRows(rows)

	dbProducts, err = repo.GetAllProductsWithIDs(ids)

	if dbProducts != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbProducts)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetVendorIDFromProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendorID := "0"
	productID := "0"

	row := mock.NewRows([]string{"vendorID"}).AddRow(vendorID)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(productID).
		WillReturnRows(row)

	dbVendorID, err := repo.GetVendorIDFromProduct(productID)

	if vendorID != dbVendorID {
		t.Errorf("expected: %v\n got: %v", dbVendorID, vendorID)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(productID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendorID, err = repo.GetVendorIDFromProduct(productID)

	if dbVendorID != "" {
		t.Errorf("expected: %v\n got: %v", "", dbVendorID)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestIsVendorExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendorName := "test"

	row := mock.NewRows([]string{"vendorName"}).AddRow(vendorName)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorName).
		WillReturnRows(row)

	err = repo.IsVendorExists(vendorName)

	if err == nil {
		t.Errorf("expected err")
	}

	// no rows
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorName).
		WillReturnError(sql.ErrNoRows)

	err = repo.IsVendorExists(vendorName)

	if err != nil {
		t.Errorf("expected nil, got: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorName).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.IsVendorExists(vendorName)

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendor := models.Vendor{
		ID:          1,
		Name:        "b",
		Description: "bb",
	}

	// good update
	mock.
		ExpectExec("UPDATE").
		WithArgs(vendor.Name, vendor.Description, vendor.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(vendor)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE").
		WithArgs(vendor.Name, vendor.Description, vendor.ID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.Update(vendor)
	if err == nil {
		t.Error("expected error")
		return
	}
}
