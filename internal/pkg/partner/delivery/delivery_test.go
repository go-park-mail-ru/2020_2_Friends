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
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/friends/internal/pkg/vendors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     2,
	}
	cookieName := &session.SessionName{Name: "sessname"}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Create(context.Background(), &session.UserID{Id: "0"}).Times(1).Return(cookieName, nil)

	handler := New(mockUserUsecase, mockProfileUsecase, mockSessionClient, mockVendorUsecase)

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

func TestCreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     2,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("", fmt.Errorf("db error"))

	handler := PartnerDelivery{
		userUsecase: mockUserUsecase,
	}

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/partners", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateProfileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     2,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(fmt.Errorf("db error"))

	handler := New(mockUserUsecase, mockProfileUsecase, mockSessionClient, mockVendorUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/partners", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateSessionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUsecase := user.NewMockUsecase(ctrl)
	mockProfileUsecase := profile.NewMockUsecase(ctrl)
	mockSessionClient := session.NewMockSessionWorkerClient(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	user := models.User{
		Login:    "testlogin",
		Password: "testpswd",
		Role:     2,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(nil)
	mockUserUsecase.EXPECT().Create(user).Times(1).Return("0", nil)
	mockProfileUsecase.EXPECT().Create("0").Times(1).Return(nil)
	mockSessionClient.EXPECT().Create(context.Background(), &session.UserID{Id: "0"}).Times(1).Return(nil, fmt.Errorf("db error"))

	handler := New(mockUserUsecase, mockProfileUsecase, mockSessionClient, mockVendorUsecase)

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/partners", body)

	handler.Create(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected: %v\n got: %v", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateBadJson(t *testing.T) {
	handler := PartnerDelivery{}

	body := bytes.NewReader([]byte(`{"login": "test", "password": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/partners", body)

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
		Role:     2,
	}

	mockUserUsecase.EXPECT().CheckIfUserExists(user).Times(1).Return(fmt.Errorf("user exists"))

	handler := PartnerDelivery{
		userUsecase: mockUserUsecase,
	}

	userJson, _ := json.Marshal(&user)
	body := bytes.NewReader(userJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/partners", body)

	handler.Create(w, r)

	expected := http.StatusConflict
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestCreateVendorSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"

	vendor := models.Vendor{
		Name:        "test",
		Description: "testtest",
		Picture:     "test.png",
	}

	vendorJson, _ := json.Marshal(&vendor)
	body := bytes.NewReader(vendorJson)

	mockVendorUsecase.EXPECT().Create(partnerID, vendor).Return(0, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.CreateVendor(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	expectedResp := models.AddResponse{
		ID: 0,
	}
	var resp models.AddResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(expectedResp, resp) {
		t.Errorf("expected: %v\n got: %v", expectedResp, resp)
	}
}

func TestCreateVendorNoPartner(t *testing.T) {
	hadnler := PartnerDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors", nil)

	hadnler.CreateVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestCreateVendorBadJson(t *testing.T) {
	handler := PartnerDelivery{}

	partnerID := "0"

	body := bytes.NewReader([]byte(`{"store_name": "test", "description": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.CreateVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestVendorDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"

	vendor := models.Vendor{
		Name:        "test",
		Description: "testtest",
		Picture:     "test.png",
	}

	vendorJson, _ := json.Marshal(&vendor)
	body := bytes.NewReader(vendorJson)

	mockVendorUsecase.EXPECT().Create(partnerID, vendor).Return(0, fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.CreateVendor(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateVendorSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "1"

	vendor := models.Vendor{
		Name:        "test",
		Description: "testtest",
		Picture:     "test.png",
	}

	vendorJson, _ := json.Marshal(&vendor)
	body := bytes.NewReader(vendorJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().Update(vendor).Times(1).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateVendor(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateVendorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "1"

	vendor := models.Vendor{
		Name:        "test",
		Description: "testtest",
		Picture:     "test.png",
	}

	vendorJson, _ := json.Marshal(&vendor)
	body := bytes.NewReader(vendorJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().Update(vendor).Times(1).Return(fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/1", body)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateVendor(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateVendorBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)

	body := bytes.NewReader([]byte(`{"store_name": "test", "description": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateVendorNotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(fmt.Errorf("err"))

	body := bytes.NewReader([]byte(`{"store_name": "test", "description": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateVendorNoPartnerID(t *testing.T) {
	handler := PartnerDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})

	handler.UpdateVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddProductSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"

	product := models.Product{
		Name:    "a",
		Price:   0,
		Picture: "c",
	}

	productJson, _ := json.Marshal(&product)
	body := bytes.NewReader(productJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().AddProduct(product).Times(1).Return(0, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.AddProductToVendor(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	expectedResp := models.AddResponse{
		ID: 0,
	}
	var resp models.AddResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(expectedResp, resp) {
		t.Errorf("expected: %v\n got: %v", expectedResp, resp)
	}
}

func TestAddProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"

	product := models.Product{
		Name:    "a",
		Price:   0,
		Picture: "c",
	}

	productJson, _ := json.Marshal(&product)
	body := bytes.NewReader(productJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().AddProduct(product).Times(1).Return(0, fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.AddProductToVendor(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddProductBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)

	body := bytes.NewReader([]byte(`{"food_name": "test", "food_price": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.AddProductToVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddProductNotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.AddProductToVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddProductNoPartnerID(t *testing.T) {
	handler := PartnerDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "0"})

	handler.AddProductToVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateProductSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"
	productID := "0"

	product := models.Product{
		Name:    "a",
		Price:   0,
		Picture: "c",
	}

	productJson, _ := json.Marshal(&product)
	body := bytes.NewReader(productJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().UpdateProduct(product).Times(1).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateProductOnVendor(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"
	productID := "0"

	product := models.Product{
		Name:    "a",
		Price:   0,
		Picture: "c",
	}

	productJson, _ := json.Marshal(&product)
	body := bytes.NewReader(productJson)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().UpdateProduct(product).Times(1).Return(fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateProductOnVendor(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateProductBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)

	body := bytes.NewReader([]byte(`{"food_name": "test", "food_price": "test"`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/vendors/0", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateProductOnVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateProductNotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID, vendorID := "0", "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.UpdateProductOnVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestUpdateProductNoPartnerID(t *testing.T) {
	handler := PartnerDelivery{}

	vendorID := "0"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})

	handler.UpdateProductOnVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteProductSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"
	productID := "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().DeleteProduct(productID).Times(1).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.DeleteProductFromVendor(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteProductError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"
	productID := "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(nil)
	mockVendorUsecase.EXPECT().DeleteProduct(productID).Times(1).Return(fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.DeleteProductFromVendor(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteProductNoPartner(t *testing.T) {
	handler := PartnerDelivery{}

	vendorID := "0"
	productID := "0"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})

	handler.DeleteProductFromVendor(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestDeleteProductNotOwned(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"
	vendorID := "0"
	productID := "0"

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, vendorID).Times(1).Return(fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/vendors/0", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": productID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.DeleteProductFromVendor(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetPartnerShopsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

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

	partnerID := "0"

	mockVendorUsecase.EXPECT().GetPartnerShops(partnerID).Times(1).Return(vendors, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/partners/vendors", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.GetPartnerShops(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respVendors []models.Vendor
	json.Unmarshal(w.Body.Bytes(), &respVendors)
	if !reflect.DeepEqual(vendors, respVendors) {
		t.Errorf("expected: %v\n got: %v", vendors, respVendors)
	}
}

func TestGetPartnerShopsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)
	handler := PartnerDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	partnerID := "0"

	mockVendorUsecase.EXPECT().GetPartnerShops(partnerID).Times(1).Return(nil, fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/partners/vendors", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler.GetPartnerShops(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetPartnerShopsNoPartner(t *testing.T) {
	handler := PartnerDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/partners/vendors", nil)

	handler.GetPartnerShops(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
