package usercenterservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/types"
)

type Service interface {
	ReadService
	BusinessService
}

type service struct {
	readService     ReadService
	businessService BusinessService
}

func NewService(db *gorm.DB, redisSync *redsync.Redsync) Service {
	return &service{
		readService:     newReadService(db),
		businessService: newBusinessService(db, redisSync),
	}
}

func (service *service) SearchUser(ctx context.Context, req SearchUserRequest) ([]types.UserCenterUser, error) {
	return service.readService.SearchUser(ctx, req)
}

func (service *service) GetUserByUID(ctx context.Context, req GetUserByUIDRequest) (types.UserCenterUser, error) {
	return service.readService.GetUserByUID(ctx, req)
}

func (service *service) GetUserByUsername(ctx context.Context, req GetUserByUsernameRequest) (types.UserCenterUser, error) {
	return service.readService.GetUserByUsername(ctx, req)
}

func (service *service) CreateUser(ctx context.Context, req CreateUserRequest) (types.UserCenterUser, error) {
	return service.businessService.CreateUser(ctx, req)
}

func (service *service) ValidatePassword(ctx context.Context, req ValidatePasswordRequest) error {
	return service.businessService.ValidatePassword(ctx, req)
}