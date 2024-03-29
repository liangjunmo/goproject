// Code generated by mockery v2.38.0. DO NOT EDIT.

package mockmutex

import (
	mutex "github.com/liangjunmo/goproject/internal/usercenter/mutex"
	mock "github.com/stretchr/testify/mock"
)

// MockMutexProvider is an autogenerated mock type for the MutexProvider type
type MockMutexProvider struct {
	mock.Mock
}

type MockMutexProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMutexProvider) EXPECT() *MockMutexProvider_Expecter {
	return &MockMutexProvider_Expecter{mock: &_m.Mock}
}

// ProvideCreateUserMutex provides a mock function with given fields: username
func (_m *MockMutexProvider) ProvideCreateUserMutex(username string) mutex.Mutex {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for ProvideCreateUserMutex")
	}

	var r0 mutex.Mutex
	if rf, ok := ret.Get(0).(func(string) mutex.Mutex); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mutex.Mutex)
		}
	}

	return r0
}

// MockMutexProvider_ProvideCreateUserMutex_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProvideCreateUserMutex'
type MockMutexProvider_ProvideCreateUserMutex_Call struct {
	*mock.Call
}

// ProvideCreateUserMutex is a helper method to define mock.On call
//   - username string
func (_e *MockMutexProvider_Expecter) ProvideCreateUserMutex(username interface{}) *MockMutexProvider_ProvideCreateUserMutex_Call {
	return &MockMutexProvider_ProvideCreateUserMutex_Call{Call: _e.mock.On("ProvideCreateUserMutex", username)}
}

func (_c *MockMutexProvider_ProvideCreateUserMutex_Call) Run(run func(username string)) *MockMutexProvider_ProvideCreateUserMutex_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockMutexProvider_ProvideCreateUserMutex_Call) Return(_a0 mutex.Mutex) *MockMutexProvider_ProvideCreateUserMutex_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMutexProvider_ProvideCreateUserMutex_Call) RunAndReturn(run func(string) mutex.Mutex) *MockMutexProvider_ProvideCreateUserMutex_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMutexProvider creates a new instance of MockMutexProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMutexProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMutexProvider {
	mock := &MockMutexProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
