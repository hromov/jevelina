// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	files "github.com/hromov/jevelina/domain/misc/files"
	mock "github.com/stretchr/testify/mock"
)

// FilesStorageService is an autogenerated mock type for the Storage type
type FilesStorageService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, url
func (_m *FilesStorageService) Delete(ctx context.Context, url string) error {
	ret := _m.Called(ctx, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PresignUrl provides a mock function with given fields: url
func (_m *FilesStorageService) PresignUrl(url string) (string, error) {
	ret := _m.Called(url)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upload provides a mock function with given fields: ctx, req
func (_m *FilesStorageService) Upload(ctx context.Context, req files.FileAddReq) (files.FileCreateReq, error) {
	ret := _m.Called(ctx, req)

	var r0 files.FileCreateReq
	if rf, ok := ret.Get(0).(func(context.Context, files.FileAddReq) files.FileCreateReq); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(files.FileCreateReq)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, files.FileAddReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFilesStorageService interface {
	mock.TestingT
	Cleanup(func())
}

// NewFilesStorageService creates a new instance of FilesStorageService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilesStorageService(t mockConstructorTestingTNewFilesStorageService) *FilesStorageService {
	mock := &FilesStorageService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}