package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/friends/internal/pkg/models"
)

func TestCheckIfUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := models.User{
		Login: "testlogin",
	}

	rows := mock.NewRows([]string{"login"}).AddRow(user.Login)

	// test on user exists
	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	err = repo.CheckIfUserExists(user)
	if err == nil {
		t.Errorf("expected err")
		return
	}

	// test on user not exists
	mock.
		ExpectQuery("SELECT login").
		WithArgs(user.Login).
		WillReturnError(sql.ErrNoRows)

	err = repo.CheckIfUserExists(user)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// test on error with db
	bdError := fmt.Errorf("bd error")

	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnError(bdError)

	err = repo.CheckIfUserExists(user)
	if err == nil {
		t.Errorf("expected err")
		return
	}
}

func TestCheckLoginAndPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	password := "testpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("an error '%s' was not expected on hashing", err)
	}

	user := models.User{
		ID:       "0",
		Login:    "testlogin",
		Password: password,
	}

	rows := mock.NewRows([]string{"id", "password"}).AddRow(user.ID, hashedPassword)

	// test on correct login and password
	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	id, err := repo.CheckLoginAndPassword(user)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	if id != user.ID {
		t.Error("expected equal")
		return
	}

	// test om not correct password
	rows = mock.NewRows([]string{"id", "password"}).AddRow(user.ID, "another password")
	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	id, err = repo.CheckLoginAndPassword(user)
	if err == nil {
		t.Errorf("expected err")
		return
	}

	if id != "" {
		t.Errorf("expected empty string")
		return
	}

	// test on user not exists
	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnError(sql.ErrNoRows)

	id, err = repo.CheckLoginAndPassword(user)
	if err == nil {
		t.Errorf("expected err")
		return
	}

	if id != "" {
		t.Errorf("expected empty string")
		return
	}

	// test on db error
	dbError := fmt.Errorf("db error")

	mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnError(dbError)

	id, err = repo.CheckLoginAndPassword(user)
	if err == nil {
		t.Errorf("expected err")
		return
	}

	if id != "" {
		t.Errorf("expected empty string")
		return
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := models.User{
		Login:    "testlogin",
		Password: "testpassword",
	}

	rows := mock.NewRows([]string{"id"}).AddRow(1)

	// successful creation
	mock.
		ExpectQuery("INSERT INTO users").
		WithArgs(user.Login, sqlmock.AnyArg()).
		WillReturnRows(rows)

	id, err := repo.Create(user)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	if id != "1" {
		t.Errorf("expected last id 1")
		return
	}

	// error on creation
	mock.
		ExpectQuery("INSERT INTO users").
		WithArgs(user.Login, sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("erorr with db"))

	id, err = repo.Create(user)
	if err == nil {
		t.Error("expected err")
		return
	}

	if id != "" {
		t.Errorf("expected empty string")
		return
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	userID := "1"
	intUserID := 1

	// successful deletion
	mock.
		ExpectExec("DELETE FROM users").
		WithArgs(intUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// error on deletion
	mock.
		ExpectExec("DELETE FROM users").
		WithArgs(intUserID).
		WillReturnError(fmt.Errorf("db error"))

	err = repo.Delete(userID)
	if err == nil {
		t.Error("expected err")
		return
	}

	// bad userID
	err = repo.Delete("a")
	if err == nil {
		t.Error("expected err")
		return
	}
}
