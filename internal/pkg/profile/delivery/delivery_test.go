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

func TestGetSuccess(t *testing.T) {
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

	handler := NewProfileDelivery(mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/profiles", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), prof.UserID)

	handler.Get(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respProf models.Profile
	_ = json.Unmarshal(w.Body.Bytes(), &respProf)
	if !reflect.DeepEqual(prof, respProf) {
		t.Errorf("expected: %v\n got: %v", prof, respProf)
	}
}

func TestGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	userID := "0"

	mockProfileUsecase.EXPECT().Get(userID).Times(1).Return(models.Profile{}, fmt.Errorf("error with profile"))

	handler := NewProfileDelivery(mockProfileUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/profiles", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.Get(w, r.WithContext(ctx))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusCreated, w.Code)
	}
}

func TestGetNoUser(t *testing.T) {
	handler := ProfileDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/profiles", nil)

	handler.Get(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := NewProfileDelivery(mockProfileUsecase)

	prof := models.Profile{
		UserID: "0",
		Name:   "testname",
		Phone:  "0000",
	}

	mockProfileUsecase.EXPECT().Update(prof).Times(1).Return(nil)

	profJson, _ := json.Marshal(&prof)
	body := bytes.NewReader(profJson)

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

	handler := NewProfileDelivery(mockProfileUsecase)

	prof := models.Profile{
		UserID: "0",
		Name:   "testname",
		Phone:  "0000",
	}

	mockProfileUsecase.EXPECT().Update(prof).Times(1).Return(fmt.Errorf("error"))

	profJson, _ := json.Marshal(&prof)
	body := bytes.NewReader(profJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), prof.UserID)

	handler.Update(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
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

func TestUpdateNoUser(t *testing.T) {
	handler := ProfileDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles", nil)

	handler.Update(w, r)

	expected := http.StatusInternalServerError
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

func TestUpdateAvatarNoUser(t *testing.T) {
	handler := ProfileDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/avatars", nil)

	handler.UpdateAvatar(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateAddressesSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := NewProfileDelivery(mockProfileUsecase)

	userID := "0"
	addresses := []string{"addr1", "addr2"}
	addrJson := fmt.Sprintf(`{"addresses": ["%s", "%s"]}`, addresses[0], addresses[1])

	mockProfileUsecase.EXPECT().UpdateAddresses(userID, addresses).Times(1).Return(nil)

	body := bytes.NewReader([]byte(addrJson))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/addresses", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.UpdateAddresses(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateAddressesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := NewProfileDelivery(mockProfileUsecase)

	userID := "0"
	addresses := []string{"addr1", "addr2"}
	addrJson := fmt.Sprintf(`{"addresses": ["%s", "%s"]}`, addresses[0], addresses[1])

	mockProfileUsecase.EXPECT().UpdateAddresses(userID, addresses).Times(1).Return(fmt.Errorf("err"))

	body := bytes.NewReader([]byte(addrJson))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/addresses", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.UpdateAddresses(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateAddressesBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileUsecase := profile.NewMockUsecase(ctrl)

	handler := NewProfileDelivery(mockProfileUsecase)

	userID := "0"
	addresses := []string{"addr1", "addr2"}
	addrJson := fmt.Sprintf(`{"addresses": ["%s", "%s"}`, addresses[0], addresses[1])

	body := bytes.NewReader([]byte(addrJson))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/addresses", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.UpdateAddresses(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateAddressesNoUser(t *testing.T) {
	handler := ProfileDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/profiles/addresses", nil)

	handler.UpdateAddresses(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
