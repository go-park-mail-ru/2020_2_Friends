package delivery

import (
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
	"github.com/friends/internal/pkg/chat"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	userID    = "10"
	partnerID = "11"
	vendorID  = 14
	orderID   = 50

	testMsgs = []models.Message{
		{
			OrderID:   orderID,
			IsYourMsg: true,
			Text:      "test1",
			SentAtStr: time.Now().Format(configs.TimeFormat),
		},
		{
			OrderID:   orderID,
			IsYourMsg: true,
			Text:      "test2",
			SentAtStr: time.Now().Format(configs.TimeFormat),
		},
	}

	testChats = []models.Chat{
		{
			OrderID:          1,
			InterlocutorID:   userID,
			InterlocutorName: "test",
			LastMsg:          "test1",
		},
		{
			OrderID:          2,
			InterlocutorID:   userID,
			InterlocutorName: "test",
			LastMsg:          "test2",
		},
	}

	dbError = fmt.Errorf("db error")
)

func TestGetChatSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUsecase := chat.NewMockUsecase(ctrl)
	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return(userID, nil)
	mockChatUsecase.EXPECT().GetChat(orderID, userID).Times(1).Return(testMsgs, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler := ChatDelivery{
		chatUsecase:  mockChatUsecase,
		orderUsecase: mockOrderUsecase,
	}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var msgs []models.Message
	_ = json.Unmarshal(w.Body.Bytes(), &msgs)
	if !reflect.DeepEqual(testMsgs, msgs) {
		t.Errorf("expected: %v\n got: %v", testMsgs, msgs)
	}
}

func TestGetChatSuccess2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUsecase := chat.NewMockUsecase(ctrl)
	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return(userID, nil)
	mockOrderUsecase.EXPECT().GetVendorIDFromOrder(orderID).Times(1).Return(vendorID, nil)
	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, strconv.Itoa(vendorID)).Times(1).Return(nil)
	mockChatUsecase.EXPECT().GetChat(orderID, partnerID).Times(1).Return(testMsgs, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := New(mockChatUsecase, mockOrderUsecase, mockVendorUsecase)

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var msgs []models.Message
	_ = json.Unmarshal(w.Body.Bytes(), &msgs)
	if !reflect.DeepEqual(testMsgs, msgs) {
		t.Errorf("expected: %v\n got: %v", testMsgs, msgs)
	}
}

func TestGetChatCheckVendorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return(userID, nil)
	mockOrderUsecase.EXPECT().GetVendorIDFromOrder(orderID).Times(1).Return(vendorID, nil)
	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, strconv.Itoa(vendorID)).Times(1).Return(dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{
		orderUsecase:  mockOrderUsecase,
		vendorUsecase: mockVendorUsecase,
	}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatGetVendorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return(userID, nil)
	mockOrderUsecase.EXPECT().GetVendorIDFromOrder(orderID).Times(1).Return(0, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{
		orderUsecase:  mockOrderUsecase,
		vendorUsecase: mockVendorUsecase,
	}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatGetChatError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUsecase := chat.NewMockUsecase(ctrl)
	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return(userID, nil)
	mockChatUsecase.EXPECT().GetChat(orderID, userID).Times(1).Return(nil, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler := ChatDelivery{
		chatUsecase:  mockChatUsecase,
		orderUsecase: mockOrderUsecase,
	}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatGetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderUsecase := order.NewMockUsecase(ctrl)

	mockOrderUsecase.EXPECT().GetUserIDFromOrder(orderID).Times(1).Return("", dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(orderID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler := ChatDelivery{
		orderUsecase: mockOrderUsecase,
	}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatWrongQuery(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "aaa"})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler := ChatDelivery{}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatNoQuery(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), userID)

	handler := ChatDelivery{}

	handler.GetChat(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetChatNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)

	handler := ChatDelivery{}

	handler.GetChat(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorChats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUsecase := chat.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, strconv.Itoa(vendorID)).Times(1).Return(nil)
	mockChatUsecase.EXPECT().GetVendorChats(strconv.Itoa(vendorID)).Times(1).Return(testChats, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(vendorID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{
		chatUsecase:   mockChatUsecase,
		vendorUsecase: mockVendorUsecase,
	}

	handler.GetVendorChats(w, r.WithContext(ctx))

	expected := http.StatusOK
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}

	var chats []models.Chat
	_ = json.Unmarshal(w.Body.Bytes(), &chats)
	if !reflect.DeepEqual(testChats, chats) {
		t.Errorf("expected: %v\n got: %v", testChats, chats)
	}
}

func TestGetVendorChatsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUsecase := chat.NewMockUsecase(ctrl)
	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, strconv.Itoa(vendorID)).Times(1).Return(nil)
	mockChatUsecase.EXPECT().GetVendorChats(strconv.Itoa(vendorID)).Times(1).Return(nil, dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(vendorID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{
		chatUsecase:   mockChatUsecase,
		vendorUsecase: mockVendorUsecase,
	}

	handler.GetVendorChats(w, r.WithContext(ctx))

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorChatsOwnerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorUsecase := vendors.NewMockUsecase(ctrl)

	mockVendorUsecase.EXPECT().CheckVendorOwner(partnerID, strconv.Itoa(vendorID)).Times(1).Return(dbError)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(vendorID)})
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{
		vendorUsecase: mockVendorUsecase,
	}

	handler.GetVendorChats(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorChatsNoQueryParams(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)
	ctx := context.WithValue(r.Context(), middleware.UserID(configs.UserID), partnerID)

	handler := ChatDelivery{}

	handler.GetVendorChats(w, r.WithContext(ctx))

	expected := http.StatusBadRequest
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}

func TestGetVendorChatsNoUser(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/chats", nil)

	handler := ChatDelivery{}

	handler.GetVendorChats(w, r)

	expected := http.StatusInternalServerError
	if w.Code != expected {
		t.Errorf("expected: %v\n got: %v", expected, w.Code)
	}
}
