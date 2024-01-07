package userserviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/goproject/service/userservice"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

func TestDefaultService(t *testing.T) {
	var (
		mutex            *mockedMutex
		mutexProvider    *mockedMutexProvider
		repository       *mockedRepository
		userCenterClient *mockedUserCenterClient
		service          *defaultService
		ctx              context.Context
	)

	beforeTest := func(t *testing.T) {
		mutex = &mockedMutex{}
		mutexProvider = &mockedMutexProvider{}
		repository = &mockedRepository{}
		userCenterClient = &mockedUserCenterClient{}

		service = newDefaultService(repository, mutexProvider, userCenterClient)

		ctx = context.Background()
	}

	t.Run("List", func(t *testing.T) {
		beforeTest(t)

		repository.
			On("List", ctx, mock.IsType(criteria{})).
			Return(pagination.DefaultPagination{}, []userservice.User{{UID: 1}}, nil)

		userCenterClient.
			On("SearchUser", ctx, mock.IsType(&usercenterproto.SearchUserRequest{})).
			Return(
				&usercenterproto.SearchUserReply{
					Users: map[uint32]*usercenterproto.User{1: {UID: 1}},
				},
				nil,
			)

		_, users, err := service.List(ctx, userservice.ListCommand{})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[0].UID)
		require.Equal(t, uint32(1), users[0].UserCenterUser.UID)
	})

	t.Run("Search", func(t *testing.T) {
		beforeTest(t)

		userCenterClient.
			On("SearchUser", ctx, mock.IsType(&usercenterproto.SearchUserRequest{})).
			Return(
				&usercenterproto.SearchUserReply{
					Users: map[uint32]*usercenterproto.User{1: {UID: 1}},
				},
				nil,
			)

		repository.
			On("Search", ctx, mock.IsType(criteria{})).
			Return(map[uint32]userservice.User{1: {UID: 1}}, nil)

		users, err := service.Search(ctx, userservice.SearchCommand{Uids: []uint32{1}})
		require.Nil(t, err)
		require.Equal(t, uint32(1), users[1].UID)
		require.Equal(t, uint32(1), users[1].UserCenterUser.UID)
	})

	t.Run("Get", func(t *testing.T) {
		beforeTest(t)

		repository.On("Get", ctx, uint32(1)).Return(userservice.User{UID: 1}, true, nil)

		userCenterClient.
			On("GetUserByUID", ctx, mock.IsType(&usercenterproto.GetUserByUIDRequest{})).
			Return(
				&usercenterproto.GetUserByUIDReply{
					User: &usercenterproto.User{UID: 1},
				},
				nil,
			)

		user, err := service.Get(ctx, userservice.GetCommand{UID: 1})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, uint32(1), user.UserCenterUser.UID)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		beforeTest(t)

		userCenterClient.
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

		repository.On("Get", ctx, uint32(1)).Return(userservice.User{UID: 1}, true, nil)

		user, err := service.GetByUsername(ctx, userservice.GetByUsernameCommand{Username: "user"})
		require.Nil(t, err)
		require.Equal(t, uint32(1), user.UID)
		require.Equal(t, "user", user.UserCenterUser.Username)
	})

	t.Run("Create", func(t *testing.T) {
		beforeTest(t)

		userCenterClient.
			On("CreateUser", ctx, mock.IsType(&usercenterproto.CreateUserRequest{})).
			Return(
				&usercenterproto.CreateUserReply{
					UID: 1,
				},
				nil,
			)

		mutex.On("Lock", ctx).Return(nil)
		mutex.On("Unlock", ctx).Return(true, nil)

		mutexProvider.On("ProvideCreateUserMutex", uint32(1)).Return(mutex)

		repository.On("Get", ctx, uint32(1)).Return(userservice.User{}, false, nil)
		repository.On("Create", ctx, mock.IsType(&userservice.User{})).Return(nil)

		uid, err := service.Create(ctx, userservice.CreateCommand{
			Username: "user",
			Password: "pass",
		})
		require.Nil(t, err)
		require.Equal(t, uint32(1), uid)
	})

	t.Run("ValidatePassword", func(t *testing.T) {
		beforeTest(t)

		userCenterClient.
			On("ValidatePassword", ctx, mock.IsType(&usercenterproto.ValidatePasswordRequest{})).
			Return(&usercenterproto.ValidatePasswordReply{}, nil)

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

func (mocked *mockedMutexProvider) ProvideCreateUserMutex(uid uint32) mutex {
	args := mocked.Called(uid)

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

func (mocked *mockedRepository) List(ctx context.Context, criteria criteria) (pagination.Pagination, []userservice.User, error) {
	args := mocked.Called(ctx, criteria)

	if args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(pagination.Pagination), args.Get(1).([]userservice.User), args.Error(2)
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

func (mocked *mockedRepository) Create(ctx context.Context, user *userservice.User) error {
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
