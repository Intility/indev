// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Serializer is an autogenerated mock type for the Serializer type
type Serializer struct {
	mock.Mock
}

type Serializer_Expecter struct {
	mock *mock.Mock
}

func (_m *Serializer) EXPECT() *Serializer_Expecter {
	return &Serializer_Expecter{mock: &_m.Mock}
}

// Marshal provides a mock function with given fields:
func (_m *Serializer) Marshal() ([]byte, error) {
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

// Serializer_Marshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Marshal'
type Serializer_Marshal_Call struct {
	*mock.Call
}

// Marshal is a helper method to define mock.On call
func (_e *Serializer_Expecter) Marshal() *Serializer_Marshal_Call {
	return &Serializer_Marshal_Call{Call: _e.mock.On("Marshal")}
}

func (_c *Serializer_Marshal_Call) Run(run func()) *Serializer_Marshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Serializer_Marshal_Call) Return(_a0 []byte, _a1 error) *Serializer_Marshal_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Serializer_Marshal_Call) RunAndReturn(run func() ([]byte, error)) *Serializer_Marshal_Call {
	_c.Call.Return(run)
	return _c
}

// Unmarshal provides a mock function with given fields: _a0
func (_m *Serializer) Unmarshal(_a0 []byte) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Unmarshal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Serializer_Unmarshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unmarshal'
type Serializer_Unmarshal_Call struct {
	*mock.Call
}

// Unmarshal is a helper method to define mock.On call
//   - _a0 []byte
func (_e *Serializer_Expecter) Unmarshal(_a0 interface{}) *Serializer_Unmarshal_Call {
	return &Serializer_Unmarshal_Call{Call: _e.mock.On("Unmarshal", _a0)}
}

func (_c *Serializer_Unmarshal_Call) Run(run func(_a0 []byte)) *Serializer_Unmarshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *Serializer_Unmarshal_Call) Return(_a0 error) *Serializer_Unmarshal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Serializer_Unmarshal_Call) RunAndReturn(run func([]byte) error) *Serializer_Unmarshal_Call {
	_c.Call.Return(run)
	return _c
}

// NewSerializer creates a new instance of Serializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSerializer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Serializer {
	mock := &Serializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}