package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/servercode"
	"github.com/liangjunmo/goproject/internal/rediskey"
)

type BusinessService interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
}

type businessService struct {
	db        *gorm.DB
	redisSync *redsync.Redsync
}

func NewBusinessService(db *gorm.DB, redisSync *redsync.Redsync) BusinessService {
	return &businessService{
		db:        db,
		redisSync: redisSync,
	}
}

func (service *businessService) CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error) {
	mutex := service.redisSync.NewMutex(
		rediskey.MutexCreateUser(cmd.Username),
		redsync.WithTries(1),
	)

	err := mutex.Lock()
	if err != nil {
		return User{}, servercode.Timeout
	}
	defer mutex.Unlock()

	user, ok, err := DbGetUserByUsername(ctx, service.db, cmd.Username)
	if err != nil {
		return User{}, err
	}

	if ok {
		return User{}, servercode.UserAlreadyExists
	}

	user = User{
		Username: cmd.Username,
		Password: cryptPassword(cmd.Password),
	}

	err = DbCreateUser(ctx, service.db, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service *businessService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	user, ok, err := DbGetUserByUsername(ctx, service.db, cmd.Username)
	if err != nil {
		return err
	}

	if !ok {
		return servercode.UserNotFound
	}

	if !comparePassword(user.Password, cmd.Password) {
		return servercode.LoginPasswordWrong
	}

	return nil
}
