// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

package prstate

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return _m.recorder
}

// UpdateState mocks base method
func (_m *MockStorage) UpdateState(ctx context.Context, owner string, name string, analysisID string, state *State) error {
	ret := _m.ctrl.Call(_m, "UpdateState", ctx, owner, name, analysisID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateState indicates an expected call of UpdateState
func (_mr *MockStorageMockRecorder) UpdateState(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateState", reflect.TypeOf((*MockStorage)(nil).UpdateState), arg0, arg1, arg2, arg3, arg4)
}

// GetState mocks base method
func (_m *MockStorage) GetState(ctx context.Context, owner string, name string, analysisID string) (*State, error) {
	ret := _m.ctrl.Call(_m, "GetState", ctx, owner, name, analysisID)
	ret0, _ := ret[0].(*State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState
func (_mr *MockStorageMockRecorder) GetState(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetState", reflect.TypeOf((*MockStorage)(nil).GetState), arg0, arg1, arg2, arg3)
}
