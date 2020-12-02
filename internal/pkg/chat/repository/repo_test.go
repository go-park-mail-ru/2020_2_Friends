package repository

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
)

var fatalError = "an error '%w' was not expected when opening a stub database connection"

var testMsg = models.Message{
	OrderID:   0,
	UserID:    "0",
	Text:      "test",
	SentAt:    time.Date(2020, 4, 10, 12, 42, 19, 58, time.Local),
	SentAtStr: "10.04.2020 12:42:19",
}

var dbError = fmt.Errorf("db error")

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	// good query
	mock.
		ExpectExec("INSERT").
		WithArgs(testMsg.OrderID, testMsg.UserID, testMsg.Text, testMsg.SentAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(testMsg)

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectExec("INSERT").
		WithArgs(testMsg.OrderID, testMsg.UserID, testMsg.Text, testMsg.SentAt).
		WillReturnError(dbError)

	err = repo.Save(testMsg)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	rows := mock.NewRows([]string{"userID", "message_text", "sent_at"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testMsg.UserID, testMsg.Text, testMsg.SentAt)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testMsg.OrderID).
		WillReturnRows(rows)

	msgs, err := repo.GetChat(testMsg.OrderID)

	for _, msg := range msgs {
		if !reflect.DeepEqual(testMsg, msg) {
			t.Errorf("expected: %v\n got: %v", testMsg, msg)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(testMsg.OrderID).
		WillReturnError(dbError)

	msgs, err = repo.GetChat(testMsg.OrderID)

	if msgs != nil {
		t.Errorf("expected: %v\n got: %v", nil, msgs)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query 2
	rows = mock.NewRows([]string{"userID"})
	for i := 0; i < 2; i++ {
		rows.AddRow(testMsg.UserID)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(testMsg.OrderID).
		WillReturnRows(rows)

	msgs, err = repo.GetChat(testMsg.OrderID)

	if msgs != nil {
		t.Errorf("expected: %v\n got: %v", nil, msgs)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestGetVendorChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalError, err)
	}
	defer db.Close()

	repo := New(db)

	chat := models.Chat{
		OrderID:        0,
		InterlocutorID: "0",
		LastMsg:        "test",
	}

	rows := mock.NewRows([]string{"orderID", "userID", "message_text"})
	for i := 0; i < 2; i++ {
		rows.AddRow(chat.OrderID, chat.InterlocutorID, chat.LastMsg)
	}

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	chats, err := repo.GetVendorChats([]int{0, 1})

	for _, c := range chats {
		if !reflect.DeepEqual(c, chat) {
			t.Errorf("expected: %v\n got: %v", chat, c)
		}
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(dbError)

	chats, err = repo.GetVendorChats([]int{0, 1})

	if chats != nil {
		t.Errorf("expected: %v\n got: %v", nil, chats)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}

	// bad query 2
	rows = mock.NewRows([]string{"orderID"})
	for i := 0; i < 2; i++ {
		rows.AddRow(chat.OrderID)
	}

	mock.
		ExpectQuery("SELECT").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	chats, err = repo.GetVendorChats([]int{0, 1})

	if chats != nil {
		t.Errorf("expected: %v\n got: %v", nil, chats)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}
