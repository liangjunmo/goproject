package userserviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice"
)

func TestDefaultService(t *testing.T) {
	var (
		mutex         *mockedMutex
		mutexProvider *mockedMutexProvider
		repository    *mockedRepository
		service       *defaultService
		ctx           context.Context
	)

	beforeTest := func(t *testing.T) {
		mutex = &mockedMutex{}
		mutexProvider = &mockedMutexProvider{}
		repository = &mockedRepository{}
		service = newDefaultService(repository, mutexProvider)
		ctx = context.Background()
	}

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		repository.On("Search", ctx, mock.IsType(criteria{})).Return(map[uint32]userservice.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, userservice.SearchCommand{})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		repository.On("Get", ctx, uint32(1)).Return(userservice.User{UID: 1}, true, nil)

		user, err := service.Get(ctx, userservice.GetCommand{UID: 1})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		repository.On("GetByUsername", ctx, "user").Return(userservice.User{Username: "user"}, true, nil)

		user, err := service.GetByUsername(ctx, userservice.GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, "user", user.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		mutex.On("Lock", ctx).Return(nil)
		mutex.On("Unlock", ctx).Return(true, nil)

		mutexProvider.On("ProvideCreateUserMutex", "user").Return(mutex)

		repository.On("GetByUsername", ctx, "user").Return(userservice.User{}, false, nil)
		repository.On("Create", ctx, mock.IsType(&userservice.User{})).Return(nil)

		_, err := service.Create(ctx, userservice.CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		user := userservice.User{Password: cryptPassword("pass")}

		repository.On("GetByUsername", ctx, "user").Return(user, true, nil)

		err := service.ValidatePassword(ctx, userservice.ValidatePasswordCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})
}

type mockedMutexProvider struct {
	mock.Mock
}

func (mocked *mockedMutexProvider) ProvideCreateUserMutex(username string) mutex {
	args := mocked.Called(username)

	return args.Get(0).(mutex)
}

type mockedMutex struct {
	mock.Mock
}

func (mocked *mockedMutex) Lock(ctx context.Context) error {
	args := mocked.Called(ctx)

	return args.Error(0)
}

func (mocked *mockedMutex) Unlock(ctx context.Context) (ok bool, err error) {
	args := mocked.Called(ctx)

	return args.Bool(0), args.Error(1)
}

type mockedRepository struct {
	mock.Mock
}

func (mocked *mockedRepository) Begin() (repository, error) {
	return mocked, nil
}

func (mocked *mockedRepository) Commit() error {
	return nil
}

func (mocked *mockedRepository) Rollback() error {
	return nil
}

func (mocked *mockedRepository) Search(ctx context.Context, criteria criteria) (map[uint32]userservice.User, error) {
	args := mocked.Called(ctx, criteria)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[uint32]userservice.User), args.Error(1)
}

func (mocked *mockedRepository) Get(ctx context.Context, uid uint32) (user userservice.User, exist bool, err error) {
	args := mocked.Called(ctx, uid)

	return args.Get(0).(userservice.User), args.Bool(1), args.Error(2)
}

func (mocked *mockedRepository) GetByUsername(ctx context.Context, username string) (user userservice.User, exist bool, err error) {
	args := mocked.Called(ctx, username)

	return args.Get(0).(userservice.User), args.Bool(1), args.Error(2)
}

func (mocked *mockedRepository) Create(ctx context.Context, user *userservice.User) error {
	args := mocked.Called(ctx, user)

	return args.Error(0)
}
