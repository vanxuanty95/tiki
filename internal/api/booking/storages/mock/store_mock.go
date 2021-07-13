// Code generated by MockGen. DO NOT EDIT.
// Source: tiki/internal/api/booking/storages (interfaces: Store)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storages "tiki/internal/api/booking/storages"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddBooking mocks base method.
func (m *MockStore) AddBooking(arg0 context.Context, arg1 *storages.Booking) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBooking", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBooking indicates an expected call of AddBooking.
func (mr *MockStoreMockRecorder) AddBooking(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBooking", reflect.TypeOf((*MockStore)(nil).AddBooking), arg0, arg1)
}

// RetrieveBookings mocks base method.
func (m *MockStore) RetrieveBookings(arg0 context.Context, arg1, arg2 string, arg3, arg4 int) ([]*storages.Booking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveBookings", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*storages.Booking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveBookings indicates an expected call of RetrieveBookings.
func (mr *MockStoreMockRecorder) RetrieveBookings(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveBookings", reflect.TypeOf((*MockStore)(nil).RetrieveBookings), arg0, arg1, arg2, arg3, arg4)
}