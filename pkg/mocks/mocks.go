// Code generated by MockGen. DO NOT EDIT.
// Source: reapp_students_assignement/pkg (interfaces: DistanceInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDistanceInterface is a mock of DistanceInterface interface.
type MockDistanceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDistanceInterfaceMockRecorder
}

// MockDistanceInterfaceMockRecorder is the mock recorder for MockDistanceInterface.
type MockDistanceInterfaceMockRecorder struct {
	mock *MockDistanceInterface
}

// NewMockDistanceInterface creates a new mock instance.
func NewMockDistanceInterface(ctrl *gomock.Controller) *MockDistanceInterface {
	mock := &MockDistanceInterface{ctrl: ctrl}
	mock.recorder = &MockDistanceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDistanceInterface) EXPECT() *MockDistanceInterfaceMockRecorder {
	return m.recorder
}

// GetDistance mocks base method.
func (m *MockDistanceInterface) GetDistance(arg0, arg1 []string) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDistance", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetDistance indicates an expected call of GetDistance.
func (mr *MockDistanceInterfaceMockRecorder) GetDistance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDistance", reflect.TypeOf((*MockDistanceInterface)(nil).GetDistance), arg0, arg1)
}