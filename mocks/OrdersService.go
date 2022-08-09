// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	orders "github.com/hromov/jevelina/useCases/orders"
	mock "github.com/stretchr/testify/mock"

	users "github.com/hromov/jevelina/domain/users"
)

// OrdersService is an autogenerated mock type for the Service type
type OrdersService struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *OrdersService) Create(_a0 context.Context, _a1 orders.Order) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, orders.Order) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateForUser provides a mock function with given fields: _a0, _a1, _a2
func (_m *OrdersService) CreateForUser(_a0 context.Context, _a1 orders.Order, _a2 users.User) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, orders.Order, users.User) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOrdersService interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrdersService creates a new instance of OrdersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrdersService(t mockConstructorTestingTNewOrdersService) *OrdersService {
	mock := &OrdersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}