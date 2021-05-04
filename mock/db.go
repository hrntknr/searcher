// Code generated by MockGen. DO NOT EDIT.
// Source: ../db.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/hrntknr/searcher/types"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// CountDocument mocks base method.
func (m *MockDB) CountDocument() (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountDocument")
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountDocument indicates an expected call of CountDocument.
func (mr *MockDBMockRecorder) CountDocument() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountDocument", reflect.TypeOf((*MockDB)(nil).CountDocument))
}

// CountTermInDocument mocks base method.
func (m *MockDB) CountTermInDocument(documentID uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountTermInDocument", documentID)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountTermInDocument indicates an expected call of CountTermInDocument.
func (mr *MockDBMockRecorder) CountTermInDocument(documentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountTermInDocument", reflect.TypeOf((*MockDB)(nil).CountTermInDocument), documentID)
}

// CreateDcoument mocks base method.
func (m *MockDB) CreateDcoument(document *types.Document) (*types.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDcoument", document)
	ret0, _ := ret[0].(*types.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDcoument indicates an expected call of CreateDcoument.
func (mr *MockDBMockRecorder) CreateDcoument(document interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDcoument", reflect.TypeOf((*MockDB)(nil).CreateDcoument), document)
}

// CreatePosting mocks base method.
func (m *MockDB) CreatePosting(posting *types.Posting) (*types.Posting, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePosting", posting)
	ret0, _ := ret[0].(*types.Posting)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePosting indicates an expected call of CreatePosting.
func (mr *MockDBMockRecorder) CreatePosting(posting interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePosting", reflect.TypeOf((*MockDB)(nil).CreatePosting), posting)
}

// CreateSentence mocks base method.
func (m *MockDB) CreateSentence(sentence *types.Sentence) (*types.Sentence, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSentence", sentence)
	ret0, _ := ret[0].(*types.Sentence)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSentence indicates an expected call of CreateSentence.
func (mr *MockDBMockRecorder) CreateSentence(sentence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSentence", reflect.TypeOf((*MockDB)(nil).CreateSentence), sentence)
}

// CreateToken mocks base method.
func (m *MockDB) CreateToken(token *types.Token) (*types.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", token)
	ret0, _ := ret[0].(*types.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockDBMockRecorder) CreateToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockDB)(nil).CreateToken), token)
}

// DeleteSentenceFromDocumentID mocks base method.
func (m *MockDB) DeleteSentenceFromDocumentID(documentID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSentenceFromDocumentID", documentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSentenceFromDocumentID indicates an expected call of DeleteSentenceFromDocumentID.
func (mr *MockDBMockRecorder) DeleteSentenceFromDocumentID(documentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSentenceFromDocumentID", reflect.TypeOf((*MockDB)(nil).DeleteSentenceFromDocumentID), documentID)
}

// DocumentFromID mocks base method.
func (m *MockDB) DocumentFromID(id uint) (*types.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DocumentFromID", id)
	ret0, _ := ret[0].(*types.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DocumentFromID indicates an expected call of DocumentFromID.
func (mr *MockDBMockRecorder) DocumentFromID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DocumentFromID", reflect.TypeOf((*MockDB)(nil).DocumentFromID), id)
}

// DocumentFromUri mocks base method.
func (m *MockDB) DocumentFromUri(uri string) (*types.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DocumentFromUri", uri)
	ret0, _ := ret[0].(*types.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DocumentFromUri indicates an expected call of DocumentFromUri.
func (mr *MockDBMockRecorder) DocumentFromUri(uri interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DocumentFromUri", reflect.TypeOf((*MockDB)(nil).DocumentFromUri), uri)
}

// PostingList mocks base method.
func (m *MockDB) PostingList(tokenID uint) ([]*types.Posting, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostingList", tokenID)
	ret0, _ := ret[0].([]*types.Posting)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostingList indicates an expected call of PostingList.
func (mr *MockDBMockRecorder) PostingList(tokenID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostingList", reflect.TypeOf((*MockDB)(nil).PostingList), tokenID)
}

// SentenceMultiFromID mocks base method.
func (m *MockDB) SentenceMultiFromID(ids []uint) ([]*types.Sentence, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SentenceMultiFromID", ids)
	ret0, _ := ret[0].([]*types.Sentence)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SentenceMultiFromID indicates an expected call of SentenceMultiFromID.
func (mr *MockDBMockRecorder) SentenceMultiFromID(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentenceMultiFromID", reflect.TypeOf((*MockDB)(nil).SentenceMultiFromID), ids)
}

// TokenFromID mocks base method.
func (m *MockDB) TokenFromID(id uint) (*types.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokenFromID", id)
	ret0, _ := ret[0].(*types.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TokenFromID indicates an expected call of TokenFromID.
func (mr *MockDBMockRecorder) TokenFromID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokenFromID", reflect.TypeOf((*MockDB)(nil).TokenFromID), id)
}

// TokenFromString mocks base method.
func (m *MockDB) TokenFromString(token string) (*types.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokenFromString", token)
	ret0, _ := ret[0].(*types.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TokenFromString indicates an expected call of TokenFromString.
func (mr *MockDBMockRecorder) TokenFromString(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokenFromString", reflect.TypeOf((*MockDB)(nil).TokenFromString), token)
}
