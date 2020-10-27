package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
)

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	userID := "0"
	productID := "1"
	vendorID := "2"

	// succesful add
	mock.
		ExpectExec("INSERT INTO carts").
		WithArgs(userID, productID, vendorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(userID, productID, vendorID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// error with db
	mock.
		ExpectExec("INSERT INTO carts").
		WithArgs(userID, productID, vendorID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.Add(userID, productID, vendorID)
	if err == nil {
		t.Error("unexpected error")
		return
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	userID := "0"
	productID := "1"

	// succesful deletion
	mock.
		ExpectExec("DELETE FROM carts").
		WithArgs(userID, productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Remove(userID, productID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// error with db
	mock.
		ExpectExec("DELETE FROM carts").
		WithArgs(userID, productID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.Remove(userID, productID)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	userID := "0"
	products := []models.Product{
		{
			ID:       0,
			VendorID: 0,
			Name:     "name1",
			Price:    "100",
			Picture:  "pic.jpg",
		},
		{
			ID:       1,
			VendorID: 0,
			Name:     "name2",
			Price:    "150",
			Picture:  "img.png",
		},
	}

	// good query
	rows := mock.NewRows([]string{"id", "vendorID", "productName", "price", "picture"})
	for _, prod := range products {
		rows.AddRow(prod.ID, prod.VendorID, prod.Name, prod.Price, prod.Picture)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	resProducts, err := repo.Get(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	if !reflect.DeepEqual(products, resProducts) {
		t.Errorf("expected: %v\ngot:%v", products, resProducts)
		return
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db error"))

	resProducts, err = repo.Get(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	if resProducts != nil {
		t.Errorf("expected: nil\ngot:%v", resProducts)
	}

	// bad query2
	rows = mock.NewRows([]string{"id", "vendorID", "productName", "price"})
	for _, prod := range products {
		rows.AddRow(prod.ID, prod.VendorID, prod.Name, prod.Price)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	resProducts, err = repo.Get(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	if resProducts != nil {
		t.Errorf("expected: nil\ngot:%v", resProducts)
	}
}

func TestGetVendorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	vendorID := 0
	userID := "1"

	row := mock.NewRows([]string{"vendorID"}).AddRow(vendorID)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(row)

	respID, err := repo.GetVendorIDFromCart(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	if respID != vendorID {
		t.Errorf("expected: %v\ngot:%v", respID, vendorID)
	}

	// no rows
	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	respID, err = repo.GetVendorIDFromCart(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	expected := 0
	if respID != expected {
		t.Errorf("expected: %v\ngot:%v", expected, vendorID)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error with db"))

	respID, err = repo.GetVendorIDFromCart(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	expected = 0
	if respID != expected {
		t.Errorf("expected: %v\ngot:%v", expected, vendorID)
	}
}
