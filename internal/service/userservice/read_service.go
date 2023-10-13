package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/dbdata"
	"github.com/liangjunmo/goproject/internal/types"
)

type ReadService interface {
	SearchUser(ctx context.Context, req SearchUserRequest) ([]types.User, error)
	GetUserByUID(ctx context.Context, req GetUserByUIDRequest) (types.User, error)
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
	if len(req.Uids) == 0 {
		return nil, nil
	}

	db := service.db.WithContext(ctx).Model(&types.User{})

	if len(req.Uids) != 0 {
		db = db.Where("uid in (?)", req.Uids)
	}

	var users []types.User

	err := db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return users, nil
}

func (service *readService) GetUserByUID(ctx context.Context, req GetUserByUIDRequest) (types.User, error) {
	user, ok, err := dbdata.GetUserByUID(ctx, service.db, req.UID)
	if err != nil {
		return types.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return types.User{}, codes.UserNotFound
	}

	return user, nil
}
