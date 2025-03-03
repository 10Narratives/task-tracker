// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TaskCompleter is an autogenerated mock type for the TaskCompleter type
type TaskCompleter struct {
	mock.Mock
}

// Complete provides a mock function with given fields: ctx, id
func (_m *TaskCompleter) Complete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Complete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTaskCompleter creates a new instance of TaskCompleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskCompleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskCompleter {
	mock := &TaskCompleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
