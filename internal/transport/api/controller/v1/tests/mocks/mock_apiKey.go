// Code generated by MockGen. DO NOT EDIT.
// Source: apiKey.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/linnoxlewis/trade-bot/internal/domain"
	dto "github.com/linnoxlewis/trade-bot/internal/domain/dto"
)

// MockApiKeyService is a mock of ApiKeyService interface.
type MockApiKeyService struct {
	ctrl     *gomock.Controller
	recorder *MockApiKeyServiceMockRecorder
}

// MockApiKeyServiceMockRecorder is the mock recorder for MockApiKeyService.
type MockApiKeyServiceMockRecorder struct {
	mock *MockApiKeyService
}

// NewMockApiKeyService creates a new mock instance.
func NewMockApiKeyService(ctrl *gomock.Controller) *MockApiKeyService {
	mock := &MockApiKeyService{ctrl: ctrl}
	mock.recorder = &MockApiKeyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiKeyService) EXPECT() *MockApiKeyServiceMockRecorder {
	return m.recorder
}

// AddApiKeys mocks base method.
func (m *MockApiKeyService) AddApiKeys(ctx context.Context, apiKeys *dto.ApiKeys) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddApiKeys", ctx, apiKeys)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddApiKeys indicates an expected call of AddApiKeys.
func (mr *MockApiKeyServiceMockRecorder) AddApiKeys(ctx, apiKeys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddApiKeys", reflect.TypeOf((*MockApiKeyService)(nil).AddApiKeys), ctx, apiKeys)
}

// ClearApiKey mocks base method.
func (m *MockApiKeyService) ClearApiKey(ctx context.Context, userId int64, exchange string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearApiKey", ctx, userId, exchange)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearApiKey indicates an expected call of ClearApiKey.
func (mr *MockApiKeyServiceMockRecorder) ClearApiKey(ctx, userId, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearApiKey", reflect.TypeOf((*MockApiKeyService)(nil).ClearApiKey), ctx, userId, exchange)
}

// DeleteApiKey mocks base method.
func (m *MockApiKeyService) DeleteApiKey(ctx context.Context, userId int64, exchange string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteApiKey", ctx, userId, exchange)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApiKey indicates an expected call of DeleteApiKey.
func (mr *MockApiKeyServiceMockRecorder) DeleteApiKey(ctx, userId, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApiKey", reflect.TypeOf((*MockApiKeyService)(nil).DeleteApiKey), ctx, userId, exchange)
}

// GetApiKeyByExchangeAndId mocks base method.
func (m *MockApiKeyService) GetApiKeyByExchangeAndId(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApiKeyByExchangeAndId", ctx, userId, exchange)
	ret0, _ := ret[0].(*domain.ApiKeys)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApiKeyByExchangeAndId indicates an expected call of GetApiKeyByExchangeAndId.
func (mr *MockApiKeyServiceMockRecorder) GetApiKeyByExchangeAndId(ctx, userId, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApiKeyByExchangeAndId", reflect.TypeOf((*MockApiKeyService)(nil).GetApiKeyByExchangeAndId), ctx, userId, exchange)
}
