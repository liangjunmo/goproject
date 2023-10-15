package usercenter

import (
	"context"

	"github.com/liangjunmo/gocode"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/service/usercenterservice"
)

type Server struct {
	usercenterproto.UnimplementedUserCenterServer
	userCenterService usercenterservice.Service
}

func NewServer(userCenterService usercenterservice.Service) *Server {
	return &Server{
		userCenterService: userCenterService,
	}
}

func (server *Server) SearchUser(ctx context.Context, req *usercenterproto.SearchUserRequest) (*usercenterproto.SearchUserReply, error) {
	users, err := server.userCenterService.SearchUser(ctx, usercenterservice.SearchUserRequest{
		Uids:      req.Uids,
		Usernames: req.Usernames,
	})
	if err != nil {
		return &usercenterproto.SearchUserReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).Error(),
				Message: err.Error(),
			},
		}, nil
	}

	rep := &usercenterproto.SearchUserReply{
		Users: make([]*usercenterproto.User, 0, len(users)),
	}

	for _, user := range users {
		rep.Users = append(rep.Users, &usercenterproto.User{
			UID:        user.UID,
			Username:   user.Username,
			CreateTime: user.CreateTime.Unix(),
			UpdateTime: user.UpdateTime.Unix(),
			DeleteTime: user.DeleteTime.Time.Unix(),
		})
	}

	return rep, nil
}

func (server *Server) GetUserByUID(ctx context.Context, req *usercenterproto.GetUserByUIDRequest) (*usercenterproto.GetUserByUIDReply, error) {
	user, err := server.userCenterService.GetUserByUID(ctx, usercenterservice.GetUserByUIDRequest{
		UID: req.UID,
	})
	if err != nil {
		return &usercenterproto.GetUserByUIDReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).Error(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.GetUserByUIDReply{
		User: &usercenterproto.User{
			UID:        user.UID,
			Username:   user.Username,
			CreateTime: user.CreateTime.Unix(),
			UpdateTime: user.UpdateTime.Unix(),
			DeleteTime: user.UpdateTime.Unix(),
		},
	}, nil
}

func (server *Server) GetUserByUsername(ctx context.Context, req *usercenterproto.GetUserByUsernameRequest) (*usercenterproto.GetUserByUsernameReply, error) {
	user, err := server.userCenterService.GetUserByUsername(ctx, usercenterservice.GetUserByUsernameRequest{
		Username: req.Username,
	})
	if err != nil {
		return &usercenterproto.GetUserByUsernameReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).Error(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.GetUserByUsernameReply{
		User: &usercenterproto.User{
			UID:        user.UID,
			Username:   user.Username,
			CreateTime: user.CreateTime.Unix(),
			UpdateTime: user.UpdateTime.Unix(),
			DeleteTime: user.UpdateTime.Unix(),
		},
	}, nil
}

func (server *Server) CreateUser(ctx context.Context, req *usercenterproto.CreateUserRequest) (*usercenterproto.CreateUserReply, error) {
	user, err := server.userCenterService.CreateUser(ctx, usercenterservice.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &usercenterproto.CreateUserReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).Error(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.CreateUserReply{
		User: &usercenterproto.User{
			UID:        user.UID,
			Username:   user.Username,
			CreateTime: user.CreateTime.Unix(),
			UpdateTime: user.UpdateTime.Unix(),
			DeleteTime: user.UpdateTime.Unix(),
		},
	}, nil
}

func (server *Server) ValidatePassword(ctx context.Context, req *usercenterproto.ValidatePasswordRequest) (*usercenterproto.ValidatePasswordReply, error) {
	err := server.userCenterService.ValidatePassword(ctx, usercenterservice.ValidatePasswordRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &usercenterproto.ValidatePasswordReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).Error(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.ValidatePasswordReply{}, nil
}
