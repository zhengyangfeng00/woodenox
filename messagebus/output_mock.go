// Code generated by MockGen. DO NOT EDIT.
// Source: output.go

// Package messagebus is a generated GoMock package.
package messagebus

import (
	gomock "github.com/golang/mock/gomock"
	iterator "github.com/zhengyangfeng00/woodenox/iterator"
	reflect "reflect"
)

// MockOutput is a mock of Output interface
type MockOutput struct {
	ctrl     *gomock.Controller
	recorder *MockOutputMockRecorder
}

// MockOutputMockRecorder is the mock recorder for MockOutput
type MockOutputMockRecorder struct {
	mock *MockOutput
}

// NewMockOutput creates a new mock instance
func NewMockOutput(ctrl *gomock.Controller) *MockOutput {
	mock := &MockOutput{ctrl: ctrl}
	mock.recorder = &MockOutputMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOutput) EXPECT() *MockOutputMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockOutput) Write(item iterator.Item) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Write", item)
}

// Write indicates an expected call of Write
func (mr *MockOutputMockRecorder) Write(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockOutput)(nil).Write), item)
}
