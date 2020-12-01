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
	"github.com/friends/internal/pkg/review"
	ownErr "github.com/friends/pkg/error"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var userID = "0"
var vendorID = "0"

var testReview = models.Review{
	VendorID:     0,
	OrderID:      1,
	Rating:       5,
	Text:         "test",
	CreatedAtStr: "10.04.2020 12:42:19",
}

var dbError = ownErr.NewClientError(fmt.Errorf("db error"))

func TestAddReviewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	mockReviewUsecase.EXPECT().AddReview(gomock.Any()).Times(1).Return(nil)

	reviewJson, _ := json.Marshal(&testReview)
	body := bytes.NewReader(reviewJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/reviews", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddReview(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddReviewBadJson(t *testing.T) {
	handler := ReviewDelivery{}

	reviewJson := []byte(`{"vendo_id: 0`)
	body := bytes.NewReader(reviewJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/reviews", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddReview(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddReviewError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	mockReviewUsecase.EXPECT().AddReview(gomock.Any()).Times(1).Return(dbError)

	reviewJson, _ := json.Marshal(&testReview)
	body := bytes.NewReader(reviewJson)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/reviews", body)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.AddReview(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestAddReviewNoUser(t *testing.T) {
	handler := ReviewDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/reviews", nil)

	handler.AddReview(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetUserReviewsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	mockReviewUsecase.EXPECT().GetUserReviews(userID).Times(1).Return([]models.Review{testReview}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/reviews", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.GetUserReviews(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var reviewsResp []models.Review
	_ = json.Unmarshal(w.Body.Bytes(), &reviewsResp)
	expectedReviews := []models.Review{testReview}
	if !reflect.DeepEqual(expectedReviews, reviewsResp) {
		t.Errorf("expected: %v\n got: %v", expectedReviews, reviewsResp)
	}
}

func TestGetUserReviewsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	mockReviewUsecase.EXPECT().GetUserReviews(userID).Times(1).Return(nil, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/reviews", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler.GetUserReviews(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetUserReviewsNoUser(t *testing.T) {
	handler := ReviewDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/reviews", nil)

	handler.GetUserReviews(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorReviewsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	expectedReviews := models.VendorReviewsResponse{
		VendorName:    "test",
		VendorPicture: "test.png",
		Reviews:       []models.Review{testReview},
	}

	mockReviewUsecase.EXPECT().GetVendorReviews(vendorID).Times(1).Return(expectedReviews, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/0/reviews", nil)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})

	handler.GetVendorReviews(w, r)

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var reviewsResp models.VendorReviewsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &reviewsResp)
	if !reflect.DeepEqual(expectedReviews, reviewsResp) {
		t.Errorf("expected: %v\n got: %v", expectedReviews, reviewsResp)
	}
}

func TestGetVendorReviewsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewUsecase := review.NewMockUsecase(ctrl)

	handler := New(mockReviewUsecase)

	mockReviewUsecase.EXPECT().GetVendorReviews(vendorID).Times(1).Return(models.VendorReviewsResponse{}, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/0/reviews", nil)
	r = mux.SetURLVars(r, map[string]string{"id": vendorID})

	handler.GetVendorReviews(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorReviewsNoVendorID(t *testing.T) {

	handler := ReviewDelivery{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/vendors/0/reviews", nil)

	handler.GetVendorReviews(w, r)

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
