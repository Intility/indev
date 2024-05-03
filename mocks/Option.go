// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	tokencache "github.com/intility/icpctl/pkg/tokencache"
	mock "github.com/stretchr/testify/mock"
)

// Option is an autogenerated mock type for the Option type
type Option struct {
	mock.Mock
}

type Option_Expecter struct {
	mock *mock.Mock
}

func (_m *Option) EXPECT() *Option_Expecter {
	return &Option_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *Option) Execute(_a0 *tokencache.TokenCache) {
	_m.Called(_a0)
}

// Option_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Option_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *tokencache.TokenCache
func (_e *Option_Expecter) Execute(_a0 interface{}) *Option_Execute_Call {
	return &Option_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *Option_Execute_Call) Run(run func(_a0 *tokencache.TokenCache)) *Option_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*tokencache.TokenCache))
	})
	return _c
}

func (_c *Option_Execute_Call) Return() *Option_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *Option_Execute_Call) RunAndReturn(run func(*tokencache.TokenCache)) *Option_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewOption creates a new instance of Option. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *Option {
	mock := &Option{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
