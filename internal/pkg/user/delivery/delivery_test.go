package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/golang/mock/gomock"
)

func TestCreateHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionUsecase := session.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
	}
	cookieName := "sessname"

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(nil)
	mockSessionUsecase.EXPECT().Create("0").Times(1).Return(cookieName, nil)

	handler := NewUserHandler(mockUserUsecase, mockSessionUsecase, mockProfileUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected: %v\n got: %v", http.StatusCreated, w.Code)
	}

	respCookie := w.Result().Cookies()[0].Value
	if respCookie != cookieName {
		t.Errorf("expected cookie: %v\n got: %v", cookieName, respCookie)
	}
}

func TestCreateHandlerBadJson(t *testing.T) {
	handler := UserHandler{}

	body := bytes.NewReader([]byte(`{"login": "test", "password": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestCreateHandlerUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(fmt.Errorf("user exists"))

	handler := UserHandler{
		userUsecase: mockUserUsecase,
	}

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	expected := http.StatusConflict
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  configs.SessionID,
		Value: "testcookie",
	}

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionUsecase := session.NewMockUsecase(ctrl)

	mockSessionUsecase.EXPECT().Check(cookie.Value).Times(1).Return("0", nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockSessionUsecase.EXPECT().Delete(cookie.Value).Times(1).Return(nil)

	handler := NewUserHandler(mockUserUsecase, mockSessionUsecase, mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteNoCookie(t *testing.T) {
	handler := UserHandler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)

	handler.Delete(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteNoCookie2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  configs.SessionID,
		Value: "testcookie",
	}

	mockSessionUsecase := session.NewMockUsecase(ctrl)

	mockSessionUsecase.EXPECT().Check(cookie.Value).Times(1).Return("", fmt.Errorf("no cookie"))

	handler := UserHandler{
		sessionUsecase: mockSessionUsecase,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteUserDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  configs.SessionID,
		Value: "testcookie",
	}

	mockSessionUsecase := session.NewMockUsecase(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)

	mockSessionUsecase.EXPECT().Check(cookie.Value).Times(1).Return("0", nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(fmt.Errorf("error with db"))

	handler := UserHandler{
		sessionUsecase: mockSessionUsecase,
		userUsecase:    mockUserUsecase,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteProfileDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  configs.SessionID,
		Value: "testcookie",
	}

	mockSessionUsecase := session.NewMockUsecase(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	mockSessionUsecase.EXPECT().Check(cookie.Value).Times(1).Return("0", nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(fmt.Errorf("error with db"))

	handler := NewUserHandler(mockUserUsecase, mockSessionUsecase, mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteSessionDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  configs.SessionID,
		Value: "testcookie",
	}

	mockSessionUsecase := session.NewMockUsecase(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	mockSessionUsecase.EXPECT().Check(cookie.Value).Times(1).Return("0", nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockSessionUsecase.EXPECT().Delete(cookie.Value).Times(1).Return(fmt.Errorf("error with db"))

	handler := NewUserHandler(mockUserUsecase, mockSessionUsecase, mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
