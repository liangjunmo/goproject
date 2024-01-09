// Code generated by mockery v2.38.0. DO NOT EDIT.

package mockrepository

import (
	context "context"

	model "github.com/liangjunmo/goproject/internal/usercenter/model"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/liangjunmo/goproject/internal/usercenter/repository"
)

// MockUserRepository is an autogenerated mock type for the UserRepository type
type MockUserRepository struct {
	mock.Mock
}

type MockUserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserRepository) EXPECT() *MockUserRepository_Expecter {
	return &MockUserRepository_Expecter{mock: &_m.Mock}
}

// Begin provides a mock function with given fields:
func (_m *MockUserRepository) Begin() (repository.UserRepository, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 repository.UserRepository
	var r1 error
	if rf, ok := ret.Get(0).(func() (repository.UserRepository, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() repository.UserRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.UserRepository)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserRepository_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type MockUserRepository_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
func (_e *MockUserRepository_Expecter) Begin() *MockUserRepository_Begin_Call {
	return &MockUserRepository_Begin_Call{Call: _e.mock.On("Begin")}
}

func (_c *MockUserRepository_Begin_Call) Run(run func()) *MockUserRepository_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserRepository_Begin_Call) Return(_a0 repository.UserRepository, _a1 error) *MockUserRepository_Begin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserRepository_Begin_Call) RunAndReturn(run func() (repository.UserRepository, error)) *MockUserRepository_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// Commit provides a mock function with given fields:
func (_m *MockUserRepository) Commit() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserRepository_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type MockUserRepository_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
func (_e *MockUserRepository_Expecter) Commit() *MockUserRepository_Commit_Call {
	return &MockUserRepository_Commit_Call{Call: _e.mock.On("Commit")}
}

func (_c *MockUserRepository_Commit_Call) Run(run func()) *MockUserRepository_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserRepository_Commit_Call) Return(_a0 error) *MockUserRepository_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserRepository_Commit_Call) RunAndReturn(run func() error) *MockUserRepository_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, user
func (_m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockUserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - user *model.User
func (_e *MockUserRepository_Expecter) Create(ctx interface{}, user interface{}) *MockUserRepository_Create_Call {
	return &MockUserRepository_Create_Call{Call: _e.mock.On("Create", ctx, user)}
}

func (_c *MockUserRepository_Create_Call) Run(run func(ctx context.Context, user *model.User)) *MockUserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.User))
	})
	return _c
}

func (_c *MockUserRepository_Create_Call) Return(_a0 error) *MockUserRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserRepository_Create_Call) RunAndReturn(run func(context.Context, *model.User) error) *MockUserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, uid
func (_m *MockUserRepository) Get(ctx context.Context, uid uint32) (model.User, bool, error) {
	ret := _m.Called(ctx, uid)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 model.User
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (model.User, bool, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) model.User); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) bool); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uint32) error); ok {
		r2 = rf(ctx, uid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockUserRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockUserRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - uid uint32
func (_e *MockUserRepository_Expecter) Get(ctx interface{}, uid interface{}) *MockUserRepository_Get_Call {
	return &MockUserRepository_Get_Call{Call: _e.mock.On("Get", ctx, uid)}
}

func (_c *MockUserRepository_Get_Call) Run(run func(ctx context.Context, uid uint32)) *MockUserRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint32))
	})
	return _c
}

func (_c *MockUserRepository_Get_Call) Return(user model.User, exist bool, err error) *MockUserRepository_Get_Call {
	_c.Call.Return(user, exist, err)
	return _c
}

func (_c *MockUserRepository_Get_Call) RunAndReturn(run func(context.Context, uint32) (model.User, bool, error)) *MockUserRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUsername provides a mock function with given fields: ctx, username
func (_m *MockUserRepository) GetByUsername(ctx context.Context, username string) (model.User, bool, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetByUsername")
	}

	var r0 model.User
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.User, bool, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, username)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockUserRepository_GetByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUsername'
type MockUserRepository_GetByUsername_Call struct {
	*mock.Call
}

// GetByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *MockUserRepository_Expecter) GetByUsername(ctx interface{}, username interface{}) *MockUserRepository_GetByUsername_Call {
	return &MockUserRepository_GetByUsername_Call{Call: _e.mock.On("GetByUsername", ctx, username)}
}

func (_c *MockUserRepository_GetByUsername_Call) Run(run func(ctx context.Context, username string)) *MockUserRepository_GetByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUserRepository_GetByUsername_Call) Return(user model.User, exist bool, err error) *MockUserRepository_GetByUsername_Call {
	_c.Call.Return(user, exist, err)
	return _c
}

func (_c *MockUserRepository_GetByUsername_Call) RunAndReturn(run func(context.Context, string) (model.User, bool, error)) *MockUserRepository_GetByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with given fields:
func (_m *MockUserRepository) Rollback() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserRepository_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type MockUserRepository_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
func (_e *MockUserRepository_Expecter) Rollback() *MockUserRepository_Rollback_Call {
	return &MockUserRepository_Rollback_Call{Call: _e.mock.On("Rollback")}
}

func (_c *MockUserRepository_Rollback_Call) Run(run func()) *MockUserRepository_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUserRepository_Rollback_Call) Return(_a0 error) *MockUserRepository_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserRepository_Rollback_Call) RunAndReturn(run func() error) *MockUserRepository_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// Search provides a mock function with given fields: ctx, criteria
func (_m *MockUserRepository) Search(ctx context.Context, criteria repository.UserCriteria) (map[uint32]model.User, error) {
	ret := _m.Called(ctx, criteria)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 map[uint32]model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.UserCriteria) (map[uint32]model.User, error)); ok {
		return rf(ctx, criteria)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.UserCriteria) map[uint32]model.User); ok {
		r0 = rf(ctx, criteria)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint32]model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.UserCriteria) error); ok {
		r1 = rf(ctx, criteria)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserRepository_Search_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Search'
type MockUserRepository_Search_Call struct {
	*mock.Call
}

// Search is a helper method to define mock.On call
//   - ctx context.Context
//   - criteria repository.UserCriteria
func (_e *MockUserRepository_Expecter) Search(ctx interface{}, criteria interface{}) *MockUserRepository_Search_Call {
	return &MockUserRepository_Search_Call{Call: _e.mock.On("Search", ctx, criteria)}
}

func (_c *MockUserRepository_Search_Call) Run(run func(ctx context.Context, criteria repository.UserCriteria)) *MockUserRepository_Search_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.UserCriteria))
	})
	return _c
}

func (_c *MockUserRepository_Search_Call) Return(_a0 map[uint32]model.User, _a1 error) *MockUserRepository_Search_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserRepository_Search_Call) RunAndReturn(run func(context.Context, repository.UserCriteria) (map[uint32]model.User, error)) *MockUserRepository_Search_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserRepository creates a new instance of MockUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserRepository {
	mock := &MockUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
