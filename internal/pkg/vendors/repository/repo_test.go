package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
	"github.com/lib/pq"
)

var (
	testVendor = models.Vendor{
		ID:          1,
		Name:        "b",
		HintContent: "b",
		Description: "cc",
		Picture:     "b.jpg",
		Radius:      3,
		Longitude:   1.0,
		Latitude:    2.0,
		Categories:  []string{"Обеды", "Десерты"},
		Products: []models.Product{
			{
				ID:          0,
				Name:        "aaa",
				Price:       0,
				Picture:     "aaa.png",
				Description: "a",
			},
			{
				ID:          1,
				Name:        "bbb",
				Price:       0,
				Picture:     "bbb.png",
				Description: "b",
			},
		},
	}

	userID = "12"

	fatalError = "an error '%s' was not expected when opening a stub database connection"
	dbError    = fmt.Errorf("db error")
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{
		"id", "vendorName", "descript", "picture", "ST_X(coordinates::geometry)",
		"ST_Y(coordinates::geometry)", "service_radius",
	})
	rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Description, testVendor.Picture, testVendor.Longitude, testVendor.Latitude, testVendor.Radius)

	// good query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(testVendor.ID).
		WillReturnRows(rows)

	productRows := mock.NewRows([]string{"id", "productName", "descript", "price", "picture"})
	for _, product := range testVendor.Products {
		productRows.AddRow(product.ID, product.Name, product.Description, product.Price, product.Picture)
	}

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(testVendor.ID).
		WillReturnRows(productRows)

	categoryRows := mock.NewRows([]string{"category"})
	for _, category := range testVendor.Categories {
		categoryRows.AddRow(category)
	}

	mock.
		ExpectQuery("SELECT category").
		WithArgs(testVendor.ID).
		WillReturnRows(categoryRows)

	dbVendor, err := repo.Get(testVendor.ID)

	if !reflect.DeepEqual(testVendor, dbVendor) {
		t.Errorf("expected: %v\n got: %v", testVendor, dbVendor)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(testVendor.ID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendor, err = repo.Get(testVendor.ID)

	if !reflect.DeepEqual(dbVendor, models.Vendor{}) {
		t.Errorf("expected: %v\n got: %v", models.Vendor{}, dbVendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{
		"id", "vendorName", "descript", "picture", "ST_X(coordinates::geometry)",
		"ST_Y(coordinates::geometry)", "service_radius",
	})
	rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Description, testVendor.Picture, testVendor.Longitude, testVendor.Latitude, testVendor.Radius)

	// bad query
	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(testVendor.ID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(testVendor.ID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendor, err := repo.Get(testVendor.ID)

	if !reflect.DeepEqual(dbVendor, models.Vendor{}) {
		t.Errorf("expected: %v\n got: %v", models.Vendor{}, dbVendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetError2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{
		"id", "vendorName", "descript", "picture", "ST_X(coordinates::geometry)",
		"ST_Y(coordinates::geometry)", "service_radius",
	})
	rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Description, testVendor.Picture, testVendor.Longitude, testVendor.Latitude, testVendor.Radius)

	productRows := mock.NewRows([]string{"id"})
	for _, product := range testVendor.Products {
		productRows.AddRow(product.ID)
	}

	mock.
		ExpectQuery("SELECT id, vendorName").
		WithArgs(testVendor.ID).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT id, productName").
		WithArgs(testVendor.ID).
		WillReturnRows(productRows)

	dbVendor, err := repo.Get(testVendor.ID)

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
			Products:    make([]models.Product, 0),
			HintContent: "a",
			Radius:      3,
			Longitude:   1.0,
			Latitude:    2.0,
			Categories:  []string{"Обеды", "Десерты"},
		},
		{
			ID:          1,
			Name:        "b",
			Description: "bb",
			Picture:     "b.jpg",
			Products:    make([]models.Product, 0),
			HintContent: "b",
			Radius:      3,
			Longitude:   1.0,
			Latitude:    2.0,
			Categories:  []string{"Обеды", "Десерты"},
		},
	}

	rows := mock.NewRows([]string{
		"id", "vendorName", "descript", "picture", "ST_X(coordinates::geometry)",
		"ST_Y(coordinates::geometry)", "service_radius",
	})

	for _, vendor := range vendors {
		rows.AddRow(
			vendor.ID, vendor.Name, vendor.Description, vendor.Picture, vendor.Longitude, vendor.Latitude, vendor.Radius,
		)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WillReturnRows(rows)

	for _, vendor := range vendors {
		categoryRows := mock.NewRows([]string{"category"})
		for _, category := range testVendor.Categories {
			categoryRows.AddRow(category)
		}

		mock.
			ExpectQuery("SELECT category").
			WithArgs(vendor.ID).
			WillReturnRows(categoryRows)
	}

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

	ids := []int{0, 1}
	products := []models.Product{
		{
			ID:          0,
			Name:        "aaa",
			Price:       0,
			Picture:     "aaa.png",
			VendorID:    1,
			Description: "a",
		},
		{
			ID:          1,
			Name:        "bbb",
			Price:       0,
			Picture:     "bbb.png",
			VendorID:    1,
			Description: "b",
		},
	}

	rows := mock.NewRows([]string{"id", "vendorID", "productName", "descript", "price", "picture"})
	for _, product := range products {
		rows.AddRow(product.ID, product.VendorID, product.Name, product.Description, product.Price, product.Picture)
	}

	// good query
	mock.
		ExpectQuery("SELECT id, vendorID").
		WithArgs(pq.Array(ids)).
		WillReturnRows(rows)

	dbProducts, err := repo.GetAllProductsWithIDsFromSameVendor(ids)

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

	dbProducts, err = repo.GetAllProductsWithIDsFromSameVendor(ids)

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

	dbProducts, err = repo.GetAllProductsWithIDsFromSameVendor(ids)

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

func TestCheckVendorOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendorID := "0"
	partnerIDs := []string{"0", "1"}

	rows := mock.NewRows([]string{"partnerID"})
	for _, id := range partnerIDs {
		rows.AddRow(id)
	}

	// is owner
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorID).
		WillReturnRows(rows)

	err = repo.CheckVendorOwner(partnerIDs[0], vendorID)
	if err != nil {
		t.Error("expected nil")
		return
	}

	// not owner
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorID).
		WillReturnRows(rows)

	err = repo.CheckVendorOwner("3", vendorID)
	if err == nil {
		t.Error("expected not nil")
		return
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.CheckVendorOwner(partnerIDs[0], vendorID)
	if err == nil {
		t.Error("expected not nil")
		return
	}

	// bad query2
	rows = mock.NewRows([]string{"partnerID", "junk"})
	for _, id := range partnerIDs {
		rows.AddRow(id, "")
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(vendorID).
		WillReturnRows(rows)

	err = repo.CheckVendorOwner(partnerIDs[0], vendorID)
	if err == nil {
		t.Error("expected not nil")
		return
	}
}

func TestAddProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	product := models.Product{
		VendorID:    1,
		Name:        "a",
		Price:       0,
		Description: "aa",
	}

	rows := mock.NewRows([]string{"id"}).AddRow(1)

	// good query
	mock.
		ExpectQuery("INSERT").
		WithArgs(product.VendorID, product.Name, product.Description, product.Price).
		WillReturnRows(rows)

	id, err := repo.AddProduct(product)

	if id != 1 {
		t.Errorf("expected 1, got: %v", id)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("INSERT").
		WithArgs(product.VendorID, product.Name, product.Price).
		WillReturnError(fmt.Errorf("db error"))

	id, err = repo.AddProduct(product)

	if id != 0 {
		t.Errorf("expected 0, got: %v", id)
	}

	if err == nil {
		t.Error("expected err")
	}
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	product := models.Product{
		ID:          0,
		Name:        "a",
		Price:       0,
		Description: "aa",
	}

	// good update
	mock.
		ExpectExec("UPDATE").
		WithArgs(product.Name, product.Description, product.Price, product.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateProduct(product)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE").
		WithArgs(product.Name, product.Price, product.ID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.UpdateProduct(product)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	productID := "0"

	// good query
	mock.
		ExpectExec("DELETE").
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteProduct(productID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad query
	mock.
		ExpectExec("DELETE").
		WithArgs(productID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.DeleteProduct(productID)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestUpdateVendorImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	vendor := models.Vendor{
		ID:      0,
		Picture: "image.png",
	}

	// good update
	mock.
		ExpectExec("UPDATE").
		WithArgs(vendor.Picture, strconv.Itoa(vendor.ID)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateVendorImage(strconv.Itoa(vendor.ID), vendor.Picture)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE").
		WithArgs(vendor.Picture, strconv.Itoa(vendor.ID)).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.UpdateVendorImage(strconv.Itoa(vendor.ID), vendor.Picture)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestUpdateProductImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	product := models.Product{
		ID:      0,
		Picture: "image.png",
	}

	// good update
	mock.
		ExpectExec("UPDATE").
		WithArgs(product.Picture, strconv.Itoa(product.ID)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateProductImage(strconv.Itoa(product.ID), product.Picture)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE").
		WithArgs(product.Picture, strconv.Itoa(product.ID)).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.UpdateProductImage(strconv.Itoa(product.ID), product.Picture)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestGetPartnerShops(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	partnerID := "0"

	vendors := []models.Vendor{
		{
			ID:          1,
			Name:        "a",
			Description: "a",
			Picture:     "a.jpg",
			Products:    make([]models.Product, 0),
		},
		{
			ID:          2,
			Name:        "b",
			Description: "bb",
			Picture:     "b.jpg",
			Products:    make([]models.Product, 0),
		},
	}

	rows := mock.NewRows([]string{"id", "vendorName", "descript", "picture"})
	for _, vendor := range vendors {
		rows.AddRow(vendor.ID, vendor.Name, vendor.Description, vendor.Picture)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(partnerID).
		WillReturnRows(rows)

	dbVendors, err := repo.GetPartnerShops(partnerID)

	if !reflect.DeepEqual(vendors, dbVendors) {
		t.Errorf("expected: %v\n got: %v", vendors, dbVendors)
	}

	if err != nil {
		t.Error("unexpected err: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(partnerID).
		WillReturnError(fmt.Errorf("db error"))

	dbVendors, err = repo.GetPartnerShops(partnerID)

	if dbVendors != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbVendors)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query2
	rows = mock.NewRows([]string{"id"})
	for _, vendor := range vendors {
		rows.AddRow(vendor.ID)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(partnerID).
		WillReturnRows(rows)

	dbVendors, err = repo.GetPartnerShops(partnerID)

	if dbVendors != nil {
		t.Errorf("expected: %v\n got: %v", nil, dbVendors)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestVendorInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{"id", "vendorName", "descript", "picture"})
	rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Description, testVendor.Picture)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(testVendor.ID)).
		WillReturnRows(rows)

	vendor, err := repo.GetVendorInfo(strconv.Itoa(testVendor.ID))
	if vendor.ID != testVendor.ID &&
		vendor.Name != testVendor.Name &&
		vendor.Description != testVendor.Description &&
		vendor.Picture != testVendor.Picture {
		t.Errorf("expected: %v\n got: %v", testVendor, vendor)
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	//bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(testVendor.ID)).
		WillReturnError(dbError)

	vendor, err = repo.GetVendorInfo(strconv.Itoa(testVendor.ID))
	expected := models.Vendor{}

	if !reflect.DeepEqual(vendor, expected) {
		t.Errorf("expected: nil\n got: %v", vendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetVendorFromProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{"id", "vendorName", "descript", "picture"})
	rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Description, testVendor.Picture)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testVendor.Products[0].ID).
		WillReturnRows(rows)

	vendor, err := repo.GetVendorFromProduct(testVendor.Products[0].ID)
	if vendor.ID != testVendor.ID &&
		vendor.Name != testVendor.Name &&
		vendor.Description != testVendor.Description &&
		vendor.Picture != testVendor.Picture {
		t.Errorf("expected: %v\n got: %v", testVendor, vendor)
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	//bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testVendor.Products[0].ID).
		WillReturnError(dbError)

	vendor, err = repo.GetVendorFromProduct(testVendor.Products[0].ID)
	expected := models.Vendor{}

	if !reflect.DeepEqual(vendor, expected) {
		t.Errorf("expected: nil\n got: %v", vendor)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	// good query
	mock.ExpectBegin()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(testVendor.ID)

	mock.
		ExpectQuery("INSERT INTO vendors").
		WithArgs(testVendor.Name, testVendor.Description, testVendor.Longitude, testVendor.Latitude, testVendor.Radius).
		WillReturnRows(rows)

	for _, category := range testVendor.Categories {
		mock.
			ExpectExec("INSERT INTO vendor_categories").
			WithArgs(testVendor.ID, category).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.
		ExpectExec("INSERT INTO vendor_partner").
		WithArgs(userID, testVendor.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, product := range testVendor.Products {
		mock.
			ExpectExec("INSERT INTO products").
			WithArgs(testVendor.ID, product.Name, product.Price, product.Picture).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	id, err := repo.Create(userID, testVendor)

	if id != testVendor.ID {
		t.Errorf("expected: %v\n got: %v", testVendor.ID, id)
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	// begin error
	mock.ExpectBegin().WillReturnError(dbError)

	id, err = repo.Create(userID, testVendor)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestCreateCommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	mock.ExpectBegin()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(testVendor.ID)

	// commit error
	mock.
		ExpectQuery("INSERT INTO vendors").
		WithArgs(testVendor.Name, testVendor.Description, testVendor.Longitude, testVendor.Latitude, testVendor.Radius).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO vendor_partner").
		WithArgs(userID, testVendor.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, product := range testVendor.Products {
		mock.
			ExpectExec("INSERT INTO products").
			WithArgs(testVendor.ID, product.Name, product.Price, product.Picture).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit().WillReturnError(dbError)

	mock.ExpectRollback()

	id, err := repo.Create(userID, testVendor)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestCreateProductsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	mock.ExpectBegin()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(testVendor.ID)

	mock.
		ExpectQuery("INSERT INTO vendors").
		WithArgs(testVendor.Name, testVendor.Description, testVendor.Longitude, testVendor.Latitude, testVendor.Radius).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO vendor_partner").
		WithArgs(userID, testVendor.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, product := range testVendor.Products {
		mock.
			ExpectExec("INSERT INTO products").
			WithArgs(testVendor.ID, product.Name, product.Price, product.Picture).
			WillReturnError(dbError)
	}

	mock.ExpectRollback()

	id, err := repo.Create(userID, testVendor)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestCreatePartnerError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	mock.ExpectBegin()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(testVendor.ID)

	mock.
		ExpectQuery("INSERT INTO vendors").
		WithArgs(testVendor.Name, testVendor.Description, testVendor.Longitude, testVendor.Latitude, testVendor.Radius).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO vendor_partner").
		WithArgs(userID, testVendor.ID).
		WillReturnError(dbError)

	mock.ExpectRollback()

	id, err := repo.Create(userID, testVendor)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestCreateVendorError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	mock.ExpectBegin()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(testVendor.ID)

	mock.
		ExpectQuery("INSERT INTO vendors").
		WithArgs(testVendor.Name, testVendor.Description, testVendor.Longitude, testVendor.Latitude, testVendor.Radius).
		WillReturnError(dbError)

	mock.ExpectRollback()

	id, err := repo.Create(userID, testVendor)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetVendorOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{"partnerID"})
	rows.AddRow(userID)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testVendor.ID).
		WillReturnRows(rows)

	id, err := repo.GetVendorOwner(testVendor.ID)

	if id != userID {
		t.Errorf("expected: %v\n got: %v", userID, id)
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testVendor.ID).
		WillReturnError(dbError)

	id, err = repo.GetVendorOwner(testVendor.ID)

	expected := ""
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}

func TestGetNeares(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := NewVendorRepository(db)

	rows := mock.NewRows([]string{"id", "vendorName", "ST_X(coordinates::geometry)", "ST_Y(coordinates::geometry)", "service_radius"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testVendor.ID, testVendor.Name, testVendor.Longitude, testVendor.Latitude, testVendor.Radius)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(float64(testVendor.Longitude), float64(testVendor.Latitude)).
		WillReturnRows(rows)

	vs, err := repo.GetNearest(float64(testVendor.Longitude), float64(testVendor.Latitude))

	for _, v := range vs {
		if v.ID != testVendor.ID &&
			v.Name != testVendor.Name &&
			v.Longitude != testVendor.Longitude &&
			v.Latitude != testVendor.Latitude &&
			v.Radius != testVendor.Radius {
			t.Errorf("expected: %v\n got: %v", testVendor, v)
		}
	}

	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(float64(testVendor.Longitude), float64(testVendor.Latitude)).
		WillReturnError(dbError)

	vs, err = repo.GetNearest(float64(testVendor.Longitude), float64(testVendor.Latitude))

	if len(vs) != 0 {
		t.Errorf("expected: []\n got: %v", vs)
	}

	if err == nil {
		t.Errorf("expected err")
	}

	// bad query 2
	rows = mock.NewRows([]string{"id"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testVendor.ID)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(float64(testVendor.Longitude), float64(testVendor.Latitude)).
		WillReturnRows(rows)

	vs, err = repo.GetNearest(float64(testVendor.Longitude), float64(testVendor.Latitude))

	if len(vs) != 0 {
		t.Errorf("expected: []\n got: %v", vs)
	}

	if err == nil {
		t.Errorf("expected err")
	}
}
