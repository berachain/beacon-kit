// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Basic is an autogenerated mock type for the Basic type
type Basic struct {
	mock.Mock
}

type Basic_Expecter struct {
	mock *mock.Mock
}

func (_m *Basic) EXPECT() *Basic_Expecter {
	return &Basic_Expecter{mock: &_m.Mock}
}

// Name provides a mock function with given fields:
func (_m *Basic) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Basic_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type Basic_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *Basic_Expecter) Name() *Basic_Name_Call {
	return &Basic_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *Basic_Name_Call) Run(run func()) *Basic_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Basic_Name_Call) Return(_a0 string) *Basic_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Basic_Name_Call) RunAndReturn(run func() string) *Basic_Name_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: ctx
func (_m *Basic) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Basic_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type Basic_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Basic_Expecter) Start(ctx interface{}) *Basic_Start_Call {
	return &Basic_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *Basic_Start_Call) Run(run func(ctx context.Context)) *Basic_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Basic_Start_Call) Return(_a0 error) *Basic_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Basic_Start_Call) RunAndReturn(run func(context.Context) error) *Basic_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewBasic creates a new instance of Basic. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBasic(t interface {
	mock.TestingT
	Cleanup(func())
}) *Basic {
	mock := &Basic{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}