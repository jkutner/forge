// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sclevine/forge/engine (interfaces: Container)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	engine "github.com/sclevine/forge/engine"
	io "io"
	reflect "reflect"
	time "time"
)

// MockContainer is a mock of Container interface
type MockContainer struct {
	ctrl     *gomock.Controller
	recorder *MockContainerMockRecorder
}

// MockContainerMockRecorder is the mock recorder for MockContainer
type MockContainerMockRecorder struct {
	mock *MockContainer
}

// NewMockContainer creates a new mock instance
func NewMockContainer(ctrl *gomock.Controller) *MockContainer {
	mock := &MockContainer{ctrl: ctrl}
	mock.recorder = &MockContainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContainer) EXPECT() *MockContainerMockRecorder {
	return m.recorder
}

// Background mocks base method
func (m *MockContainer) Background() error {
	ret := m.ctrl.Call(m, "Background")
	ret0, _ := ret[0].(error)
	return ret0
}

// Background indicates an expected call of Background
func (mr *MockContainerMockRecorder) Background() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Background", reflect.TypeOf((*MockContainer)(nil).Background))
}

// Close mocks base method
func (m *MockContainer) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockContainerMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockContainer)(nil).Close))
}

// CloseAfterStream mocks base method
func (m *MockContainer) CloseAfterStream(arg0 *engine.Stream) error {
	ret := m.ctrl.Call(m, "CloseAfterStream", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseAfterStream indicates an expected call of CloseAfterStream
func (mr *MockContainerMockRecorder) CloseAfterStream(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseAfterStream", reflect.TypeOf((*MockContainer)(nil).CloseAfterStream), arg0)
}

// Commit mocks base method
func (m *MockContainer) Commit(arg0 string) (string, error) {
	ret := m.ctrl.Call(m, "Commit", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Commit indicates an expected call of Commit
func (mr *MockContainerMockRecorder) Commit(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockContainer)(nil).Commit), arg0)
}

// ExtractTo mocks base method
func (m *MockContainer) ExtractTo(arg0 io.Reader, arg1 string) error {
	ret := m.ctrl.Call(m, "ExtractTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExtractTo indicates an expected call of ExtractTo
func (mr *MockContainerMockRecorder) ExtractTo(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractTo", reflect.TypeOf((*MockContainer)(nil).ExtractTo), arg0, arg1)
}

// HealthCheck mocks base method
func (m *MockContainer) HealthCheck() <-chan string {
	ret := m.ctrl.Call(m, "HealthCheck")
	ret0, _ := ret[0].(<-chan string)
	return ret0
}

// HealthCheck indicates an expected call of HealthCheck
func (mr *MockContainerMockRecorder) HealthCheck() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockContainer)(nil).HealthCheck))
}

// ID mocks base method
func (m *MockContainer) ID() string {
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockContainerMockRecorder) ID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockContainer)(nil).ID))
}

// Mkdir mocks base method
func (m *MockContainer) Mkdir(arg0 string) error {
	ret := m.ctrl.Call(m, "Mkdir", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mkdir indicates an expected call of Mkdir
func (mr *MockContainerMockRecorder) Mkdir(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mkdir", reflect.TypeOf((*MockContainer)(nil).Mkdir), arg0)
}

// Shell mocks base method
func (m *MockContainer) Shell(arg0 engine.TTY, arg1 ...string) error {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Shell", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shell indicates an expected call of Shell
func (mr *MockContainerMockRecorder) Shell(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shell", reflect.TypeOf((*MockContainer)(nil).Shell), varargs...)
}

// Start mocks base method
func (m *MockContainer) Start(arg0 string, arg1 io.Writer, arg2 <-chan time.Time) (int64, error) {
	ret := m.ctrl.Call(m, "Start", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start
func (mr *MockContainerMockRecorder) Start(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockContainer)(nil).Start), arg0, arg1, arg2)
}

// StreamFileFrom mocks base method
func (m *MockContainer) StreamFileFrom(arg0 string) (engine.Stream, error) {
	ret := m.ctrl.Call(m, "StreamFileFrom", arg0)
	ret0, _ := ret[0].(engine.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamFileFrom indicates an expected call of StreamFileFrom
func (mr *MockContainerMockRecorder) StreamFileFrom(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamFileFrom", reflect.TypeOf((*MockContainer)(nil).StreamFileFrom), arg0)
}

// StreamFileTo mocks base method
func (m *MockContainer) StreamFileTo(arg0 engine.Stream, arg1 string) error {
	ret := m.ctrl.Call(m, "StreamFileTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StreamFileTo indicates an expected call of StreamFileTo
func (mr *MockContainerMockRecorder) StreamFileTo(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamFileTo", reflect.TypeOf((*MockContainer)(nil).StreamFileTo), arg0, arg1)
}

// StreamTarFrom mocks base method
func (m *MockContainer) StreamTarFrom(arg0 string) (engine.Stream, error) {
	ret := m.ctrl.Call(m, "StreamTarFrom", arg0)
	ret0, _ := ret[0].(engine.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamTarFrom indicates an expected call of StreamTarFrom
func (mr *MockContainerMockRecorder) StreamTarFrom(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamTarFrom", reflect.TypeOf((*MockContainer)(nil).StreamTarFrom), arg0)
}

// StreamTarTo mocks base method
func (m *MockContainer) StreamTarTo(arg0 engine.Stream, arg1 string) error {
	ret := m.ctrl.Call(m, "StreamTarTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StreamTarTo indicates an expected call of StreamTarTo
func (mr *MockContainerMockRecorder) StreamTarTo(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamTarTo", reflect.TypeOf((*MockContainer)(nil).StreamTarTo), arg0, arg1)
}
