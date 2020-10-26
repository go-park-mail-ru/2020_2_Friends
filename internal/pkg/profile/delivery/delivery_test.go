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
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/golang/mock/gomock"
)

func TestGetHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	prof := models.Profile{
		UserID:    "0",
		Name:      "testname",
		Phone:     "0000",
		Addresses: []string{"addr1", "addr2"},
		Points:    100,
		Avatar:    "avatar.jpg",
	}

	mockProfileUsecase.EXPECT().Get(prof.UserID).Times(1).Return(prof, nil)

	handler := ProfileDelivery{
		profUsecase: mockProfileUsecase,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/profiles", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), prof.UserID)

	handler.Get(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respProf models.Profile
	json.Unmarshal(w.Body.Bytes(), &respProf)
	if !reflect.DeepEqual(prof, respProf) {
		t.Errorf("expected: %v\n got: %v", prof, respProf)
	}
}

func TestGetHandlerBad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	userID := "0"

	mockProfileUsecase.EXPECT().Get(userID).Times(1).Return(models.Profile{}, fmt.Errorf("error with profile"))

	handler := ProfileDelivery{
		profUsecase: mockProfileUsecase,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/profiles", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.Get(w, r.WithContext(ctx))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected: %v\n got: %v", http.StatusCreated, w.Code)
	}
}

func TestUpdateHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := ProfileDelivery{
		profUsecase: mockProfileUsecase,
	}

	prof := models.Profile{
		UserID: "0",
		Name:   "testname",
		Phone:  "0000",
	}

	mockProfileUsecase.EXPECT().Update(prof).Times(1).Return(nil)

	body := bytes.NewReader([]byte(fmt.Sprintf(`{"userId": "%s", "name": "%s", "phone": "%s"}`, prof.UserID, prof.Name, prof.Phone)))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), prof.UserID)

	handler.Update(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := ProfileDelivery{
		profUsecase: mockProfileUsecase,
	}

	prof := models.Profile{
		UserID: "0",
		Name:   "testname",
		Phone:  "0000",
	}

	mockProfileUsecase.EXPECT().Update(prof).Times(1).Return(fmt.Errorf("error"))

	body := bytes.NewReader([]byte(fmt.Sprintf(`{"userId": "%s", "name": "%s", "phone": "%s"}`, prof.UserID, prof.Name, prof.Phone)))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), prof.UserID)

	handler.Update(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateHandlerBadJson(t *testing.T) {
	handler := ProfileDelivery{}

	userID := "0"

	body := bytes.NewReader([]byte(fmt.Sprintf(`{"userId": "%s"`, userID)))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.Update(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateAvatarError(t *testing.T) {
	handler := ProfileDelivery{}

	userID := "0"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/avatars", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.UpdateAvatar(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
