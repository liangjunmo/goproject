package usercenterservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/dbdata"
	"github.com/liangjunmo/goproject/internal/types"
)

type ReadService interface {
	SearchUser(ctx context.Context, req SearchUserRequest) ([]types.UserCenterUser, error)
	GetUserByUID(ctx context.Context, req GetUserByUIDRequest) (types.UserCenterUser, error)
	GetUserByUsername(ctx context.Context, req GetUserByUsernameRequest) (types.UserCenterUser, error)
}

type readService struct {
	db *gorm.DB
}

func newReadService(db *gorm.DB) ReadService {
	return &readService{
		db: db,
	}
}

func (service *readService) SearchUser(ctx context.Context, req SearchUserRequest) ([]types.UserCenterUser, error) {
	if len(req.Uids) == 0 && len(req.Usernames) == 0 {
		return nil, nil
	}

	db := service.db.WithContext(ctx).Model(&types.UserCenterUser{})

	if len(req.Uids) != 0 {
		db = db.Where("id in (?)", req.Uids)
	}

	if len(req.Usernames) != 0 {
		db = db.Where("username in (?)", req.Usernames)
	}

	var users []types.UserCenterUser

	err := db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return users, nil
}

func (service *readService) GetUserByUID(ctx context.Context, req GetUserByUIDRequest) (types.UserCenterUser, error) {
	user, ok, err := dbdata.GetUserCenterUserByUID(ctx, service.db, req.UID)
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return types.UserCenterUser{}, codes.UserNotFound
	}

	return user, nil
}

func (service *readService) GetUserByUsername(ctx context.Context, req GetUserByUsernameRequest) (types.UserCenterUser, error) {
	user, ok, err := dbdata.GetUserCenterUserByUsername(ctx, service.db, req.Username)
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return types.UserCenterUser{}, codes.UserNotFound
	}

	return user, nil
}
