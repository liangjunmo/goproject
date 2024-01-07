package userserviceimpl

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice"
)

func ProvideService(db *gorm.DB, redisClient *redis.Client) userservice.Service {
	return newDefaultService(
		newDefaultRepository(db),
		newDefaultMutexProvider(redsync.New(goredis.NewPool(redisClient))),
	)
}

type defaultService struct {
	log           *logrus.Entry
	repository    repository
	mutexProvider mutexProvider
}

func newDefaultService(repository repository, mutexProvider mutexProvider) *defaultService {
	return &defaultService{
		log:           logrus.WithField("tag", "usercenter.userservice.service"),
		repository:    repository,
		mutexProvider: mutexProvider,
	}
}

func (service *defaultService) Search(ctx context.Context, cmd userservice.SearchCommand) (map[uint32]userservice.User, error) {
	users, err := service.repository.Search(ctx, criteria{
		uids:     cmd.Uids,
		username: cmd.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return users, nil
}

func (service *defaultService) Get(ctx context.Context, cmd userservice.GetCommand) (userservice.User, error) {
	user, exist, err := service.repository.Get(ctx, cmd.UID)
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return userservice.User{}, codes.UserNotFound
	}

	return user, nil
}

func (service *defaultService) GetByUsername(ctx context.Context, cmd userservice.GetByUsernameCommand) (userservice.User, error) {
	user, exist, err := service.repository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return userservice.User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return userservice.User{}, codes.UserNotFound
	}

	return user, nil
}

func (service *defaultService) Create(ctx context.Context, cmd userservice.CreateCommand) (uid uint32, err error) {
	m := service.mutexProvider.ProvideCreateUserMutex(cmd.Username)

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

	user, exist, err := service.repository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if exist {
		return user.UID, nil
	}

	user = userservice.User{}
	user.Username = cmd.Username
	user.Password = cryptPassword(cmd.Password)

	err = service.repository.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return user.UID, nil
}

func (service *defaultService) ValidatePassword(ctx context.Context, cmd userservice.ValidatePasswordCommand) error {
	user, exist, err := service.repository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return codes.UserNotFound
	}

	if !comparePassword(user.Password, cmd.Password) {
		return codes.LoginFailedPasswordWrong
	}

	return nil
}

func cryptPassword(plain string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b)
}

func comparePassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
