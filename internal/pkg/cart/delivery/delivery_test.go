package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/cart"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/golang/mock/gomock"
)

func TestAddToCartSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"
	productID := "1"

	mockCartUsecase.EXPECT().Add(userID, productID).Times(1).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	q := r.URL.Query()
	q.Add(configs.ProductID, productID)
	r.URL.RawQuery = q.Encode()
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddToCart(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddToCartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"
	productID := "1"

	mockCartUsecase.EXPECT().Add(userID, productID).Times(1).Return(fmt.Errorf("error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	q := r.URL.Query()
	q.Add(configs.ProductID, productID)
	r.URL.RawQuery = q.Encode()
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddToCart(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddToCartNoQueryParam(t *testing.T) {
	handler := CartDelivery{}
	userID := "0"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddToCart(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddToCartNoUser(t *testing.T) {
	handler := CartDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)

	handler.AddToCart(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestRemoveFromCartSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"
	productID := "1"

	mockCartUsecase.EXPECT().Remove(userID, productID).Times(1).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	q := r.URL.Query()
	q.Add(configs.ProductID, productID)
	r.URL.RawQuery = q.Encode()
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.RemoveFromCart(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestRemoveFromCartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"
	productID := "1"

	mockCartUsecase.EXPECT().Remove(userID, productID).Times(1).Return(fmt.Errorf("error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	q := r.URL.Query()
	q.Add(configs.ProductID, productID)
	r.URL.RawQuery = q.Encode()
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.RemoveFromCart(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestRemoveFromCartNoQueryParam(t *testing.T) {
	handler := CartDelivery{}
	userID := "0"

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.RemoveFromCart(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestRemoveFromCartNoUser(t *testing.T) {
	handler := CartDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)

	handler.RemoveFromCart(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetCartSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"

	products := []models.Product{
		{
			Name:    "a",
			Price:   "b",
			Picture: "c",
		},
		{
			Name:    "c",
			Price:   "a",
			Picture: "b",
		},
	}

	mockCartUsecase.EXPECT().Get(userID).Times(1).Return(products, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.GetCart(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var respProducts []models.Product
	json.Unmarshal(w.Body.Bytes(), &respProducts)
	if !reflect.DeepEqual(products, respProducts) {
		t.Errorf("expected: %v\n got: %v", products, respProducts)
	}
}

func TestGetCartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartUsecase := cart.NewMockUsecase(ctrl)
	handler := NewCartDelivery(mockCartUsecase)

	userID := "0"

	mockCartUsecase.EXPECT().Get(userID).Times(1).Return(nil, fmt.Errorf("err"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.GetCart(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetCartNoUser(t *testing.T) {
	handler := CartDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/cart", nil)

	handler.GetCart(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
