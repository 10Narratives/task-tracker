// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TaskRemover is an autogenerated mock type for the TaskRemover type
type TaskRemover struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TaskRemover) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTaskRemover creates a new instance of TaskRemover. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskRemover(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskRemover {
	mock := &TaskRemover{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
