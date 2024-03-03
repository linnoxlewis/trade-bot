// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/linnoxlewis/trade-bot/internal/domain"
	service "github.com/linnoxlewis/trade-bot/internal/service"
)

// MockOrderRepo is a mock of OrderRepo interface.
type MockOrderRepo struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepoMockRecorder
}

// MockOrderRepoMockRecorder is the mock recorder for MockOrderRepo.
type MockOrderRepoMockRecorder struct {
	mock *MockOrderRepo
}

// NewMockOrderRepo creates a new mock instance.
func NewMockOrderRepo(ctrl *gomock.Controller) *MockOrderRepo {
	mock := &MockOrderRepo{ctrl: ctrl}
	mock.recorder = &MockOrderRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepo) EXPECT() *MockOrderRepoMockRecorder {
	return m.recorder
}

// ActivateOrder mocks base method.
func (m *MockOrderRepo) ActivateOrder(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateOrder indicates an expected call of ActivateOrder.
func (mr *MockOrderRepoMockRecorder) ActivateOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateOrder", reflect.TypeOf((*MockOrderRepo)(nil).ActivateOrder), ctx, id)
}

// Atomic mocks base method.
func (m *MockOrderRepo) Atomic(ctx context.Context, fn func(context.Context, service.OrderRepo) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Atomic", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Atomic indicates an expected call of Atomic.
func (mr *MockOrderRepoMockRecorder) Atomic(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Atomic", reflect.TypeOf((*MockOrderRepo)(nil).Atomic), ctx, fn)
}

// CancelOrder mocks base method.
func (m *MockOrderRepo) CancelOrder(ctx context.Context, id int64, symbol, exchange string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelOrder", ctx, id, symbol, exchange)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelOrder indicates an expected call of CancelOrder.
func (mr *MockOrderRepoMockRecorder) CancelOrder(ctx, id, symbol, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelOrder", reflect.TypeOf((*MockOrderRepo)(nil).CancelOrder), ctx, id, symbol, exchange)
}

// CreateOrder mocks base method.
func (m *MockOrderRepo) CreateOrder(ctx context.Context, order *domain.Order) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderRepoMockRecorder) CreateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderRepo)(nil).CreateOrder), ctx, order)
}

// CreateOrderWithSettings mocks base method.
func (m *MockOrderRepo) CreateOrderWithSettings(ctx context.Context, order *domain.Order, settings *domain.Settings) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrderWithSettings", ctx, order, settings)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrderWithSettings indicates an expected call of CreateOrderWithSettings.
func (mr *MockOrderRepoMockRecorder) CreateOrderWithSettings(ctx, order, settings interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrderWithSettings", reflect.TypeOf((*MockOrderRepo)(nil).CreateOrderWithSettings), ctx, order, settings)
}

// ExecuteOrder mocks base method.
func (m *MockOrderRepo) ExecuteOrder(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecuteOrder indicates an expected call of ExecuteOrder.
func (mr *MockOrderRepoMockRecorder) ExecuteOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteOrder", reflect.TypeOf((*MockOrderRepo)(nil).ExecuteOrder), ctx, id)
}

// GetActiveSymbols mocks base method.
func (m *MockOrderRepo) GetActiveSymbols(ctx context.Context) (domain.SymbolList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveSymbols", ctx)
	ret0, _ := ret[0].(domain.SymbolList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveSymbols indicates an expected call of GetActiveSymbols.
func (mr *MockOrderRepoMockRecorder) GetActiveSymbols(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveSymbols", reflect.TypeOf((*MockOrderRepo)(nil).GetActiveSymbols), ctx)
}

// GetActiveTpSlOrders mocks base method.
func (m *MockOrderRepo) GetActiveTpSlOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveTpSlOrders", ctx, exchange)
	ret0, _ := ret[0].([]*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveTpSlOrders indicates an expected call of GetActiveTpSlOrders.
func (mr *MockOrderRepoMockRecorder) GetActiveTpSlOrders(ctx, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveTpSlOrders", reflect.TypeOf((*MockOrderRepo)(nil).GetActiveTpSlOrders), ctx, exchange)
}

// GetLimitOrders mocks base method.
func (m *MockOrderRepo) GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimitOrders", ctx, exchange)
	ret0, _ := ret[0].([]*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLimitOrders indicates an expected call of GetLimitOrders.
func (mr *MockOrderRepoMockRecorder) GetLimitOrders(ctx, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimitOrders", reflect.TypeOf((*MockOrderRepo)(nil).GetLimitOrders), ctx, exchange)
}

// GetOpposingTpSlOrder mocks base method.
func (m *MockOrderRepo) GetOpposingTpSlOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpposingTpSlOrder", ctx, order)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpposingTpSlOrder indicates an expected call of GetOpposingTpSlOrder.
func (mr *MockOrderRepoMockRecorder) GetOpposingTpSlOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpposingTpSlOrder", reflect.TypeOf((*MockOrderRepo)(nil).GetOpposingTpSlOrder), ctx, order)
}

// GetOrder mocks base method.
func (m *MockOrderRepo) GetOrder(ctx context.Context, orderId int64, symbol, exchange string) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, orderId, symbol, exchange)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockOrderRepoMockRecorder) GetOrder(ctx, orderId, symbol, exchange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockOrderRepo)(nil).GetOrder), ctx, orderId, symbol, exchange)
}

// GetTpSlOrderByBaseOrder mocks base method.
func (m *MockOrderRepo) GetTpSlOrderByBaseOrder(ctx context.Context, id int64, symbol, exchange, tpsl string) (*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTpSlOrderByBaseOrder", ctx, id, symbol, exchange, tpsl)
	ret0, _ := ret[0].(*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTpSlOrderByBaseOrder indicates an expected call of GetTpSlOrderByBaseOrder.
func (mr *MockOrderRepoMockRecorder) GetTpSlOrderByBaseOrder(ctx, id, symbol, exchange, tpsl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTpSlOrderByBaseOrder", reflect.TypeOf((*MockOrderRepo)(nil).GetTpSlOrderByBaseOrder), ctx, id, symbol, exchange, tpsl)
}

// GetTpSlOrdersByBaseOrder mocks base method.
func (m *MockOrderRepo) GetTpSlOrdersByBaseOrder(ctx context.Context, id int64) ([]*domain.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTpSlOrdersByBaseOrder", ctx, id)
	ret0, _ := ret[0].([]*domain.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTpSlOrdersByBaseOrder indicates an expected call of GetTpSlOrdersByBaseOrder.
func (mr *MockOrderRepoMockRecorder) GetTpSlOrdersByBaseOrder(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTpSlOrdersByBaseOrder", reflect.TypeOf((*MockOrderRepo)(nil).GetTpSlOrdersByBaseOrder), ctx, id)
}

// UpdateTpSl mocks base method.
func (m *MockOrderRepo) UpdateTpSl(ctx context.Context, id int64, price string, settings *domain.Settings) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTpSl", ctx, id, price, settings)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTpSl indicates an expected call of UpdateTpSl.
func (mr *MockOrderRepoMockRecorder) UpdateTpSl(ctx, id, price, settings interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTpSl", reflect.TypeOf((*MockOrderRepo)(nil).UpdateTpSl), ctx, id, price, settings)
}