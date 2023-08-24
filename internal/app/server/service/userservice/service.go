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

func (service *service) ListUser(ctx context.Context, cmd ListUserCommand) (pagination.Pagination, []User, error) {
	return service.listService.ListUser(ctx, cmd)
}

func (service *service) SearchUser(ctx context.Context, cmd SearchUserCommand) ([]User, error) {
	return service.readService.SearchUser(ctx, cmd)
}

func (service *service) GetUser(ctx context.Context, cmd GetUserCommand) (User, error) {
	return service.readService.GetUser(ctx, cmd)
}

func (service *service) CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error) {
	return service.businessService.CreateUser(ctx, cmd)
}

func (service *service) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	return service.businessService.ValidatePassword(ctx, cmd)
}
