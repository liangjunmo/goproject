package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/app/server/servercode"
	"github.com/liangjunmo/goproject/internal/app/server/service/userservice"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
	"github.com/liangjunmo/goproject/internal/pkg/timeutil"
)

type UserHandler struct {
	*BaseHandler
	userHubService userservice.HubService
}

func NewUserHandler(userHubService userservice.HubService) *UserHandler {
	return &UserHandler{
		userHubService: userHubService,
	}
}

type ListUserRequest struct {
	pagination.Request
}

type ListUserResponse struct {
	pagination.Pagination
	List []ListUserData `json:"list"`
}

type ListUserData struct {
	Uid        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *UserHandler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()
	l := log.WithContext(ctx)

	var (
		req  ListUserRequest
		resp ListUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	p, users, err := handler.userHubService.ListUser(ctx, userservice.ListUserCommand{
		PaginationRequest: req.Request,
	})
	if err != nil {
		l.WithError(err).Error(err)
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
			Uid:        user.Id,
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
	l := log.WithContext(ctx)

	var (
		req  SearchUserRequest
		resp SearchUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	users, err := handler.userHubService.SearchUser(ctx, userservice.SearchUserCommand{
		Uids:      req.Uids,
		Usernames: req.Usernames,
	})
	if err != nil {
		l.WithError(err).Error(err)
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
			Uid:        user.Id,
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
	l := log.WithContext(ctx)

	var (
		req  GetUserRequest
		resp GetUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	user, err := handler.userHubService.GetUser(ctx, userservice.GetUserCommand{
		Uid:      req.Uid,
		Username: req.Username,
	})
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = GetUserResponse{
		Uid:        user.Id,
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
	l := log.WithContext(ctx)

	var (
		req  CreateUserRequest
		resp CreateUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	user, err := handler.userHubService.CreateUser(ctx, userservice.CreateUserCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		l.WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = CreateUserResponse{
		Uid: user.Id,
	}

	handler.Response(c, resp, nil)
}