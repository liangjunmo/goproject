package userservice

import (
	"context"

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

func (service *service) ListUser(ctx context.Context, params ListUserParams) (pagination.Pagination, []User, error) {
	return service.listService.ListUser(ctx, params)
}

func (service *service) SearchUser(ctx context.Context, params SearchUserParams) ([]User, error) {
	return service.readService.SearchUser(ctx, params)
}

func (service *service) GetUser(ctx context.Context, params GetUserParams) (User, error) {
	return service.readService.GetUser(ctx, params)
}

func (service *service) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	return service.businessService.CreateUser(ctx, params)
}

func (service *service) ValidatePassword(ctx context.Context, params ValidatePasswordParams) error {
	return service.businessService.ValidatePassword(ctx, params)
}
