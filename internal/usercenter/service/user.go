package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/usercenter/model"
	"github.com/liangjunmo/goproject/internal/usercenter/mutex"
	"github.com/liangjunmo/goproject/internal/usercenter/repository"
)

type UserService interface {
	Search(ctx context.Context, cmd SearchCommand) (map[uint32]model.User, error)
	Get(ctx context.Context, cmd GetCommand) (model.User, error)
	GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (model.User, error)
	Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

func NewUserService(mutexProvider mutex.MutexProvider, userRepository repository.UserRepository) UserService {
	return newUserService(mutexProvider, userRepository)
}

type userService struct {
	log            *logrus.Entry
	mutexProvider  mutex.MutexProvider
	userRepository repository.UserRepository
}

func newUserService(mutexProvider mutex.MutexProvider, userRepository repository.UserRepository) *userService {
	return &userService{
		log:            logrus.WithField("tag", "usercenter.user_service"),
		mutexProvider:  mutexProvider,
		userRepository: userRepository,
	}
}

type SearchCommand struct {
	Uids     []uint32
	Username string
}

func (service *userService) Search(ctx context.Context, cmd SearchCommand) (map[uint32]model.User, error) {
	users, err := service.userRepository.Search(ctx, repository.UserCriteria{
		Uids:     cmd.Uids,
		Username: cmd.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return users, nil
}

type GetCommand struct {
	UID uint32
}

func (service *userService) Get(ctx context.Context, cmd GetCommand) (model.User, error) {
	user, exist, err := service.userRepository.Get(ctx, cmd.UID)
	if err != nil {
		return model.User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return model.User{}, codes.UserNotFound
	}

	return user, nil
}

type GetByUsernameCommand struct {
	Username string
}

func (service *userService) GetByUsername(ctx context.Context, cmd GetByUsernameCommand) (model.User, error) {
	user, exist, err := service.userRepository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return model.User{}, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if !exist {
		return model.User{}, codes.UserNotFound
	}

	return user, nil
}

type CreateCommand struct {
	Username string
	Password string
}

func (service *userService) Create(ctx context.Context, cmd CreateCommand) (uid uint32, err error) {
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

	user, exist, err := service.userRepository.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	if exist {
		return user.UID, nil
	}

	user = model.User{}
	user.Username = cmd.Username
	user.Password = cryptPassword(cmd.Password)

	err = service.userRepository.Create(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", codes.InternalServerError, err)
	}

	return user.UID, nil
}

type ValidatePasswordCommand struct {
	Username string
	Password string
}

func (service *userService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	user, exist, err := service.userRepository.GetByUsername(ctx, cmd.Username)
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
