package usercenterservice

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/dbdata"
	"github.com/liangjunmo/goproject/internal/redismutex"
	"github.com/liangjunmo/goproject/internal/types"
)

type BusinessService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (types.UserCenterUser, error)
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

func (service *businessService) CreateUser(ctx context.Context, req CreateUserRequest) (types.UserCenterUser, error) {
	mutex := redismutex.NewCreateUserCenterUserMutex(service.redisSync, req.Username)

	err := mutex.Lock()
	if err != nil {
		return types.UserCenterUser{}, codes.Timeout
	}
	defer mutex.Unlock()

	user, ok, err := dbdata.GetUserCenterUserByUsername(ctx, service.db, req.Username)
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if ok {
		return user, nil
	}

	user = types.UserCenterUser{
		Username: req.Username,
		Password: cryptPassword(req.Password),
	}

	err = dbdata.CreateUserCenterUser(ctx, service.db, &user)
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return user, nil
}

func (service *businessService) ValidatePassword(ctx context.Context, req ValidatePasswordRequest) error {
	user, ok, err := dbdata.GetUserCenterUserByUsername(ctx, service.db, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return codes.UserNotFound
	}

	if !comparePassword(user.Password, req.Password) {
		return codes.LoginPasswordWrong
	}

	return nil
}
