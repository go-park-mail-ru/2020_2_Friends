package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     1,
	}
	cookieName := &session.SessionName{Name: "sessname"}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Create(context.Background(), &session.UserID{Id: "0"}).Times(1).Return(cookieName, nil)

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected: %v\n got: %v", http.StatusCreated, w.Code)
	}

	respCookie := &session.SessionName{Name: w.Result().Cookies()[0].Value}
	if !reflect.DeepEqual(respCookie, cookieName) {
		t.Errorf("expected cookie: %v\n got: %v", cookieName, respCookie)
	}
}

func TestCreateHandlerUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     1,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("", fmt.Errorf("db error"))

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateHandlerProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     1,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(fmt.Errorf("db error"))

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateHandlerSessionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     1,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Create(context.Background(), &session.UserID{Id: "0"}).Times(1).Return(nil, fmt.Errorf("db error"))

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/users", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
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
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	mockSessionClient.EXPECT().Check(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(&session.UserID{Id: "0"}, nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Delete(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(&session.DeleteResponse{}, nil)

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

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

	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)

	mockSessionClient.EXPECT().Check(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(nil, fmt.Errorf("no cookie"))

	handler := UserHandler{
		sessionClient: mockSessionClient,
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

	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)

	mockSessionClient.EXPECT().Check(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(&session.UserID{Id: "0"}, nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(fmt.Errorf("error with db"))

	handler := UserHandler{
		sessionClient: mockSessionClient,
		userUsecase:   mockUserUsecase,
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

	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	mockSessionClient.EXPECT().Check(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(&session.UserID{Id: "0"}, nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(fmt.Errorf("error with db"))

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

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

	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	mockSessionClient.EXPECT().Check(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(&session.UserID{Id: "0"}, nil)
	mockUserUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockProfileUsecase.EXPECT().Delete("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Delete(context.Background(), &session.SessionName{Name: cookie.Value}).Times(1).Return(nil, fmt.Errorf("db error"))

	handler := NewUserHandler(mockUserUsecase, mockSessionClient, mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/users", nil)
	r.AddCookie(&cookie)

	handler.Delete(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
