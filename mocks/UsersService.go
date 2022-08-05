// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	users "github.com/hromov/jevelina/domain/users"
	mock "github.com/stretchr/testify/mock"
)

// UsersService is an autogenerated mock type for the Service type
type UsersService struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *UsersService) Create(_a0 context.Context, _a1 users.ChangeUser) (users.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 users.User
	if rf, ok := ret.Get(0).(func(context.Context, users.ChangeUser) users.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, users.ChangeUser) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRole provides a mock function with given fields: _a0, _a1
func (_m *UsersService) CreateRole(_a0 context.Context, _a1 users.Role) (users.Role, error) {
	ret := _m.Called(_a0, _a1)

	var r0 users.Role
	if rf, ok := ret.Get(0).(func(context.Context, users.Role) users.Role); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(users.Role)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, users.Role) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *UsersService) Delete(_a0 context.Context, _a1 uint64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *UsersService) Get(_a0 context.Context, _a1 uint64) (users.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 users.User
	if rf, ok := ret.Get(0).(func(context.Context, uint64) users.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(users.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: _a0
func (_m *UsersService) List(_a0 context.Context) ([]users.User, error) {
	ret := _m.Called(_a0)

	var r0 []users.User
	if rf, ok := ret.Get(0).(func(context.Context) []users.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *UsersService) Update(_a0 context.Context, _a1 users.ChangeUser) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, users.ChangeUser) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUsersService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUsersService creates a new instance of UsersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUsersService(t mockConstructorTestingTNewUsersService) *UsersService {
	mock := &UsersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
