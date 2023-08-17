package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/server/servercode"
)

type ReadService interface {
	SearchUser(ctx context.Context, cmd SearchUserCommand) ([]User, error)
	GetUser(ctx context.Context, cmd GetUserCommand) (User, error)
}

type readService struct {
	db *gorm.DB
}

func NewReadService(db *gorm.DB) ReadService {
	return &readService{
		db: db,
	}
}

func (service *readService) SearchUser(ctx context.Context, cmd SearchUserCommand) ([]User, error) {
	db := service.db.WithContext(ctx).Model(&User{})

	if len(cmd.Uids) != 0 {
		db = db.Where("id in (?)", cmd.Uids)
	}

	if len(cmd.Usernames) != 0 {
		db = db.Where("username in (?)", cmd.Usernames)
	}

	var users []User

	err := db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	return users, nil
}

func (service *readService) GetUser(ctx context.Context, cmd GetUserCommand) (User, error) {
	db := service.db.WithContext(ctx).Model(&User{})

	if cmd.Uid != 0 {
		db = db.Where("id = ?", cmd.Uid)
	}

	if cmd.Username != "" {
		db = db.Where("username = ?", cmd.Username)
	}

	var user User

	err := db.Limit(1).Scan(&user).Error
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	if user.Id == 0 {
		return User{}, servercode.UserNotFound
	}

	return user, nil
}
