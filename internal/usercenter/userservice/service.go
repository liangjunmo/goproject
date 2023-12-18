package userservice

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
)

type Service interface {
	Search(ctx context.Context, cmd SearchCommand) (map[uint32]User, error)
	Get(ctx context.Context, cmd GetCommand) (User, error)
	GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (User, error)
	Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

func ProvideService(db *gorm.DB, redisClient *redis.Client) Service {
	repository := newDefaultRepository(db)
	mutexProvider := newDefaultMutexProvider(redsync.New(goredis.NewPool(redisClient)))

	return newDefaultService(repository, mutexProvider)
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

type SearchCommand struct {
	Uids     []uint32
	Username string
}

func (service *defaultService) Search(ctx context.Context, cmd SearchCommand) (map[uint32]User, error) {
	users, err := service.repository.Search(ctx, criteria{
		uids:     cmd.Uids,
		username: cmd.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return users, nil
}

type GetCommand struct {
	UID uint32
}

func (service *defaultService) Get(ctx context.Context, cmd GetCommand) (User, error) {
	user, exist, err := service.repository.Get(ctx, cmd.UID)
	if err != nil {
		return User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return User{}, codes.UserNotFound
	}

	return user, nil
}

type GetByUsernameCommand struct {
	Username string
}

func (service *defaultService) GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (User, error) {
	user, exist, err := service.repository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return User{}, codes.UserNotFound
	}

	return user, nil
}

type CreateCommand struct {
	Username string
	Password string
}

func (service *defaultService) Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error) {
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

	user = User{}
	user.Username = cmd.Username
	user.Password = cryptPassword(cmd.Password)

	err = service.repository.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return user.UID, nil
}

type ValidatePasswordCommand struct {
	Username string
	Password string
}

func (service *defaultService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
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
