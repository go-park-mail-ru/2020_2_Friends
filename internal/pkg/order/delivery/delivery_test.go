package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
	websocketpool "github.com/friends/internal/pkg/websocketPool"
	ownErr "github.com/friends/pkg/error"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	dbError = fmt.Errorf("db error")

	testOrder = models.OrderRequest{
		ProductIDs: []int{1, 2, 3},
		Address:    "test addr",
	}

	response = models.OrderResponse{
		ID:           10,
		UserID:       50,
		VendorName:   "test",
		CreatedAtStr: time.Now().Format(configs.TimeFormat),
		Address:      "test addr",
		Status:       "ready",
		Price:        1000,
		Reviewed:     false,
		Products: []models.OrderProduct{
			{
				Name:    "test1",
				Price:   400,
				Picture: "1.jpg",
			},
			{
				Name:    "test2",
				Price:   600,
				Picture: "2.jpg",
			},
		},
	}

	vendorID = "15"

	testStatus = models.OrderStatusRequest{
		Status: "ready",
	}

	wsPool = websocketpool.NewWebsocketPool()
)

func TestAddOrderSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().AddOrder(strconv.Itoa(response.UserID), gomock.Any()).Return(response.ID, nil)

	orderJson, _ := json.Marshal(&testOrder)
	body := bytes.NewReader(orderJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/orders", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.AddOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}

	expectedResp := models.IDResponse{
		ID: response.ID,
	}
	var resp models.IDResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(expectedResp, resp) {
		t.Errorf("expected: %v\n got: %v", expectedResp, resp)
	}
}

func TestAddOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().AddOrder(strconv.Itoa(response.UserID), gomock.Any()).Return(0, ownErr.NewServerError(dbError))

	orderJson, _ := json.Marshal(&testOrder)
	body := bytes.NewReader(orderJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/orders", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.AddOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestAddOrderBadJson(t *testing.T) {
	body := bytes.NewReader([]byte(`{"id":0`))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/orders", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{}

	handler.AddOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestAddOrderNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/orders", nil)

	handler := OrderDelivery{}

	handler.AddOrder(w, r)

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetOrderSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetOrder(strconv.Itoa(response.UserID), strconv.Itoa(response.ID)).Return(response, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(response.ID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.GetOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}

	var resp models.OrderResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(response, resp) {
		t.Errorf("expected: %v\n got: %v", response, resp)
	}
}

func TestGetOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetOrder(strconv.Itoa(response.UserID), strconv.Itoa(response.ID)).Return(models.OrderResponse{}, ownErr.NewServerError(dbError))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(response.ID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.GetOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetOrderNoQueryParam(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders/10", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{}

	handler.GetOrder(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetOrderNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders/10", nil)

	handler := OrderDelivery{}

	handler.GetOrder(w, r)

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetUserOrdersSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	arrResp := []models.OrderResponse{response, response}
	mockOrderUsecase.EXPECT().GetUserOrders(strconv.Itoa(response.UserID)).Return(arrResp, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.GetUserOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}

	var resp []models.OrderResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(arrResp, resp) {
		t.Errorf("expected: %v\n got: %v", arrResp, resp)
	}
}

func TestGetUserOrdersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserOrders(strconv.Itoa(response.UserID)).Return([]models.OrderResponse{}, ownErr.NewServerError(dbError))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.GetUserOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetUserOrdersNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)

	handler := OrderDelivery{}

	handler.GetUserOrders(w, r)

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetVendorOrdersSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)
	vendorResp := models.VendorOrdersResponse{
		VendorName:    "test1",
		VendorPicture: "test.png",
		Orders:        []models.OrderResponse{response, response},
	}
	mockOrderUsecase.EXPECT().GetVendorOrders(vendorID).Times(1).Return(vendorResp, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := New(mockOrderUsecase, mockVendorUsecase, wsPool)

	handler.GetVendorOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}

	var resp models.VendorOrdersResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !reflect.DeepEqual(vendorResp, resp) {
		t.Errorf("expected: %v\n got: %v", vendorResp, resp)
	}
}

func TestGetVendorOrdersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)
	mockOrderUsecase.EXPECT().GetVendorOrders(vendorID).Times(1).Return(models.VendorOrdersResponse{}, ownErr.NewServerError(dbError))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := New(mockOrderUsecase, mockVendorUsecase, wsPool)

	handler.GetVendorOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetVendorOrdersNotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := New(mockOrderUsecase, mockVendorUsecase, wsPool)

	handler.GetVendorOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetVendorOrdersNoQueryParams(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{}

	handler.GetVendorOrders(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestGetVendorOrdersNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)

	handler := OrderDelivery{}

	handler.GetVendorOrders(w, r)

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)
	mockOrderUsecase.EXPECT().UpdateOrderStatus(strconv.Itoa(response.ID), testStatus.Status).Times(1).Return(nil)
	mockOrderUsecase.EXPECT().GetUserIDFromOrder(response.ID).Times(1).Return("0", nil)

	statusJson, _ := json.Marshal(&testStatus)
	body := bytes.NewReader(statusJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": strconv.Itoa(response.ID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := New(mockOrderUsecase, mockVendorUsecase, wsPool)

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)
	mockOrderUsecase.EXPECT().UpdateOrderStatus(strconv.Itoa(response.ID), testStatus.Status).Times(1).Return(dbError)

	statusJson, _ := json.Marshal(&testStatus)
	body := bytes.NewReader(statusJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID, "id": strconv.Itoa(response.ID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := New(mockOrderUsecase, mockVendorUsecase, wsPool)

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusNoOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)

	statusJson, _ := json.Marshal(&testStatus)
	body := bytes.NewReader(statusJson)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(nil)

	body := bytes.NewReader([]byte(`{"status":"ready`))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", body)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusNotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(strconv.Itoa(response.UserID), vendorID).Times(1).Return(dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	r = mux.SetURLVars(r, map[string]string{"vendorID": vendorID})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusNoQueryParamVendorID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), strconv.Itoa(response.UserID))

	handler := OrderDelivery{}

	handler.UpdateOrderStatus(w, r.WithContext(ctx))

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}

func TestUpdateOrderStatusNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/orders", nil)

	handler := OrderDelivery{}

	handler.UpdateOrderStatus(w, r)

	expectedCode := http.StatusInternalServerError
	if w.Code != expectedCode {
		t.Errorf("expected: %v\n got: %v", expectedCode, w.Code)
	}
}
