// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	context "context"

	math "github.com/berachain/beacon-kit/primitives/math"
	mock "github.com/stretchr/testify/mock"
)

// AvailabilityStore is an autogenerated mock type for the AvailabilityStore type
type AvailabilityStore[BlobSidecarsT any] struct {
	mock.Mock
}

type AvailabilityStore_Expecter[BlobSidecarsT any] struct {
	mock *mock.Mock
}

func (_m *AvailabilityStore[BlobSidecarsT]) EXPECT() *AvailabilityStore_Expecter[BlobSidecarsT] {
	return &AvailabilityStore_Expecter[BlobSidecarsT]{mock: &_m.Mock}
}

// IsDataAvailable provides a mock function with given fields: _a0, _a1
func (_m *AvailabilityStore[BlobSidecarsT]) IsDataAvailable(_a0 context.Context, _a1 math.U64) bool {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for IsDataAvailable")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, math.U64) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// AvailabilityStore_IsDataAvailable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsDataAvailable'
type AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT any] struct {
	*mock.Call
}

// IsDataAvailable is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 math.U64
func (_e *AvailabilityStore_Expecter[BlobSidecarsT]) IsDataAvailable(_a0 interface{}, _a1 interface{}) *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT] {
	return &AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT]{Call: _e.mock.On("IsDataAvailable", _a0, _a1)}
}

func (_c *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT]) Run(run func(_a0 context.Context, _a1 math.U64)) *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(math.U64))
	})
	return _c
}

func (_c *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT]) Return(_a0 bool) *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT]) RunAndReturn(run func(context.Context, math.U64) bool) *AvailabilityStore_IsDataAvailable_Call[BlobSidecarsT] {
	_c.Call.Return(run)
	return _c
}

// Persist provides a mock function with given fields: _a0, _a1
func (_m *AvailabilityStore[BlobSidecarsT]) Persist(_a0 math.U64, _a1 BlobSidecarsT) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Persist")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(math.U64, BlobSidecarsT) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AvailabilityStore_Persist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Persist'
type AvailabilityStore_Persist_Call[BlobSidecarsT any] struct {
	*mock.Call
}

// Persist is a helper method to define mock.On call
//   - _a0 math.U64
//   - _a1 BlobSidecarsT
func (_e *AvailabilityStore_Expecter[BlobSidecarsT]) Persist(_a0 interface{}, _a1 interface{}) *AvailabilityStore_Persist_Call[BlobSidecarsT] {
	return &AvailabilityStore_Persist_Call[BlobSidecarsT]{Call: _e.mock.On("Persist", _a0, _a1)}
}

func (_c *AvailabilityStore_Persist_Call[BlobSidecarsT]) Run(run func(_a0 math.U64, _a1 BlobSidecarsT)) *AvailabilityStore_Persist_Call[BlobSidecarsT] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(math.U64), args[1].(BlobSidecarsT))
	})
	return _c
}

func (_c *AvailabilityStore_Persist_Call[BlobSidecarsT]) Return(_a0 error) *AvailabilityStore_Persist_Call[BlobSidecarsT] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AvailabilityStore_Persist_Call[BlobSidecarsT]) RunAndReturn(run func(math.U64, BlobSidecarsT) error) *AvailabilityStore_Persist_Call[BlobSidecarsT] {
	_c.Call.Return(run)
	return _c
}

// NewAvailabilityStore creates a new instance of AvailabilityStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAvailabilityStore[BlobSidecarsT any](t interface {
	mock.TestingT
	Cleanup(func())
}) *AvailabilityStore[BlobSidecarsT] {
	mock := &AvailabilityStore[BlobSidecarsT]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
