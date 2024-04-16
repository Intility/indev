// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	credentialstore "github.com/intility/minctl/pkg/credentialstore"
	mock "github.com/stretchr/testify/mock"
)

// FileSystemOption is an autogenerated mock type for the FileSystemOption type
type FileSystemOption struct {
	mock.Mock
}

type FileSystemOption_Expecter struct {
	mock *mock.Mock
}

func (_m *FileSystemOption) EXPECT() *FileSystemOption_Expecter {
	return &FileSystemOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *FileSystemOption) Execute(_a0 *credentialstore.FilesystemCredentialStore) {
	_m.Called(_a0)
}

// FileSystemOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type FileSystemOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *credentialstore.FilesystemCredentialStore
func (_e *FileSystemOption_Expecter) Execute(_a0 interface{}) *FileSystemOption_Execute_Call {
	return &FileSystemOption_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *FileSystemOption_Execute_Call) Run(run func(_a0 *credentialstore.FilesystemCredentialStore)) *FileSystemOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*credentialstore.FilesystemCredentialStore))
	})
	return _c
}

func (_c *FileSystemOption_Execute_Call) Return() *FileSystemOption_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *FileSystemOption_Execute_Call) RunAndReturn(run func(*credentialstore.FilesystemCredentialStore)) *FileSystemOption_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewFileSystemOption creates a new instance of FileSystemOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileSystemOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileSystemOption {
	mock := &FileSystemOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
