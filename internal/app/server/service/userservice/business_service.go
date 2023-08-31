package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/app/server/datautil"
	"github.com/liangjunmo/goproject/internal/app/server/types"
	"github.com/liangjunmo/goproject/internal/rediskey"
)

type BusinessService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (types.User, error)
	ValidatePassword(ctx context.Context, req ValidatePasswordRequest) error
}

type businessService struct {
	db        *gorm.DB
	redisSync *redsync.Redsync
}

func newBusinessService(db *gorm.DB, redisSync *redsync.Redsync) BusinessService {
	return &businessService{
		db:        db,
		redisSync: redisSync,
	}
}

func (service *businessService) CreateUser(ctx context.Context, req CreateUserRequest) (types.User, error) {
	mutex := service.redisSync.NewMutex(
		rediskey.MutexCreateUser(req.Username),
		redsync.WithTries(1),
	)

	err := mutex.Lock()
	if err != nil {
		return types.User{}, codes.Timeout
	}
	defer mutex.Unlock()

	user, ok, err := datautil.DbGetUserByUsername(ctx, service.db, req.Username)
	if err != nil {
		return types.User{}, err
	}

	if ok {
		return types.User{}, codes.UserAlreadyExists
	}

	user = types.User{
		Username: req.Username,
		Password: cryptPassword(req.Password),
	}

	err = datautil.DbCreateUser(ctx, service.db, &user)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (service *businessService) ValidatePassword(ctx context.Context, req ValidatePasswordRequest) error {
	user, ok, err := datautil.DbGetUserByUsername(ctx, service.db, req.Username)
	if err != nil {
		return err
	}

	if !ok {
		return codes.UserNotFound
	}

	if !comparePassword(user.Password, req.Password) {
		return codes.LoginPasswordWrong
	}

	return nil
}
