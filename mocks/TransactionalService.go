// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TransactionalService is an autogenerated mock type for the TransactionalService type
type TransactionalService struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx, fn
func (_m *TransactionalService) Exec(ctx context.Context, fn func(context.Context) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactionalService creates a new instance of TransactionalService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionalService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionalService {
	mock := &TransactionalService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}