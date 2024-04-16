// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Marshaler is an autogenerated mock type for the Marshaler type
type Marshaler struct {
	mock.Mock
}

type Marshaler_Expecter struct {
	mock *mock.Mock
}

func (_m *Marshaler) EXPECT() *Marshaler_Expecter {
	return &Marshaler_Expecter{mock: &_m.Mock}
}

// Marshal provides a mock function with given fields:
func (_m *Marshaler) Marshal() ([]byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Marshal")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Marshaler_Marshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Marshal'
type Marshaler_Marshal_Call struct {
	*mock.Call
}

// Marshal is a helper method to define mock.On call
func (_e *Marshaler_Expecter) Marshal() *Marshaler_Marshal_Call {
	return &Marshaler_Marshal_Call{Call: _e.mock.On("Marshal")}
}

func (_c *Marshaler_Marshal_Call) Run(run func()) *Marshaler_Marshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Marshaler_Marshal_Call) Return(_a0 []byte, _a1 error) *Marshaler_Marshal_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Marshaler_Marshal_Call) RunAndReturn(run func() ([]byte, error)) *Marshaler_Marshal_Call {
	_c.Call.Return(run)
	return _c
}

// NewMarshaler creates a new instance of Marshaler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMarshaler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Marshaler {
	mock := &Marshaler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
