package usermanager

import (
	"context"

	"github.com/liangjunmo/goproject/internal/service/userservice"
	"github.com/liangjunmo/goproject/internal/types"
)

type Manager struct {
	userService userservice.Service
}

func NewManager(userService userservice.Service) *Manager {
	return &Manager{
		userService: userService,
	}
}

func (manager *Manager) CreateUser(ctx context.Context, username, password string) (types.User, error) {
	user, err := manager.userService.CreateUser(ctx, userservice.CreateUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return types.User{}, err
	}

	// do other things

	return user, nil
}
