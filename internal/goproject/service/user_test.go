package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/goproject/model"
	"github.com/liangjunmo/goproject/internal/goproject/mutex"
	"github.com/liangjunmo/goproject/internal/goproject/repository"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

func TestUserService(t *testing.T) {
	var (
		mutexMocked            *mockedMutex
		mutexProviderMocked    *mockedMutexProvider
		repositoryMocked       *mockedRepository
		userCenterClientMocked *mockedUserCenterClient
		service                *userService
		ctx                    context.Context
	)

	beforeTest := func(t *testing.T) {
		mutexMocked = &mockedMutex{}
		mutexProviderMocked = &mockedMutexProvider{}
		repositoryMocked = &mockedRepository{}
		userCenterClientMocked = &mockedUserCenterClient{}

		service = newUserService(repositoryMocked, mutexProviderMocked, userCenterClientMocked)

		ctx = context.Background()
	}

	t.Run("List", func(t *testing.T) {
		beforeTest(t)

		repositoryMocked.
			On("List", ctx, mock.IsType(repository.UserCriteria{})).
			Return(pagination.DefaultPagination{}, []model.User{{UID: 1}}, nil)

		userCenterClientMocked.
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

		userCenterClientMocked.
			On("SearchUser", ctx, mock.IsType(&usercenterproto.SearchUserRequest{})).
			Return(
				&usercenterproto.SearchUserReply{
					Users: map[uint32]*usercenterproto.User{1: {UID: 1}},
				},
				nil,
			)

		repositoryMocked.
			On("Search", ctx, mock.IsType(repository.UserCriteria{})).
			Return(map[uint32]model.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, SearchCommand{Uids: []uint32{1}})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
		require.Equal(t, uint32(1), users[1].UserCenterUser.UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		repositoryMocked.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		userCenterClientMocked.
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

		userCenterClientMocked.
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

		repositoryMocked.On("Get", ctx, uint32(1)).Return(model.User{UID: 1}, true, nil)

		user, err := service.GetByUsername(ctx, GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, "user", user.UserCenterUser.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		userCenterClientMocked.
			On("CreateUser", ctx, mock.IsType(&usercenterproto.CreateUserRequest{})).
			Return(
				&usercenterproto.CreateUserReply{
					UID: 1,
				},
				nil,
			)

		mutexMocked.On("Lock", ctx).Return(nil)
		mutexMocked.On("Unlock", ctx).Return(true, nil)

		mutexProviderMocked.On("ProvideCreateUserMutex", uint32(1)).Return(mutexMocked)

		repositoryMocked.On("Get", ctx, uint32(1)).Return(model.User{}, false, nil)
		repositoryMocked.On("Create", ctx, mock.IsType(&model.User{})).Return(nil)

		uid, err := service.Create(ctx, CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), uid)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		userCenterClientMocked.
			On("ValidatePassword", ctx, mock.IsType(&usercenterproto.ValidatePasswordRequest{})).
			Return(&usercenterproto.ValidatePasswordReply{}, nil)

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

func (mocked *mockedMutexProvider) ProvideCreateUserMutex(uid uint32) mutex.Mutex {
	args := mocked.Called(uid)

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

func (mocked *mockedRepository) List(ctx context.Context, criteria repository.UserCriteria) (pagination.Pagination, []model.User, error) {
	args := mocked.Called(ctx, criteria)

	if args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(pagination.Pagination), args.Get(1).([]model.User), args.Error(2)
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

func (mocked *mockedRepository) Create(ctx context.Context, user *model.User) error {
	args := mocked.Called(ctx, user)

	return args.Error(0)
}

type mockedUserCenterClient struct {
	mock.Mock
}

func (mocked *mockedUserCenterClient) SearchUser(ctx context.Context, in *usercenterproto.SearchUserRequest, opts ...grpc.CallOption) (*usercenterproto.SearchUserReply, error) {
	args := mocked.Called(ctx, in)

	return args.Get(0).(*usercenterproto.SearchUserReply), args.Error(1)
}

func (mocked *mockedUserCenterClient) GetUserByUID(ctx context.Context, in *usercenterproto.GetUserByUIDRequest, opts ...grpc.CallOption) (*usercenterproto.GetUserByUIDReply, error) {
	args := mocked.Called(ctx, in)

	return args.Get(0).(*usercenterproto.GetUserByUIDReply), args.Error(1)
}

func (mocked *mockedUserCenterClient) GetUserByUsername(ctx context.Context, in *usercenterproto.GetUserByUsernameRequest, opts ...grpc.CallOption) (*usercenterproto.GetUserByUsernameReply, error) {
	args := mocked.Called(ctx, in)

	return args.Get(0).(*usercenterproto.GetUserByUsernameReply), args.Error(1)
}

func (mocked *mockedUserCenterClient) CreateUser(ctx context.Context, in *usercenterproto.CreateUserRequest, opts ...grpc.CallOption) (*usercenterproto.CreateUserReply, error) {
	args := mocked.Called(ctx, in)

	return args.Get(0).(*usercenterproto.CreateUserReply), args.Error(1)
}

func (mocked *mockedUserCenterClient) ValidatePassword(ctx context.Context, in *usercenterproto.ValidatePasswordRequest, opts ...grpc.CallOption) (*usercenterproto.ValidatePasswordReply, error) {
	args := mocked.Called(ctx, in)

	return args.Get(0).(*usercenterproto.ValidatePasswordReply), args.Error(1)
}
