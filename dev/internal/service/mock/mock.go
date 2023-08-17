// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	constant "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	model "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetNearestStores mocks base method.
func (m *MockService) GetNearestStores(ctx context.Context, req constant.GetNearestStoresRequest) ([]constant.NearestStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNearestStores", ctx, req)
	ret0, _ := ret[0].([]constant.NearestStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNearestStores indicates an expected call of GetNearestStores.
func (mr *MockServiceMockRecorder) GetNearestStores(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNearestStores", reflect.TypeOf((*MockService)(nil).GetNearestStores), ctx, req)
}

// SelectStore mocks base method.
func (m *MockService) SelectStore(ctx context.Context, token *model.UserToken, req constant.SelectStoreRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectStore", ctx, token, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectStore indicates an expected call of SelectStore.
func (mr *MockServiceMockRecorder) SelectStore(ctx, token, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectStore", reflect.TypeOf((*MockService)(nil).SelectStore), ctx, token, req)
}

// StoreState mocks base method.
func (m *MockService) StoreState(ctx context.Context, token *model.UserToken, req constant.StoreStateRequest) (*model.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreState", ctx, token, req)
	ret0, _ := ret[0].(*model.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreState indicates an expected call of StoreState.
func (mr *MockServiceMockRecorder) StoreState(ctx, token, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreState", reflect.TypeOf((*MockService)(nil).StoreState), ctx, token, req)
}