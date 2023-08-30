package userservice

import (
	"context"

	"github.com/liangjunmo/goproject/internal/app/server/types"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
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

func NewService(listService ListService, readService ReadService, businessService BusinessService) Service {
	return &service{
		listService:     listService,
		readService:     readService,
		businessService: businessService,
	}
}

func (service *service) ListUser(ctx context.Context, req ListUserRequest) (pagination.Pagination, []types.User, error) {
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
