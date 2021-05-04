// Code generated by MockGen. DO NOT EDIT.
// Source: ../tokenizer.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTokenizer is a mock of Tokenizer interface.
type MockTokenizer struct {
	ctrl     *gomock.Controller
	recorder *MockTokenizerMockRecorder
}

// MockTokenizerMockRecorder is the mock recorder for MockTokenizer.
type MockTokenizerMockRecorder struct {
	mock *MockTokenizer
}

// NewMockTokenizer creates a new mock instance.
func NewMockTokenizer(ctrl *gomock.Controller) *MockTokenizer {
	mock := &MockTokenizer{ctrl: ctrl}
	mock.recorder = &MockTokenizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenizer) EXPECT() *MockTokenizerMockRecorder {
	return m.recorder
}

// Analyze mocks base method.
func (m *MockTokenizer) Analyze(text []string) [][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Analyze", text)
	ret0, _ := ret[0].([][]string)
	return ret0
}

// Analyze indicates an expected call of Analyze.
func (mr *MockTokenizerMockRecorder) Analyze(text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Analyze", reflect.TypeOf((*MockTokenizer)(nil).Analyze), text)
}