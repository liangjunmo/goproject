package usercenter

import (
	"context"
	"fmt"
	"time"

	"github.com/liangjunmo/gocode"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/service/usercenterservice"
	"github.com/liangjunmo/goproject/internal/types"
)

type Adapter struct {
	client usercenterproto.UserCenterClient
}

func NewAdapter(client usercenterproto.UserCenterClient) *Adapter {
	return &Adapter{
		client: client,
	}
}

func (adapter *Adapter) SearchUser(ctx context.Context, req usercenterservice.SearchUserRequest) ([]types.UserCenterUser, error) {
	rep, err := adapter.client.SearchUser(ctx, &usercenterproto.SearchUserRequest{
		Uids:      req.Uids,
		Usernames: req.Usernames,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return nil, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	users := make([]types.UserCenterUser, 0, len(rep.Users))

	for _, u := range rep.Users {
		var deleteTime gorm.DeletedAt

		if u.DeleteTime != 0 {
			deleteTime.Time = time.Unix(u.DeleteTime, 0)
		}
		users = append(users, types.UserCenterUser{
			UID:        u.UID,
			CreateTime: time.Unix(u.CreateTime, 0),
			UpdateTime: time.Unix(u.UpdateTime, 0),
			DeleteTime: deleteTime,
			Username:   u.Username,
		})
	}

	return users, nil
}

func (adapter *Adapter) GetUserByUID(ctx context.Context, req usercenterservice.GetUserByUIDRequest) (types.UserCenterUser, error) {
	rep, err := adapter.client.GetUserByUID(ctx, &usercenterproto.GetUserByUIDRequest{
		UID: req.UID,
	})
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	var deleteTime gorm.DeletedAt

	if rep.User.DeleteTime != 0 {
		deleteTime.Time = time.Unix(rep.User.DeleteTime, 0)
	}

	return types.UserCenterUser{
		UID:        rep.User.UID,
		CreateTime: time.Unix(rep.User.CreateTime, 0),
		UpdateTime: time.Unix(rep.User.UpdateTime, 0),
		DeleteTime: deleteTime,
		Username:   rep.User.Username,
	}, nil
}

func (adapter *Adapter) GetUserByUsername(ctx context.Context, req usercenterservice.GetUserByUsernameRequest) (types.UserCenterUser, error) {
	rep, err := adapter.client.GetUserByUsername(ctx, &usercenterproto.GetUserByUsernameRequest{
		Username: req.Username,
	})
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	var deleteTime gorm.DeletedAt

	if rep.User.DeleteTime != 0 {
		deleteTime.Time = time.Unix(rep.User.DeleteTime, 0)
	}

	return types.UserCenterUser{
		UID:        rep.User.UID,
		CreateTime: time.Unix(rep.User.CreateTime, 0),
		UpdateTime: time.Unix(rep.User.UpdateTime, 0),
		DeleteTime: deleteTime,
		Username:   rep.User.Username,
	}, nil
}

func (adapter *Adapter) CreateUser(ctx context.Context, req usercenterservice.CreateUserRequest) (types.UserCenterUser, error) {
	rep, err := adapter.client.CreateUser(ctx, &usercenterproto.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return types.UserCenterUser{}, fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	var deleteTime gorm.DeletedAt

	if rep.User.DeleteTime != 0 {
		deleteTime.Time = time.Unix(rep.User.DeleteTime, 0)
	}

	return types.UserCenterUser{
		UID:        rep.User.UID,
		CreateTime: time.Unix(rep.User.CreateTime, 0),
		UpdateTime: time.Unix(rep.User.UpdateTime, 0),
		DeleteTime: deleteTime,
		Username:   rep.User.Username,
	}, nil
}

func (adapter *Adapter) ValidatePassword(ctx context.Context, req usercenterservice.ValidatePasswordRequest) error {
	rep, err := adapter.client.ValidatePassword(ctx, &usercenterproto.ValidatePasswordRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if rep.Error != nil {
		return fmt.Errorf("%w: %s", gocode.Code(rep.Error.Code), rep.Error.Message)
	}

	return nil
}
