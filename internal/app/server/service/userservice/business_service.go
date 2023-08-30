package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/rediskey"
)

type BusinessService interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	ValidatePassword(ctx context.Context, params ValidatePasswordParams) error
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

func (service *businessService) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	mutex := service.redisSync.NewMutex(
		rediskey.MutexCreateUser(params.Username),
		redsync.WithTries(1),
	)

	err := mutex.Lock()
	if err != nil {
		return User{}, codes.Timeout
	}
	defer mutex.Unlock()

	user, ok, err := DbGetUserByUsername(ctx, service.db, params.Username)
	if err != nil {
		return User{}, err
	}

	if ok {
		return User{}, codes.UserAlreadyExists
	}

	user = User{
		Username: params.Username,
		Password: cryptPassword(params.Password),
	}

	err = DbCreateUser(ctx, service.db, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service *businessService) ValidatePassword(ctx context.Context, params ValidatePasswordParams) error {
	user, ok, err := DbGetUserByUsername(ctx, service.db, params.Username)
	if err != nil {
		return err
	}

	if !ok {
		return codes.UserNotFound
	}

	if !comparePassword(user.Password, params.Password) {
		return codes.LoginPasswordWrong
	}

	return nil
}
