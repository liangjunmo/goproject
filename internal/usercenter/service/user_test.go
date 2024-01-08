package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/usercenter/model"
	"github.com/liangjunmo/goproject/internal/usercenter/mutex"
	"github.com/liangjunmo/goproject/internal/usercenter/repository"
)

func TestUserService(t *testing.T) {
	var (
		mutexMocked         *mockedMutex
		mutexProviderMocked *mockedMutexProvider
		repositoryMocked    *mockedRepository
		service             *userService
		ctx                 context.Context
	)

	beforeTest := func(t *testing.T) {
		mutexMocked = &mockedMutex{}
		mutexProviderMocked = &mockedMutexProvider{}
		repositoryMocked = &mockedRepository{}

		service = newUserService(repositoryMocked, mutexProviderMocked)

		ctx = context.Background()
	}

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		repositoryMocked.
			On("Search", ctx, mock.IsType(repository.UserCriteria{})).
			Return(map[uint32]model.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, SearchCommand{})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		repositoryMocked.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		user, err := service.Get(ctx, GetCommand{UID: 1})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		repositoryMocked.On("GetByUsername", ctx, "user").Return(model.User{Username: "user"}, true, nil)

		user, err := service.GetByUsername(ctx, GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, "user", user.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		mutexMocked.On("Lock", ctx).Return(nil)
		mutexMocked.On("Unlock", ctx).Return(true, nil)

		mutexProviderMocked.On("ProvideCreateUserMutex", "user").Return(mutexMocked)

		repositoryMocked.On("GetByUsername", ctx, "user").Return(model.User{}, false, nil)
		repositoryMocked.On("Create", ctx, mock.IsType(&model.User{})).Return(nil)

		_, err := service.Create(ctx, CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		user := model.User{Password: cryptPassword("pass")}

		repositoryMocked.On("GetByUsername", ctx, "user").Return(user, true, nil)

		err := service.ValidatePassword(ctx, ValidatePasswordCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})
}

type mockedMutexProvider struct {
	mock.Mock
}

func (mocked *mockedMutexProvider) ProvideCreateUserMutex(username string) mutex.Mutex {
	args := mocked.Called(username)

	return args.Get(0).(mutex.Mutex)
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

func (mocked *mockedRepository) Begin() (repository.UserRepository, error) {
	return mocked, nil
}

func (mocked *mockedRepository) Commit() error {
	return nil
}

func (mocked *mockedRepository) Rollback() error {
	return nil
}

func (mocked *mockedRepository) Search(ctx context.Context, criteria repository.UserCriteria) (map[uint32]model.User, error) {
	args := mocked.Called(ctx, criteria)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[uint32]model.User), args.Error(1)
}

func (mocked *mockedRepository) Get(ctx context.Context, uid uint32) (user model.User, exist bool, err error) {
	args := mocked.Called(ctx, uid)

	return args.Get(0).(model.User), args.Bool(1), args.Error(2)
}

func (mocked *mockedRepository) GetByUsername(ctx context.Context, username string) (user model.User, exist bool, err error) {
	args := mocked.Called(ctx, username)

	return args.Get(0).(model.User), args.Bool(1), args.Error(2)
}

func (mocked *mockedRepository) Create(ctx context.Context, user *model.User) error {
	args := mocked.Called(ctx, user)

	return args.Error(0)
}
