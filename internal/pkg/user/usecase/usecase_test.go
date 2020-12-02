package usecase

import (
	"fmt"
	"testing"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
	"github.com/golang/mock/gomock"
)

var testUser = models.User{
	Login:    "testlogin",
	Password: "testpassword",
	Role:     1,
}

var userID = "0"

var dbError = fmt.Errorf("db error")

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockUserRepo)

	// without error
	mockUserRepo.EXPECT().Create(testUser).Times(1).Return(userID, nil)

	id, err := userUsecase.Create(testUser)

	expected := userID
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// with error
	mockUserRepo.EXPECT().Create(testUser).Times(1).Return("", dbError)

	id, err = userUsecase.Create(testUser)

	expected = ""
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestCheckIfUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockUserRepo)

	// without error
	mockUserRepo.EXPECT().CheckIfUserExists(testUser).Times(1).Return(nil)

	err := userUsecase.CheckIfUserExists(testUser)

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// with error
	mockUserRepo.EXPECT().CheckIfUserExists(testUser).Times(1).Return(dbError)

	err = userUsecase.CheckIfUserExists(testUser)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockUserRepo)

	// without error
	mockUserRepo.EXPECT().CheckLoginAndPassword(testUser).Times(1).Return(userID, nil)

	id, err := userUsecase.Verify(testUser)

	expected := userID
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// with error
	mockUserRepo.EXPECT().CheckLoginAndPassword(testUser).Times(1).Return("", dbError)

	id, err = userUsecase.Verify(testUser)

	expected = ""
	if id != expected {
		t.Errorf("expected: %v\n got: %v", expected, id)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockUserRepo)

	// without error
	mockUserRepo.EXPECT().Delete(userID).Times(1).Return(nil)

	err := userUsecase.Delete(userID)

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// with error
	mockUserRepo.EXPECT().Delete(userID).Times(1).Return(dbError)

	err = userUsecase.Delete(userID)

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}

func TestCheckUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := user.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockUserRepo)

	// without error
	mockUserRepo.EXPECT().CheckUsersRole(userID).Times(1).Return(testUser.Role, nil)

	role, err := userUsecase.CheckUsersRole(userID)

	expected := testUser.Role
	if role != expected {
		t.Errorf("expected: %v\n got: %v", expected, role)
	}

	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// with error
	mockUserRepo.EXPECT().CheckUsersRole(userID).Times(1).Return(0, dbError)

	role, err = userUsecase.CheckUsersRole(userID)

	expected = 0
	if role != expected {
		t.Errorf("expected: %v\n got: %v", expected, role)
	}

	if err == nil {
		t.Errorf("expected error. Got nil")
	}
}
