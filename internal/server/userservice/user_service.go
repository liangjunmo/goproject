package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/rediskey"
	"github.com/liangjunmo/goproject/internal/server/servercode"
)

type UserService interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error)
}

type userService struct {
	db        *gorm.DB
	redisSync *redsync.Redsync
}

func NewUserService(db *gorm.DB, redisSync *redsync.Redsync) UserService {
	return &userService{
		db:        db,
		redisSync: redisSync,
	}
}

func (service *userService) CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error) {
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
	}

	err = DbCreateUser(ctx, service.db, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
