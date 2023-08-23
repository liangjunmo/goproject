package userservice

import (
	"context"

	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type HubService interface {
	ListService
	ReadService
	UserService
}

type hubService struct {
	listService ListService
	readService ReadService
	userService UserService
}

func NewHubService(listService ListService, readService ReadService, userService UserService) HubService {
	return &hubService{
		listService: listService,
		readService: readService,
		userService: userService,
	}
}

func (service *hubService) ListUser(ctx context.Context, cmd ListUserCommand) (pagination.Pagination, []User, error) {
	return service.listService.ListUser(ctx, cmd)
}

func (service *hubService) SearchUser(ctx context.Context, cmd SearchUserCommand) ([]User, error) {
	return service.readService.SearchUser(ctx, cmd)
}

func (service *hubService) GetUser(ctx context.Context, cmd GetUserCommand) (User, error) {
	return service.readService.GetUser(ctx, cmd)
}

func (service *hubService) CreateUser(ctx context.Context, cmd CreateUserCommand) (User, error) {
	return service.userService.CreateUser(ctx, cmd)
}

func (service *hubService) ValidatePassword(ctx context.Context, cmd ValidatePasswordCommand) error {
	return service.userService.ValidatePassword(ctx, cmd)
}
