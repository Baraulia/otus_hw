// Code generated by MockGen. DO NOT EDIT.
// Source: routers.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"
	time "time"

	models "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockApplicationInterface is a mock of ApplicationInterface interface.
type MockApplicationInterface struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationInterfaceMockRecorder
}

// MockApplicationInterfaceMockRecorder is the mock recorder for MockApplicationInterface.
type MockApplicationInterfaceMockRecorder struct {
	mock *MockApplicationInterface
}

// NewMockApplicationInterface creates a new mock instance.
func NewMockApplicationInterface(ctrl *gomock.Controller) *MockApplicationInterface {
	mock := &MockApplicationInterface{ctrl: ctrl}
	mock.recorder = &MockApplicationInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationInterface) EXPECT() *MockApplicationInterfaceMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockApplicationInterface) CreateEvent(ctx context.Context, eventDTO models.Event) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, eventDTO)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockApplicationInterfaceMockRecorder) CreateEvent(ctx, eventDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockApplicationInterface)(nil).CreateEvent), ctx, eventDTO)
}

// DeleteEvent mocks base method.
func (m *MockApplicationInterface) DeleteEvent(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockApplicationInterfaceMockRecorder) DeleteEvent(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockApplicationInterface)(nil).DeleteEvent), ctx, id)
}

// GetListEventsDuringDay mocks base method.
func (m *MockApplicationInterface) GetListEventsDuringDay(ctx context.Context, day time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListEventsDuringDay", ctx, day)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListEventsDuringDay indicates an expected call of GetListEventsDuringDay.
func (mr *MockApplicationInterfaceMockRecorder) GetListEventsDuringDay(ctx, day interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListEventsDuringDay", reflect.TypeOf((*MockApplicationInterface)(nil).GetListEventsDuringDay), ctx, day)
}

// GetListEventsDuringFewDays mocks base method.
func (m *MockApplicationInterface) GetListEventsDuringFewDays(ctx context.Context, start time.Time, amountDays int) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListEventsDuringFewDays", ctx, start, amountDays)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListEventsDuringFewDays indicates an expected call of GetListEventsDuringFewDays.
func (mr *MockApplicationInterfaceMockRecorder) GetListEventsDuringFewDays(ctx, start, amountDays interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListEventsDuringFewDays", reflect.TypeOf((*MockApplicationInterface)(nil).GetListEventsDuringFewDays), ctx, start, amountDays)
}

// UpdateEvent mocks base method.
func (m *MockApplicationInterface) UpdateEvent(ctx context.Context, eventDTO models.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, eventDTO)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockApplicationInterfaceMockRecorder) UpdateEvent(ctx, eventDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockApplicationInterface)(nil).UpdateEvent), ctx, eventDTO)
}
