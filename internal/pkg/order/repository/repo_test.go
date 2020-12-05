package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
)

var (
	fatalError = "an error '%w' was not expected when opening a stub database connection"

	userID = "1"

	dbError = fmt.Errorf("db error")

	request = models.OrderRequest{
		VendorID:   5,
		VendorName: "test",
		CreatedAt:  time.Now(),
		Address:    "test addr",
		Price:      1000,
		Products: []models.OrderProduct{
			{
				Name:    "test product",
				Price:   1000,
				Picture: "test.png",
			},
		},
	}

	response = models.OrderResponse{
		ID:           10,
		UserID:       50,
		VendorID:     20,
		VendorName:   "test",
		CreatedAt:    time.Now(),
		CreatedAtStr: time.Now().Format(configs.TimeFormat),
		Address:      "test addr",
		Status:       "ready",
		Price:        1000,
		Reviewed:     false,
		Products: []models.OrderProduct{
			{
				Name:    "test1",
				Price:   400,
				Picture: "1.jpg",
			},
			{
				Name:    "test2",
				Price:   600,
				Picture: "2.jpg",
			},
		},
	}
)

func TestAddOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	// good query
	mock.ExpectBegin()

	orderID := 100
	rows := mock.NewRows([]string{"orderID"}).AddRow(orderID)

	mock.
		ExpectQuery("INSERT INTO orders").
		WithArgs(userID, request.VendorID, request.VendorName, request.CreatedAt, request.Address, request.Price).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO products_in_order").
		WithArgs(orderID, request.Products[0].Name, request.Products[0].Price, request.Products[0].Picture).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := repo.AddOrder(userID, request)

	if id != orderID {
		t.Errorf("expected: %v\n got: %v", orderID, id)
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.ExpectBegin()

	mock.
		ExpectQuery("INSERT INTO orders").
		WithArgs(userID, request.VendorID, request.VendorName, request.CreatedAt, request.Address, request.Price).
		WillReturnError(dbError)

	mock.ExpectRollback()

	id, err = repo.AddOrder(userID, request)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// begin return error
	mock.ExpectBegin().WillReturnError(dbError)

	id, err = repo.AddOrder(userID, request)

	expected = 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestAddOrderBadQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	orderID := 100
	rows := mock.NewRows([]string{"orderID"}).AddRow(orderID)

	// bad query 2
	mock.ExpectBegin()

	mock.
		ExpectQuery("INSERT INTO orders").
		WithArgs(userID, request.VendorID, request.VendorName, request.CreatedAt, request.Address, request.Price).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO products_in_order").
		WithArgs(orderID, request.Products[0].Name, request.Products[0].Price, request.Products[0].Picture).
		WillReturnError(dbError)

	mock.ExpectRollback()

	id, err := repo.AddOrder(userID, request)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestAddOrderCommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	orderID := 100
	rows := mock.NewRows([]string{"orderID"}).AddRow(orderID)

	mock.ExpectBegin()

	mock.
		ExpectQuery("INSERT INTO orders").
		WithArgs(userID, request.VendorID, request.VendorName, request.CreatedAt, request.Address, request.Price).
		WillReturnRows(rows)

	mock.
		ExpectExec("INSERT INTO products_in_order").
		WithArgs(orderID, request.Products[0].Name, request.Products[0].Price, request.Products[0].Picture).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit().WillReturnError(dbError)

	mock.ExpectRollback()

	id, err := repo.AddOrder(userID, request)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "vendorID", "vendorName", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.VendorID, response.VendorName, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.ID)).
		WillReturnRows(rows)

	productRows := mock.NewRows([]string{"productName", "price", "picture"})
	for _, prod := range response.Products {
		productRows.AddRow(prod.Name, prod.Price, prod.Picture)
	}

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnRows(productRows)

	order, err := repo.GetOrder(strconv.Itoa(response.ID))

	if !reflect.DeepEqual(response, order) {
		t.Errorf("expected: %v\n got: %v", response, order)
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.ID)).
		WillReturnError(dbError)

	order, err = repo.GetOrder(strconv.Itoa(response.ID))

	expected := models.OrderResponse{}
	if !reflect.DeepEqual(expected, order) {
		t.Errorf("expected: %v\n got: %v", expected, order)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetOrderProductsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "vendorID", "vendorName", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.VendorID, response.VendorName, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.ID)).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnError(dbError)

	order, err := repo.GetOrder(strconv.Itoa(response.ID))

	expected := models.OrderResponse{}
	if !reflect.DeepEqual(expected, order) {
		t.Errorf("expected: %v\n got: %v", expected, order)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetVendorOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnRows(rows)

	productRows := mock.NewRows([]string{"productName", "price", "picture"})
	for _, prod := range response.Products {
		productRows.AddRow(prod.Name, prod.Price, prod.Picture)
	}

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnRows(productRows)

	orders, err := repo.GetVendorOrders(strconv.Itoa(response.VendorID))

	for _, order := range orders {
		order.VendorID = response.VendorID
		order.VendorName = response.VendorName
		if !reflect.DeepEqual(response, order) {
			t.Errorf("expected: %v\n got: %v", response, order)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnError(dbError)

	orders, err = repo.GetVendorOrders(strconv.Itoa(response.VendorID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query 2
	rows = mock.NewRows([]string{"id"})
	rows.AddRow(response.ID)

	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnRows(rows)

	orders, err = repo.GetVendorOrders(strconv.Itoa(response.VendorID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetVendorOrdersProductsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnError(dbError)

	orders, err := repo.GetVendorOrders(strconv.Itoa(response.VendorID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetUserOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "vendorName", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.VendorName, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.UserID)).
		WillReturnRows(rows)

	productRows := mock.NewRows([]string{"productName", "price", "picture"})
	for _, prod := range response.Products {
		productRows.AddRow(prod.Name, prod.Price, prod.Picture)
	}

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnRows(productRows)

	orders, err := repo.GetUserOrders(strconv.Itoa(response.UserID))

	for _, order := range orders {
		order.VendorID = response.VendorID
		if !reflect.DeepEqual(response, order) {
			t.Errorf("expected: %v\n got: %v", response, order)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.UserID)).
		WillReturnError(dbError)

	orders, err = repo.GetUserOrders(strconv.Itoa(response.UserID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query 2
	rows = mock.NewRows([]string{"id"})
	rows.AddRow(response.ID)

	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.UserID)).
		WillReturnRows(rows)

	orders, err = repo.GetUserOrders(strconv.Itoa(response.UserID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetUserOrdersProductsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id", "userID", "vendorName", "createdAt", "clientAddress", "orderStatus", "price", "reviewed"})
	rows.AddRow(response.ID, response.UserID, response.VendorName, response.CreatedAt, response.Address, response.Status, response.Price, response.Reviewed)

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.UserID)).
		WillReturnRows(rows)

	mock.
		ExpectQuery("SELECT productName").
		WithArgs(response.ID).
		WillReturnError(dbError)

	orders, err := repo.GetUserOrders(strconv.Itoa(response.UserID))

	if len(orders) != 0 {
		t.Errorf("expected [], got: %v", orders)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestCheckOrderByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"userID"})
	rows.AddRow(strconv.Itoa(response.UserID))

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.ID)).
		WillReturnRows(rows)

	ok := repo.CheckOrderByUser(strconv.Itoa(response.UserID), strconv.Itoa(response.ID))

	if !ok {
		t.Errorf("expected true")
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.ID)).
		WillReturnError(dbError)

	ok = repo.CheckOrderByUser(strconv.Itoa(response.UserID), strconv.Itoa(response.ID))

	if ok {
		t.Errorf("expected false")
	}
}

func TestGetVendorOrderIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id"})
	for i := 0; i < 2; i++ {
		rows.AddRow(response.ID)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnRows(rows)

	ids, err := repo.GetVendorOrdersIDs(strconv.Itoa(response.VendorID))
	for _, id := range ids {
		if id != response.ID {
			t.Errorf("expected: %v\n got: %v", response.ID, id)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnError(dbError)

	ids, err = repo.GetVendorOrdersIDs(strconv.Itoa(response.VendorID))

	if len(ids) != 0 {
		t.Errorf("expected [], got: %v", ids)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query 2
	rows = mock.NewRows([]string{"id", "userID"})
	for i := 0; i < 2; i++ {
		rows.AddRow(response.ID, response.UserID)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(response.VendorID)).
		WillReturnRows(rows)

	ids, err = repo.GetVendorOrdersIDs(strconv.Itoa(response.VendorID))

	if len(ids) != 0 {
		t.Errorf("expected [], got: %v", ids)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetVendorIDFromOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"vendorID"})
	rows.AddRow(response.VendorID)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnRows(rows)

	id, err := repo.GetVendorIDFromOrder(response.ID)

	if id != response.VendorID {
		t.Errorf("expected: %v\n, got: %v", response.VendorID, id)
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// sql no rows
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnError(sql.ErrNoRows)

	id, err = repo.GetVendorIDFromOrder(response.ID)

	expected := 0
	if id != expected {
		t.Errorf("expected: %v\n, got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnError(dbError)

	id, err = repo.GetVendorIDFromOrder(response.ID)

	expected = 0
	if id != expected {
		t.Errorf("expected: %v\n, got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetUserIDFromOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	userID := strconv.Itoa(response.UserID)

	rows := mock.NewRows([]string{"userID"})
	rows.AddRow(userID)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnRows(rows)

	id, err := repo.GetUserIDFromOrder(response.ID)

	if id != userID {
		t.Errorf("expected: %v\n, got: %v", userID, id)
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// sql no rows
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnError(sql.ErrNoRows)

	id, err = repo.GetUserIDFromOrder(response.ID)

	expected := ""
	if id != expected {
		t.Errorf("expected: %v\n, got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnError(dbError)

	id, err = repo.GetUserIDFromOrder(response.ID)

	expected = ""
	if id != expected {
		t.Errorf("expected: %v\n, got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	// good query
	mock.
		ExpectExec("UPDATE").
		WithArgs(response.Status, strconv.Itoa(response.ID)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateOrderStatus(strconv.Itoa(response.ID), response.Status)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectExec("UPDATE").
		WithArgs(response.Status, strconv.Itoa(response.ID)).
		WillReturnError(dbError)

	err = repo.UpdateOrderStatus(strconv.Itoa(response.ID), response.Status)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestSetOrderReviewStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	// good query
	mock.
		ExpectExec("UPDATE").
		WithArgs(response.Reviewed, response.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SetOrderReviewStatus(response.ID, response.Reviewed)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// bad query
	mock.
		ExpectExec("UPDATE").
		WithArgs(response.Reviewed, response.ID).
		WillReturnError(dbError)

	err = repo.SetOrderReviewStatus(response.ID, response.Reviewed)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetProductsFromOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(response.ID)

	mock.
		ExpectQuery("SELECT").
		WithArgs(response.ID).
		WillReturnRows(rows)

	err = repo.GetProductsFromOrder(&response)
	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}
