package repository

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
)

var fatalError = "an error '%w' was not expected when opening a stub database connection"

var testReview = models.Review{
	UserID:       "0",
	VendorID:     0,
	OrderID:      1,
	Rating:       5,
	Text:         "test",
	CreatedAt:    time.Date(2020, 4, 10, 12, 42, 19, 58, time.Local),
	CreatedAtStr: "10.04.2020 12:42:19",
}

var dbError = fmt.Errorf("db error")

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	// good query
	mock.
		ExpectExec("INSERT").
		WithArgs(testReview.UserID, testReview.OrderID, testReview.VendorID, testReview.Rating, testReview.Text, testReview.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddReview(testReview)

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectExec("INSERT").
		WithArgs(testReview.UserID, testReview.OrderID, testReview.VendorID, testReview.Rating, testReview.Text, testReview.CreatedAt).
		WillReturnError(dbError)

	err = repo.AddReview(testReview)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetUserReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"userID", "orderID", "rating", "review_text", "created_at"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testReview.UserID, testReview.OrderID, testReview.Rating, testReview.Text, testReview.CreatedAt)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testReview.UserID).
		WillReturnRows(rows)

	reviews, err := repo.GetUserReviews(testReview.UserID)

	for _, review := range reviews {
		if !reflect.DeepEqual(testReview, review) {
			t.Errorf("expected: %v\n got: %v", testReview, review)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testReview.UserID).
		WillReturnError(dbError)

	reviews, err = repo.GetUserReviews(testReview.UserID)

	if reviews != nil {
		t.Errorf("expected: %v\n got: %v", nil, reviews)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	rows = mock.NewRows([]string{"userID"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testReview.UserID)
	}

	// bad query 2
	mock.
		ExpectQuery("SELECT").
		WithArgs(testReview.UserID).
		WillReturnRows(rows)

	reviews, err = repo.GetUserReviews(testReview.UserID)

	if reviews != nil {
		t.Errorf("expected: %v\n got: %v", nil, reviews)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetVendorReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"userID", "orderID", "rating", "review_text", "created_at"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testReview.UserID, testReview.OrderID, testReview.Rating, testReview.Text, testReview.CreatedAt)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(testReview.VendorID)).
		WillReturnRows(rows)

	reviews, err := repo.GetVendorReviews(strconv.Itoa(testReview.VendorID))

	for _, review := range reviews {
		if !reflect.DeepEqual(testReview, review) {
			t.Errorf("expected: %v\n got: %v", testReview, review)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(testReview.VendorID)).
		WillReturnError(dbError)

	reviews, err = repo.GetVendorReviews(strconv.Itoa(testReview.VendorID))

	if reviews != nil {
		t.Errorf("expected: %v\n got: %v", nil, reviews)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	rows = mock.NewRows([]string{"userID"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testReview.UserID)
	}

	// bad query 2
	mock.
		ExpectQuery("SELECT").
		WithArgs(strconv.Itoa(testReview.VendorID)).
		WillReturnRows(rows)

	reviews, err = repo.GetVendorReviews(strconv.Itoa(testReview.VendorID))

	if reviews != nil {
		t.Errorf("expected: %v\n got: %v", nil, reviews)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}
