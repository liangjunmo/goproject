package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/goproject/model"
	"github.com/liangjunmo/goproject/internal/goproject/repository"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
	mockmutex "github.com/liangjunmo/goproject/mocks/goproject/mutex"
	mockrepository "github.com/liangjunmo/goproject/mocks/goproject/repository"
	mockusercenterproto "github.com/liangjunmo/goproject/mocks/usercenterproto"
)

func TestUserService(t *testing.T) {
	var (
		mockMutex            *mockmutex.MockMutex
		mockMutexProvider    *mockmutex.MockMutexProvider
		mockUserRepository   *mockrepository.MockUserRepository
		mockUserCenterClient *mockusercenterproto.MockUserCenterClient
		service              *userService
		ctx                  context.Context
	)

	beforeTest := func(t *testing.T) {
		mockMutex = &mockmutex.MockMutex{}
		mockMutexProvider = &mockmutex.MockMutexProvider{}
		mockUserRepository = &mockrepository.MockUserRepository{}
		mockUserCenterClient = &mockusercenterproto.MockUserCenterClient{}

		service = newUserService(mockMutexProvider, mockUserRepository, mockUserCenterClient)

		ctx = context.Background()
	}

	t.Run("List", func(t *testing.T) {
		beforeTest(t)

		mockUserRepository.
			On("List", ctx, mock.IsType(repository.UserCriteria{})).
			Return(pagination.DefaultPagination{}, []model.User{{UID: 1}}, nil)

		mockUserCenterClient.
			On("SearchUser", ctx, mock.IsType(&usercenterproto.SearchUserRequest{})).
			Return(
				&usercenterproto.SearchUserReply{
					Users: map[uint32]*usercenterproto.User{1: {UID: 1}},
				},
				nil,
			)

		_, users, err := service.List(ctx, ListCommand{})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[0].UID)
		require.Equal(t, uint32(1), users[0].UserCenterUser.UID)
	})

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		mockUserCenterClient.
			On("SearchUser", ctx, mock.IsType(&usercenterproto.SearchUserRequest{})).
			Return(
				&usercenterproto.SearchUserReply{
					Users: map[uint32]*usercenterproto.User{1: {UID: 1}},
				},
				nil,
			)

		mockUserRepository.
			On("Search", ctx, mock.IsType(repository.UserCriteria{})).
			Return(map[uint32]model.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, SearchCommand{Uids: []uint32{1}})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
		require.Equal(t, uint32(1), users[1].UserCenterUser.UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		mockUserRepository.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		mockUserCenterClient.
			On("GetUserByUID", ctx, mock.IsType(&usercenterproto.GetUserByUIDRequest{})).
			Return(
				&usercenterproto.GetUserByUIDReply{
					User: &usercenterproto.User{UID: 1},
				},
				nil,
			)

		user, err := service.Get(ctx, GetCommand{UID: 1})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, uint32(1), user.UserCenterUser.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		mockUserCenterClient.
			On("GetUserByUsername", ctx, mock.IsType(&usercenterproto.GetUserByUsernameRequest{})).
			Return(
				&usercenterproto.GetUserByUsernameReply{
					User: &usercenterproto.User{
						UID:      1,
						Username: "user",
					},
				},
				nil,
			)

		mockUserRepository.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		user, err := service.GetByUsername(ctx, GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, "user", user.UserCenterUser.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		mockUserCenterClient.
			On("CreateUser", ctx, mock.IsType(&usercenterproto.CreateUserRequest{})).
			Return(
				&usercenterproto.CreateUserReply{
					UID: 1,
				},
				nil,
			)

		mockMutex.On("Lock", ctx).Return(nil)
		mockMutex.On("Unlock", ctx).Return(true, nil)

		mockMutexProvider.On("ProvideCreateUserMutex", uint32(1)).Return(mockMutex)

		mockUserRepository.On("Get", ctx, uint32(1)).Return(model.User{}, false, nil)
		mockUserRepository.On("Create", ctx, mock.IsType(&model.User{})).Return(nil)

		uid, err := service.Create(ctx, CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), uid)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		mockUserCenterClient.
			On("ValidatePassword", ctx, mock.IsType(&usercenterproto.ValidatePasswordRequest{})).
			Return(&usercenterproto.ValidatePasswordReply{}, nil)

		err := service.ValidatePassword(ctx, ValidatePasswordCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
	})
}
