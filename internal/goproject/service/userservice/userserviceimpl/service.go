package userserviceimpl

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/liangjunmo/gocode"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/goproject/service/userservice"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

func ProvideService(db *gorm.DB, redisClient *redis.Client, userCenterClient usercenterproto.UserCenterClient) userservice.Service {
	return newDefaultService(
		newDefaultRepository(db),
		newDefaultMutexProvider(redsync.New(goredis.NewPool(redisClient))),
		userCenterClient,
	)
}

type defaultService struct {
	log              *logrus.Entry
	repository       repository
	mutexProvider    mutexProvider
	userCenterClient usercenterproto.UserCenterClient
}

func newDefaultService(repository repository, mutexProvider mutexProvider, userCenterClient usercenterproto.UserCenterClient) *defaultService {
	return &defaultService{
		log:              logrus.WithField("tag", "goproject.userservice.service"),
		repository:       repository,
		mutexProvider:    mutexProvider,
		userCenterClient: userCenterClient,
	}
}

func (service *defaultService) List(ctx context.Context, cmd userservice.ListCommand) (pagination.Pagination, []userservice.User, error) {
	p, users, err := service.repository.List(ctx, criteria{
		Request: cmd.Request,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	uids := make([]uint32, 0, len(users))

	for _, u := range users {
		uids = append(uids, u.UID)
	}

	rep, err := service.userCenterClient.SearchUser(ctx, &usercenterproto.SearchUserRequest{
		Uids: uids,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return nil, nil, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	for i := range users {
		users[i].UserCenterUser = rep.Users[users[i].UID]
	}

	return p, users, nil
}

func (service *defaultService) Search(ctx context.Context, cmd userservice.SearchCommand) (map[uint32]userservice.User, error) {
	if len(cmd.Uids) == 0 && cmd.Username == "" {
		return nil, nil
	}

	rep, err := service.userCenterClient.SearchUser(ctx, &usercenterproto.SearchUserRequest{
		Uids:     cmd.Uids,
		Username: cmd.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return nil, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	if len(rep.Users) == 0 {
		return nil, nil
	}

	uids := make([]uint32, 0, len(rep.Users))

	for uid := range rep.Users {
		uids = append(uids, uid)
	}

	users, err := service.repository.Search(ctx, criteria{
		uids: uids,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	for uid, u := range users {
		u.UserCenterUser = rep.Users[uid]
		users[uid] = u
	}

	return users, nil
}

func (service *defaultService) Get(ctx context.Context, cmd userservice.GetCommand) (userservice.User, error) {
	user, exist, err := service.repository.Get(ctx, cmd.UID)
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !exist {
		return userservice.User{}, codes.UserNotFound
	}

	rep, err := service.userCenterClient.GetUserByUID(ctx, &usercenterproto.GetUserByUIDRequest{
		UID: user.UID,
	})
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return userservice.User{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	user.UserCenterUser = rep.User

	return user, nil
}

func (service *defaultService) GetByUsername(ctx context.Context, cmd userservice.GetByUsernameCommand) (userservice.User, error) {
	rep, err := service.userCenterClient.GetUserByUsername(ctx, &usercenterproto.GetUserByUsernameRequest{
		Username: cmd.Username,
	})
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return userservice.User{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	user, exist, err := service.repository.Get(ctx, rep.User.UID)
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !exist {
		return userservice.User{}, codes.UserNotFound
	}

	user.UserCenterUser = rep.User

	return user, nil
}

func (service *defaultService) Create(ctx context.Context, cmd userservice.CreateCommand) (uid uint32, err error) {
	rep, err := service.userCenterClient.CreateUser(ctx, &usercenterproto.CreateUserRequest{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return 0, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	m := service.mutexProvider.ProvideCreateUserMutex(rep.UID)

	err = m.Lock(ctx)
	if err != nil {
		return 0, codes.Timeout
	}
	defer func() {
		ok, err := m.Unlock(ctx)
		if !ok || err != nil {
			service.log.WithContext(ctx).WithError(err).Error("unlock")
		}
	}()

	_, exist, err := service.repository.Get(ctx, rep.UID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if exist {
		return 0, codes.UserAlreadyExists
	}

	user := userservice.User{
		UID: rep.UID,
	}

	err = service.repository.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return user.UID, nil
}

func (service *defaultService) ValidatePassword(ctx context.Context, cmd userservice.ValidatePasswordCommand) error {
	rep, err := service.userCenterClient.ValidatePassword(ctx, &usercenterproto.ValidatePasswordRequest{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	return nil
}

func (service *defaultService) taskToRunExample(ctx context.Context, log *logrus.Entry) {
	log.WithContext(ctx).Info("runExample")
}
