package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/app/server/manager/usermanager"
	"github.com/liangjunmo/goproject/internal/app/server/service/userservice"
	"github.com/liangjunmo/goproject/internal/pkg/pageutil"
	"github.com/liangjunmo/goproject/internal/pkg/timeutil"
)

type UserHandler struct {
	*BaseHandler
	userService userservice.Service
	userManager *usermanager.Manager
}

func NewUserHandler(userService userservice.Service, userManager *usermanager.Manager) *UserHandler {
	return &UserHandler{
		userService: userService,
		userManager: userManager,
	}
}

type ListUserRequest struct {
	pageutil.Request
}

type ListUserResponse struct {
	pageutil.Pagination
	List []ListUserData `json:"list"`
}

type ListUserData struct {
	Uid        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *UserHandler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  ListUserRequest
		resp ListUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	p, users, err := handler.userService.ListUser(ctx, userservice.ListUserRequest{
		PaginationRequest: req.Request,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = ListUserResponse{
		Pagination: p,
		List:       make([]ListUserData, 0, len(users)),
	}

	if len(users) == 0 {
		handler.Response(c, resp, nil)
		return
	}

	for _, user := range users {
		resp.List = append(resp.List, ListUserData{
			Uid:        user.Uid,
			Username:   user.Username,
			CreateTime: user.CreateTime.Format(timeutil.LayoutTime),
		})
	}

	handler.Response(c, resp, nil)
}

type SearchUserRequest struct {
	Uids      []uint32 `json:"uids"`
	Usernames []string `json:"usernames"`
}

type SearchUserResponse []SearchUserData

type SearchUserData struct {
	Uid        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *UserHandler) SearchUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  SearchUserRequest
		resp SearchUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	users, err := handler.userService.SearchUser(ctx, userservice.SearchUserRequest{
		Uids:      req.Uids,
		Usernames: req.Usernames,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = make([]SearchUserData, 0, len(users))

	if len(users) == 0 {
		handler.Response(c, resp, nil)
		return
	}

	for _, user := range users {
		resp = append(resp, SearchUserData{
			Uid:        user.Uid,
			Username:   user.Username,
			CreateTime: user.CreateTime.Format(timeutil.LayoutTime),
		})
	}

	handler.Response(c, resp, nil)
}

type GetUserRequest struct {
	Uid      uint32 `json:"uid"`
	Username string `json:"username"`
}

type GetUserResponse struct {
	Uid        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  GetUserRequest
		resp GetUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	user, err := handler.userService.GetUser(ctx, userservice.GetUserRequest{
		Uid:      req.Uid,
		Username: req.Username,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = GetUserResponse{
		Uid:        user.Uid,
		Username:   user.Username,
		CreateTime: user.CreateTime.Format(timeutil.LayoutTime),
	}

	handler.Response(c, resp, nil)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	Uid uint32 `json:"uid"`
}

func (handler *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  CreateUserRequest
		resp CreateUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	user, err := handler.userManager.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = CreateUserResponse{
		Uid: user.Uid,
	}

	handler.Response(c, resp, nil)
}
