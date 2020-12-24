package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/vendors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	testVendors = []models.Vendor{
		{
			ID:          0,
			Name:        "a",
			Description: "b",
			Picture:     "c",
		},
		{
			ID:          1,
			Name:        "c",
			Description: "a",
			Picture:     "b",
		},
	}

	Longitude = "75.4"
	Latitude  = "37.2"

	dbError = fmt.Errorf("db error")
)

func TestGetVendorSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	vendor := models.Vendor{
		ID:          0,
		Name:        "a",
		Description: "b",
		Picture:     "c",
	}

	mockVendorUsecase.EXPECT().Get(vendor.ID).Times(1).Return(vendor, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(vendor.ID)})

	handler.GetVendor(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respVendor models.Vendor
	_ = json.Unmarshal(w.Body.Bytes(), &respVendor)
	if !reflect.DeepEqual(vendor, respVendor) {
		t.Errorf("expected: %v\n got: %v", vendor, respVendor)
	}
}

func TestGetVendorDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	mockVendorUsecase.EXPECT().Get(0).Times(1).Return(models.Vendor{}, fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})

	handler.GetVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorBadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/a", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "a"})

	handler.GetVendor(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorNoID(t *testing.T) {
	handler := VendorDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/a", nil)

	handler.GetVendor(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetAllSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	vendors := []models.Vendor{
		{
			ID:          0,
			Name:        "a",
			Description: "b",
			Picture:     "c",
		},
		{
			ID:          1,
			Name:        "c",
			Description: "a",
			Picture:     "b",
		},
	}

	mockVendorUsecase.EXPECT().GetAll().Times(1).Return(vendors, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)

	handler.GetAll(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respVendors []models.Vendor
	_ = json.Unmarshal(w.Body.Bytes(), &respVendors)
	if !reflect.DeepEqual(vendors, respVendors) {
		t.Errorf("expected: %v\n got: %v", vendors, respVendors)
	}
}

func TestGetAllDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	mockVendorUsecase.EXPECT().GetAll().Times(1).Return(nil, fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)

	handler.GetAll(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetNearestSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	mockVendorUsecase.EXPECT().GetNearest(gomock.Any(), gomock.Any()).Times(1).Return(testVendors, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	q.Add(configs.Longitude, Longitude)
	q.Add(configs.Latitude, Latitude)
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respVendors []models.Vendor
	_ = json.Unmarshal(w.Body.Bytes(), &respVendors)
	if !reflect.DeepEqual(testVendors, respVendors) {
		t.Errorf("expected: %v\n got: %v", testVendors, respVendors)
	}
}

func TestGetNearestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := NewVendorDelivery(mockVendorUsecase)

	mockVendorUsecase.EXPECT().GetNearest(gomock.Any(), gomock.Any()).Times(1).Return(nil, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	q.Add(configs.Longitude, Longitude)
	q.Add(configs.Latitude, Latitude)
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetNearestWrongQuery(t *testing.T) {
	handler := VendorDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	q.Add(configs.Longitude, Longitude)
	q.Add(configs.Latitude, "aaa")
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetNearestWrongQuery2(t *testing.T) {
	handler := VendorDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	q.Add(configs.Longitude, "aaa")
	q.Add(configs.Latitude, Latitude)
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetNearestNoQuery(t *testing.T) {
	handler := VendorDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	q.Add(configs.Longitude, Longitude)
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
func TestGetNearestNoQuery2(t *testing.T) {
	handler := VendorDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors", nil)
	q := r.URL.Query()
	r.URL.RawQuery = q.Encode()

	handler.GetNearest(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
