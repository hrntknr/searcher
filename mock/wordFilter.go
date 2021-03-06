// Code generated by MockGen. DO NOT EDIT.
// Source: ../wordFilter.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWordFilter is a mock of WordFilter interface.
type MockWordFilter struct {
	ctrl     *gomock.Controller
	recorder *MockWordFilterMockRecorder
}

// MockWordFilterMockRecorder is the mock recorder for MockWordFilter.
type MockWordFilterMockRecorder struct {
	mock *MockWordFilter
}

// NewMockWordFilter creates a new mock instance.
func NewMockWordFilter(ctrl *gomock.Controller) *MockWordFilter {
	mock := &MockWordFilter{ctrl: ctrl}
	mock.recorder = &MockWordFilterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWordFilter) EXPECT() *MockWordFilterMockRecorder {
	return m.recorder
}

// Filter mocks base method.
func (m *MockWordFilter) Filter(arg0 [][]string) [][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Filter", arg0)
	ret0, _ := ret[0].([][]string)
	return ret0
}

// Filter indicates an expected call of Filter.
func (mr *MockWordFilterMockRecorder) Filter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockWordFilter)(nil).Filter), arg0)
}
