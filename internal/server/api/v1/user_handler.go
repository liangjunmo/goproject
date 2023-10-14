package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/manager"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
	"github.com/liangjunmo/goproject/internal/pkg/timeutil"
)

type UserHandler struct {
	*BaseHandler
	userManager *manager.UserManager
}

func NewUserHandler(userManager *manager.UserManager) *UserHandler {
	return &UserHandler{
		userManager: userManager,
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
	UID        uint32 `json:"uid"`
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

	p, users, err := handler.userManager.ListUser(ctx, req.Request)
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
			UID:        user.UID,
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
	UID        uint32 `json:"uid"`
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

	users, err := handler.userManager.SearchUser(ctx, req.Uids, req.Usernames)
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
			UID:        user.UID,
			Username:   user.Username,
			CreateTime: user.CreateTime.Format(timeutil.LayoutTime),
		})
	}

	handler.Response(c, resp, nil)
}

type GetUserRequest struct {
	UID uint32 `json:"uid" binding:"required"`
}

type GetUserResponse struct {
	UID        uint32 `json:"uid"`
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

	user, err := handler.userManager.GetUserByUID(ctx, req.UID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	resp = GetUserResponse{
		UID:        user.UID,
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
	UID uint32 `json:"uid"`
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
		UID: user.UID,
	}

	handler.Response(c, resp, nil)
}
