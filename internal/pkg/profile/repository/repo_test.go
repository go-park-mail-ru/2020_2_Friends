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

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	userID := "0"

	// successful creation
	mock.
		ExpectExec("INSERT INTO profiles").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// error on creation
	mock.
		ExpectExec("INSERT INTO profiles").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error with db"))

	err = repo.Create(userID)
	if err == nil {
		t.Error("expected err")
		return
	}
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	profile := dbProfile{
		UserID:    "0",
		Name:      sql.NullString{String: "testname", Valid: true},
		Phone:     sql.NullString{String: "0000", Valid: true},
		Avatar:    sql.NullString{String: "avatar.jpg", Valid: true},
		Points:    sql.NullInt64{Int64: 0, Valid: true},
		Addresses: pq.StringArray([]string{"addr1", "addr2"}),
	}

	rows := mock.NewRows([]string{"userID", "username", "phone", "addresses", "points", "avatar"})
	rows.AddRow(profile.UserID, profile.Name, profile.Phone, profile.Addresses, profile.Points, profile.Avatar)

	// good query
	mock.
		ExpectQuery("SELECT").
		WithArgs(profile.UserID).
		WillReturnRows(rows)

	resultProfile, err := repo.Get(profile.UserID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	expectedProfile := fromDBToApp(profile)
	if !reflect.DeepEqual(resultProfile, expectedProfile) {
		t.Errorf("expected %v\n got: %v", expectedProfile, resultProfile)
		return
	}

	// bad query
	mock.
		ExpectQuery("SELECT").
		WithArgs(profile.UserID).
		WillReturnError(fmt.Errorf("error with db"))

	resultProfile, err = repo.Get(profile.UserID)
	if err == nil {
		t.Error("expected error")
		return
	}

	expectedProfile = models.Profile{}
	if !reflect.DeepEqual(resultProfile, expectedProfile) {
		t.Errorf("expected %v\n got: %v", expectedProfile, resultProfile)
		return
	}

	// sql no rows
	mock.
		ExpectQuery("SELECT").
		WithArgs(profile.UserID).
		WillReturnError(sql.ErrNoRows)

	resultProfile, err = repo.Get(profile.UserID)
	if err == nil {
		t.Error("expected error")
		return
	}

	expectedProfile = models.Profile{}
	if !reflect.DeepEqual(resultProfile, expectedProfile) {
		t.Errorf("expected %v\n got: %v", expectedProfile, resultProfile)
		return
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	profile := dbProfile{
		UserID:    "0",
		Name:      sql.NullString{String: "testname", Valid: true},
		Phone:     sql.NullString{String: "0000", Valid: true},
		Avatar:    sql.NullString{String: "avatar.jpg", Valid: true},
		Points:    sql.NullInt64{Int64: 0, Valid: true},
		Addresses: pq.StringArray([]string{"addr1", "addr2"}),
	}

	// good update
	mock.
		ExpectExec("UPDATE profiles").
		WithArgs(profile.Name, profile.Phone, profile.Addresses, profile.Points, profile.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(fromDBToApp(profile))
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE profiles").
		WithArgs(profile.Name, profile.Phone, profile.Addresses, profile.Points, profile.UserID).
		WillReturnError(fmt.Errorf("error with db"))

	err = repo.Update(fromDBToApp(profile))
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestUpdateAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	userID := "0"
	link := "avatar.jpg"

	// good update
	mock.
		ExpectExec("UPDATE profiles").
		WithArgs(link, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateAvatar(userID, link)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE profiles").
		WithArgs(link, userID).
		WillReturnError(fmt.Errorf("error with db"))

	err = repo.UpdateAvatar(userID, link)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestUpdateAddresses(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	userID := "0"
	addresses := []string{"addr1", "addr2"}

	// good update
	mock.
		ExpectExec("UPDATE").
		WithArgs(pq.StringArray(addresses), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateAddresses(userID, addresses)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad update
	mock.
		ExpectExec("UPDATE").
		WithArgs(pq.StringArray(addresses), userID).
		WillReturnError(fmt.Errorf("error with db"))

	err = repo.UpdateAddresses(userID, addresses)
	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProfileRepository(db)

	userID := "0"

	// good query
	mock.
		ExpectExec("DELETE FROM profiles").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(userID)
	if err != nil {
		t.Error("unexpected err: %w", err)
		return
	}

	// bad query
	mock.
		ExpectExec("DELETE FROM profiles").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("error with db"))

	err = repo.Delete(userID)
	if err == nil {
		t.Error("expected error")
		return
	}
}
