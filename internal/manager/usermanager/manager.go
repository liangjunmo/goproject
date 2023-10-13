package usermanager

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/helper/userhelper"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
	"github.com/liangjunmo/goproject/internal/service/usercenterservice"
	"github.com/liangjunmo/goproject/internal/service/userservice"
	"github.com/liangjunmo/goproject/internal/types"
)

type Manager struct {
	userCenterService usercenterservice.Service
	userService       userservice.Service
}

func NewManager(userCenterService usercenterservice.Service, userService userservice.Service) *Manager {
	return &Manager{
		userCenterService: userCenterService,
		userService:       userService,
	}
}

func (manager *Manager) ListUser(ctx context.Context, preq pagination.Request) (pagination.Pagination, []types.UserDetail, error) {
	p, users, err := manager.userService.ListUser(ctx, userservice.ListUserRequest{
		PaginationRequest: preq,
	})
	if err != nil {
		return pagination.Pagination{}, nil, err
	}

	if len(users) == 0 {
		return pagination.Pagination{}, nil, nil
	}

	uids := userhelper.FetchUserUids(users)

	ucUsers, err := manager.userCenterService.SearchUser(ctx, usercenterservice.SearchUserRequest{
		Uids: uids,
	})
	if err != nil {
		return pagination.Pagination{}, nil, err
	}

	ucUserMap := userhelper.UserCenterUserToMap(ucUsers)

	userDetailList := make([]types.UserDetail, 0, len(users))

	for _, user := range users {
		userDetailList = append(userDetailList, types.UserDetail{
			UID:        user.UID,
			Username:   ucUserMap[user.UID].Username,
			CreateTime: user.CreateTime,
			UpdateTime: user.UpdateTime,
			DeleteTime: user.DeleteTime,
		})
	}

	return p, userDetailList, nil
}

func (manager *Manager) SearchUser(ctx context.Context, uids []uint32, usernames []string) ([]types.UserDetail, error) {
	if len(uids) == 0 && len(usernames) == 0 {
		return nil, nil
	}

	if len(usernames) != 0 {
		ucUsers, err := manager.userCenterService.SearchUser(ctx, usercenterservice.SearchUserRequest{
			Usernames: usernames,
		})
		if err != nil {
			return nil, err
		}

		uids = userhelper.FetchUserCenterUserUids(ucUsers)
	}

	if len(uids) == 0 {
		return nil, nil
	}

	users, err := manager.userService.SearchUser(ctx, userservice.SearchUserRequest{
		Uids: uids,
	})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	uids = userhelper.FetchUserUids(users)

	ucUsers, err := manager.userCenterService.SearchUser(ctx, usercenterservice.SearchUserRequest{
		Uids: uids,
	})
	if err != nil {
		return nil, err
	}

	ucUserMap := userhelper.UserCenterUserToMap(ucUsers)

	userDetailList := make([]types.UserDetail, 0, len(users))

	for _, user := range users {
		userDetailList = append(userDetailList, types.UserDetail{
			UID:        user.UID,
			Username:   ucUserMap[user.UID].Username,
			CreateTime: user.CreateTime,
			UpdateTime: user.UpdateTime,
			DeleteTime: user.DeleteTime,
		})
	}

	return userDetailList, nil
}

func (manager *Manager) GetUserByUID(ctx context.Context, uid uint32) (types.UserDetail, error) {
	ucUser, err := manager.userCenterService.GetUserByUID(ctx, usercenterservice.GetUserByUIDRequest{
		UID: uid,
	})
	if err != nil {
		return types.UserDetail{}, err
	}

	user, err := manager.userService.GetUserByUID(ctx, userservice.GetUserByUIDRequest{
		UID: uid,
	})
	if err != nil {
		return types.UserDetail{}, err
	}

	return types.UserDetail{
		UID:        user.UID,
		Username:   ucUser.Username,
		CreateTime: user.CreateTime,
		UpdateTime: time.Time{},
		DeleteTime: gorm.DeletedAt{},
	}, nil
}

func (manager *Manager) GetUserByUsername(ctx context.Context, username string) (types.UserDetail, error) {
	ucUser, err := manager.userCenterService.GetUserByUsername(ctx, usercenterservice.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return types.UserDetail{}, err
	}

	user, err := manager.userService.GetUserByUID(ctx, userservice.GetUserByUIDRequest{
		UID: ucUser.UID,
	})
	if err != nil {
		return types.UserDetail{}, err
	}

	return types.UserDetail{
		UID:        user.UID,
		Username:   ucUser.Username,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
		DeleteTime: user.DeleteTime,
	}, nil
}

func (manager *Manager) CreateUser(ctx context.Context, username, password string) (types.User, error) {
	ucUser, err := manager.userCenterService.CreateUser(ctx, usercenterservice.CreateUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return types.User{}, err
	}

	user, err := manager.userService.CreateUser(ctx, userservice.CreateUserRequest{
		UID: ucUser.UID,
	})
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (manager *Manager) ValidatePassword(ctx context.Context, username, password string) error {
	err := manager.userCenterService.ValidatePassword(ctx, usercenterservice.ValidatePasswordRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	return nil
}
