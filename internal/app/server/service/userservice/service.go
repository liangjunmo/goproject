package userservice

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/types"
	"github.com/liangjunmo/goproject/internal/pkg/pageutil"
)

type Service interface {
	ListService
	ReadService
	BusinessService
}

type service struct {
	listService     ListService
	readService     ReadService
	businessService BusinessService
}

func NewService(db *gorm.DB, redisSync *redsync.Redsync) Service {
	return &service{
		listService:     newListService(db),
		readService:     newReadService(db),
		businessService: newBusinessService(db, redisSync),
	}
}

func (service *service) ListUser(ctx context.Context, req ListUserRequest) (pageutil.Pagination, []types.User, error) {
	return service.listService.ListUser(ctx, req)
}

func (service *service) SearchUser(ctx context.Context, req SearchUserRequest) ([]types.User, error) {
	return service.readService.SearchUser(ctx, req)
}

func (service *service) GetUser(ctx context.Context, req GetUserRequest) (types.User, error) {
	return service.readService.GetUser(ctx, req)
}

func (service *service) CreateUser(ctx context.Context, req CreateUserRequest) (types.User, error) {
	return service.businessService.CreateUser(ctx, req)
}

func (service *service) ValidatePassword(ctx context.Context, req ValidatePasswordRequest) error {
	return service.businessService.ValidatePassword(ctx, req)
}
