// Code generated by mockery v2.38.0. DO NOT EDIT.

package mockusercenterproto

import mock "github.com/stretchr/testify/mock"

// MockUnsafeUserCenterServer is an autogenerated mock type for the UnsafeUserCenterServer type
type MockUnsafeUserCenterServer struct {
	mock.Mock
}

type MockUnsafeUserCenterServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUnsafeUserCenterServer) EXPECT() *MockUnsafeUserCenterServer_Expecter {
	return &MockUnsafeUserCenterServer_Expecter{mock: &_m.Mock}
}

// mustEmbedUnimplementedUserCenterServer provides a mock function with given fields:
func (_m *MockUnsafeUserCenterServer) mustEmbedUnimplementedUserCenterServer() {
	_m.Called()
}

// MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedUserCenterServer'
type MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedUserCenterServer is a helper method to define mock.On call
func (_e *MockUnsafeUserCenterServer_Expecter) mustEmbedUnimplementedUserCenterServer() *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call {
	return &MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call{Call: _e.mock.On("mustEmbedUnimplementedUserCenterServer")}
}

func (_c *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call) Run(run func()) *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call) Return() *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call) RunAndReturn(run func()) *MockUnsafeUserCenterServer_mustEmbedUnimplementedUserCenterServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUnsafeUserCenterServer creates a new instance of MockUnsafeUserCenterServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUnsafeUserCenterServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUnsafeUserCenterServer {
	mock := &MockUnsafeUserCenterServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
