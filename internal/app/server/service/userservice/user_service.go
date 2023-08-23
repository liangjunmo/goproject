package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/servercode"
	"github.com/liangjunmo/goproject/internal/rediskey"
)

type UserService interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error)
	ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error
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
		Password: service.cryptPassword(cmd.Password),
	}

	err = DbCreateUser(ctx, service.db, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service *userService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	user, ok, err := DbGetUserByUsername(ctx, service.db, cmd.Username)
	if err != nil {
		return err
	}

	if !ok {
		return servercode.UserNotFound
	}

	if !service.comparePassword(user.Password, cmd.Password) {
		return servercode.LoginPasswordWrong
	}

	return nil
}

func (service *userService) cryptPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func (service *userService) comparePassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
