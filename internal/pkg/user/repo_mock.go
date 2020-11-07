// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/friends/internal/pkg/user (interfaces: Repository)

// Package user is a generated GoMock package.
package user

import (
	models "github.com/friends/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CheckIfUserExists mocks base method
func (m *MockRepository) CheckIfUserExists(arg0 models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfUserExists", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckIfUserExists indicates an expected call of CheckIfUserExists
func (mr *MockRepositoryMockRecorder) CheckIfUserExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfUserExists", reflect.TypeOf((*MockRepository)(nil).CheckIfUserExists), arg0)
}

// CheckLoginAndPassword mocks base method
func (m *MockRepository) CheckLoginAndPassword(arg0 models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLoginAndPassword", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLoginAndPassword indicates an expected call of CheckLoginAndPassword
func (mr *MockRepositoryMockRecorder) CheckLoginAndPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLoginAndPassword", reflect.TypeOf((*MockRepository)(nil).CheckLoginAndPassword), arg0)
}

// CheckUsersRole mocks base method
func (m *MockRepository) CheckUsersRole(arg0 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUsersRole", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUsersRole indicates an expected call of CheckUsersRole
func (mr *MockRepositoryMockRecorder) CheckUsersRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUsersRole", reflect.TypeOf((*MockRepository)(nil).CheckUsersRole), arg0)
}

// Create mocks base method
func (m *MockRepository) Create(arg0 models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRepositoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockRepository) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), arg0)
}
