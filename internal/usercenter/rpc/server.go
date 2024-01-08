package rpc

import (
	"context"

	"github.com/liangjunmo/gocode"
	"github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/usercenter/service"
)

type Server struct {
	usercenterproto.UnimplementedUserCenterServer

	log         *logrus.Entry
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		log:         logrus.WithField("tag", "usercenter.rpc.server"),
		userService: userService,
	}
}

func (server *Server) SearchUser(ctx context.Context, req *usercenterproto.SearchUserRequest) (*usercenterproto.SearchUserReply, error) {
	users, err := server.userService.Search(ctx, service.SearchCommand{
		Uids:     req.Uids,
		Username: req.Username,
	})
	if err != nil {
		server.log.WithContext(ctx).WithError(err).Error(err)

		return &usercenterproto.SearchUserReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).String(),
				Message: err.Error(),
			},
		}, nil
	}

	rep := &usercenterproto.SearchUserReply{
		Users: make(map[uint32]*usercenterproto.User),
	}

	for _, u := range users {
		rep.Users[u.UID] = &usercenterproto.User{
			UID:        u.UID,
			Username:   u.Username,
			CreateTime: u.CreateTime.Unix(),
			UpdateTime: u.UpdateTime.Unix(),
		}
	}

	return rep, nil
}

func (server *Server) GetUserByUID(ctx context.Context, req *usercenterproto.GetUserByUIDRequest) (*usercenterproto.GetUserByUIDReply, error) {
	user, err := server.userService.Get(ctx, service.GetCommand{UID: req.UID})
	if err != nil {
		server.log.WithContext(ctx).WithError(err).Error(err)

		return &usercenterproto.GetUserByUIDReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).String(),
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
		},
	}, nil
}

func (server *Server) GetUserByUsername(ctx context.Context, req *usercenterproto.GetUserByUsernameRequest) (*usercenterproto.GetUserByUsernameReply, error) {
	user, err := server.userService.GetByUsername(ctx, service.GetByUsernameCommand{Username: req.Username})
	if err != nil {
		server.log.WithContext(ctx).WithError(err).Error(err)

		return &usercenterproto.GetUserByUsernameReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).String(),
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
		},
	}, nil
}

func (server *Server) CreateUser(ctx context.Context, req *usercenterproto.CreateUserRequest) (*usercenterproto.CreateUserReply, error) {
	uid, err := server.userService.Create(ctx, service.CreateCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		server.log.WithContext(ctx).WithError(err).Error(err)

		return &usercenterproto.CreateUserReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).String(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.CreateUserReply{
		UID: uid,
	}, nil
}

func (server *Server) ValidatePassword(ctx context.Context, req *usercenterproto.ValidatePasswordRequest) (*usercenterproto.ValidatePasswordReply, error) {
	err := server.userService.ValidatePassword(ctx, service.ValidatePasswordCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		server.log.WithContext(ctx).WithError(err).Error(err)

		return &usercenterproto.ValidatePasswordReply{
			Error: &usercenterproto.Error{
				Code:    gocode.Parse(err).String(),
				Message: err.Error(),
			},
		}, nil
	}

	return &usercenterproto.ValidatePasswordReply{}, nil
}
