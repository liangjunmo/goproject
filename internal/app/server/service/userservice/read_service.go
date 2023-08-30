package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/app/server/types"
)

type ReadService interface {
	SearchUser(ctx context.Context, req SearchUserRequest) ([]types.User, error)
	GetUser(ctx context.Context, req GetUserRequest) (types.User, error)
}

type readService struct {
	db *gorm.DB
}

func newReadService(db *gorm.DB) ReadService {
	return &readService{
		db: db,
	}
}

func (service *readService) SearchUser(ctx context.Context, req SearchUserRequest) ([]types.User, error) {
	db := service.db.WithContext(ctx).Model(&types.User{})

	if len(req.Uids) != 0 {
		db = db.Where("id in (?)", req.Uids)
	}

	if len(req.Usernames) != 0 {
		db = db.Where("username in (?)", req.Usernames)
	}

	var users []types.User

	err := db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return users, nil
}

func (service *readService) GetUser(ctx context.Context, req GetUserRequest) (types.User, error) {
	db := service.db.WithContext(ctx).Model(&types.User{})

	if req.Uid != 0 {
		db = db.Where("id = ?", req.Uid)
	}

	if req.Username != "" {
		db = db.Where("username = ?", req.Username)
	}

	var user types.User

	err := db.Limit(1).Scan(&user).Error
	if err != nil {
		return types.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if user.Uid == 0 {
		return types.User{}, codes.UserNotFound
	}

	return user, nil
}
