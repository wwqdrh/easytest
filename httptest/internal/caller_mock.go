// Code generated by MockGen. DO NOT EDIT.
// Source: caller.go

// Package internal is a generated GoMock package.
package internal

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIHTTPCtx is a mock of IHTTPCtx interface.
type MockIHTTPCtx struct {
	ctrl     *gomock.Controller
	recorder *MockIHTTPCtxMockRecorder
}

// MockIHTTPCtxMockRecorder is the mock recorder for MockIHTTPCtx.
type MockIHTTPCtxMockRecorder struct {
	mock *MockIHTTPCtx
}

// NewMockIHTTPCtx creates a new mock instance.
func NewMockIHTTPCtx(ctrl *gomock.Controller) *MockIHTTPCtx {
	mock := &MockIHTTPCtx{ctrl: ctrl}
	mock.recorder = &MockIHTTPCtxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHTTPCtx) EXPECT() *MockIHTTPCtxMockRecorder {
	return m.recorder
}

// GetEnv mocks base method.
func (m *MockIHTTPCtx) GetEnv(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnv", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetEnv indicates an expected call of GetEnv.
func (mr *MockIHTTPCtxMockRecorder) GetEnv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnv", reflect.TypeOf((*MockIHTTPCtx)(nil).GetEnv), arg0)
}

// GetRequest mocks base method.
func (m *MockIHTTPCtx) GetRequest() *http.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRequest")
	ret0, _ := ret[0].(*http.Request)
	return ret0
}

// GetRequest indicates an expected call of GetRequest.
func (mr *MockIHTTPCtxMockRecorder) GetRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequest", reflect.TypeOf((*MockIHTTPCtx)(nil).GetRequest))
}

// GetResponse mocks base method.
func (m *MockIHTTPCtx) GetResponse() *http.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResponse")
	ret0, _ := ret[0].(*http.Response)
	return ret0
}

// GetResponse indicates an expected call of GetResponse.
func (mr *MockIHTTPCtxMockRecorder) GetResponse() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResponse", reflect.TypeOf((*MockIHTTPCtx)(nil).GetResponse))
}

// SetEnv mocks base method.
func (m *MockIHTTPCtx) SetEnv(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetEnv", arg0, arg1)
}

// SetEnv indicates an expected call of SetEnv.
func (mr *MockIHTTPCtxMockRecorder) SetEnv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEnv", reflect.TypeOf((*MockIHTTPCtx)(nil).SetEnv), arg0, arg1)
}

// MockIInstance is a mock of IInstance interface.
type MockIInstance struct {
	ctrl     *gomock.Controller
	recorder *MockIInstanceMockRecorder
}

// MockIInstanceMockRecorder is the mock recorder for MockIInstance.
type MockIInstanceMockRecorder struct {
	mock *MockIInstance
}

// NewMockIInstance creates a new mock instance.
func NewMockIInstance(ctrl *gomock.Controller) *MockIInstance {
	mock := &MockIInstance{ctrl: ctrl}
	mock.recorder = &MockIInstanceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIInstance) EXPECT() *MockIInstanceMockRecorder {
	return m.recorder
}

// GetAttr mocks base method.
func (m *MockIInstance) GetAttr(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttr", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetAttr indicates an expected call of GetAttr.
func (mr *MockIInstanceMockRecorder) GetAttr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttr", reflect.TypeOf((*MockIInstance)(nil).GetAttr), arg0)
}
