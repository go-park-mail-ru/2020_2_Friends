package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
		WithArgs(userID, productID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(userID, productID, vendorID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// error with db
	mock.
		ExpectExec("INSERT INTO carts").
		WithArgs(userID, productID).
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
	ids := []string{"1", "2"}

	// good query
	rows := mock.NewRows([]string{"productID"})
	for _, id := range ids {
		rows.AddRow(id)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	resIDs, err := repo.GetProductIDs(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	if !reflect.DeepEqual(ids, resIDs) {
		t.Errorf("expected: %v\ngot: %v", ids, resIDs)
		return
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db error"))

	resIDs, err = repo.GetProductIDs(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	if resIDs != nil {
		t.Errorf("expected: nil\ngot: %v", resIDs)
	}

	// bad query2
	rows = mock.NewRows([]string{"productID", "vendorID"})
	for _, id := range ids {
		rows.AddRow(id, "0")
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	resIDs, err = repo.GetProductIDs(userID)
	if err == nil {
		t.Error("expected error")
		return
	}

	if resIDs != nil {
		t.Errorf("expected: nil\ngot: %v", resIDs)
	}
}

func TestGetVendorID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewCartRepository(db)

	vendorID := "0"
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
		t.Errorf("expected: %v\ngot: %v", vendorID, respID)
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

	expected := ""
	if respID != expected {
		t.Errorf("expected: %v\ngot: %v", expected, respID)
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

	expected = ""
	if respID != expected {
		t.Errorf("expected: %v\ngot: %v", expected, respID)
	}
}
