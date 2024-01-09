package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/usercenter/model"
	"github.com/liangjunmo/goproject/internal/usercenter/repository"
	mockmutex "github.com/liangjunmo/goproject/mocks/usercenter/mutex"
	mockrepository "github.com/liangjunmo/goproject/mocks/usercenter/repository"
)

func TestUserService(t *testing.T) {
	var (
		mockMutex          *mockmutex.MockMutex
		mockMutexProvider  *mockmutex.MockMutexProvider
		mockUserRepository *mockrepository.MockUserRepository
		service            *userService
		ctx                context.Context
	)

	beforeTest := func(t *testing.T) {
		mockMutex = &mockmutex.MockMutex{}
		mockMutexProvider = &mockmutex.MockMutexProvider{}
		mockUserRepository = &mockrepository.MockUserRepository{}

		service = newUserService(mockMutexProvider, mockUserRepository)

		ctx = context.Background()
	}

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		mockUserRepository.
			On("Search", ctx, mock.IsType(repository.UserCriteria{})).
			Return(map[uint32]model.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, SearchCommand{})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		mockUserRepository.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		user, err := service.Get(ctx, GetCommand{UID: 1})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		mockUserRepository.On("GetByUsername", ctx, "user").Return(model.User{Username: "user"}, true, nil)

		user, err := service.GetByUsername(ctx, GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, "user", user.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		mockMutex.On("Lock", ctx).Return(nil)
		mockMutex.On("Unlock", ctx).Return(true, nil)

		mockMutexProvider.On("ProvideCreateUserMutex", "user").Return(mockMutex)

		mockUserRepository.On("GetByUsername", ctx, "user").Return(model.User{}, false, nil)
		mockUserRepository.On("Create", ctx, mock.IsType(&model.User{})).Return(nil)

		_, err := service.Create(ctx, CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		user := model.User{Password: cryptPassword("pass")}

		mockUserRepository.On("GetByUsername", ctx, "user").Return(user, true, nil)

		err := service.ValidatePassword(ctx, ValidatePasswordCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})
}
