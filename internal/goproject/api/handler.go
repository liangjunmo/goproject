package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gocode"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/goproject/accountservice"
	"github.com/liangjunmo/goproject/internal/goproject/userservice"
)

var (
	ginCtxUserKey = "user_claims"
)

type handler struct {
	config         Config
	accountService accountservice.Service
	userService    userservice.Service
}

func newHandler(config Config, accountService accountservice.Service, userService userservice.Service) *handler {
	return &handler{
		config:         config,
		accountService: accountService,
		userService:    userService,
	}
}

func (handler *handler) responseDefault(c *gin.Context, data interface{}, err error) {
	handler.response(c, 200, data, err)
}

func (handler *handler) response(c *gin.Context, httpStatusCode int, data interface{}, err error) {
	c.JSON(httpStatusCode, handler.buildResponseBody(c, data, err))
}

func (handler *handler) buildResponseBody(c *gin.Context, data interface{}, err error) gin.H {
	if data == nil {
		data = map[string]interface{}{}
	}

	code := gocode.Parse(err)
	if errors.Is(code, gocode.SuccessCode) {
		code = codes.OK
	} else if errors.Is(code, gocode.DefaultCode) {
		code = codes.Unknown
	}

	body := gin.H{
		"data": data,
		"code": code,
		"msg":  codes.Translate(code, codes.Language(c.GetHeader("Accept-Language"))),
	}

	if handler.config.Debug {
		body["error"] = nil
		if err != nil {
			body["error"] = err.Error()
		}

		body["request_id"] = c.Request.Context().Value(handler.config.TracingIDKey)
	}

	return body
}

func (handler *handler) getUserClaims(c *gin.Context) *accountservice.UserJwtClaims {
	user, _ := c.Get(ginCtxUserKey)
	return user.(*accountservice.UserJwtClaims)
}

func (handler *handler) Ping(c *gin.Context) {
	c.String(200, "pong")
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Ticket      string `json:"ticket"`
	FailedCount uint32 `json:"failed_count"`
}

func (handler *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  LoginRequest
		resp LoginResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	ticket, failedCount, err := handler.accountService.Login(ctx, accountservice.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		resp.FailedCount = failedCount
		handler.responseDefault(c, resp, err)
		return
	}

	resp = LoginResponse{
		Ticket:      ticket,
		FailedCount: failedCount,
	}

	handler.responseDefault(c, resp, nil)
}

type CreateTokenRequest struct {
	Ticket string `json:"ticket" binding:"required"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (handler *handler) CreateToken(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  CreateTokenRequest
		resp CreateTokenResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	token, err := handler.accountService.CreateToken(ctx, accountservice.CreateTokenCommand{Ticket: req.Ticket})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, err)
		return
	}

	resp = CreateTokenResponse{Token: token}

	handler.responseDefault(c, resp, nil)
}

func (handler *handler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()

	claims, err := handler.accountService.Authorize(ctx, accountservice.AuthorizeCommand{Token: c.GetHeader("Authorization")})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		c.Abort()
		handler.responseDefault(c, nil, err)
		return
	}

	c.Set(ginCtxUserKey, claims)
	c.Next()
}

type ListUserRequest struct {
	PaginationRequest
}

type ListUserResponse struct {
	Pagination Pagination     `json:"pagination"`
	List       []ListUserData `json:"list"`
}

type ListUserData struct {
	UID        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *handler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  ListUserRequest
		resp ListUserResponse
	)

	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	p, users, err := handler.userService.List(ctx, userservice.ListCommand{Request: req.PaginationRequest})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, err)
		return
	}

	resp = ListUserResponse{
		Pagination: Pagination{
			Page:            p.GetPage(),
			CapacityPerPage: p.GetCapacityPerPage(),
			TotalPages:      p.GetTotalPages(),
			TotalRecords:    p.GetTotalRecords(),
		},
		List: make([]ListUserData, 0, len(users)),
	}

	if len(users) == 0 {
		handler.responseDefault(c, resp, nil)
		return
	}

	for _, user := range users {
		resp.List = append(resp.List, ListUserData{
			UID:        user.UID,
			Username:   user.UserCenterUser.Username,
			CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	handler.responseDefault(c, resp, nil)
}

type SearchUserRequest struct {
	Uids     []string `form:"uids[]"`
	Username string   `form:"username"`
}

type SearchUserResponse []SearchUserData

type SearchUserData struct {
	UID        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *handler) SearchUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  SearchUserRequest
		resp SearchUserResponse
	)

	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	uids := make([]uint32, 0, len(req.Uids))

	for _, uid := range req.Uids {
		uids = append(uids, cast.ToUint32(uid))
	}

	users, err := handler.userService.Search(ctx, userservice.SearchCommand{
		Uids:     uids,
		Username: req.Username,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, err)
		return
	}

	resp = make([]SearchUserData, 0, len(users))

	if len(users) == 0 {
		handler.responseDefault(c, resp, nil)
		return
	}

	for _, user := range users {
		resp = append(resp, SearchUserData{
			UID:        user.UID,
			Username:   user.UserCenterUser.Username,
			CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	handler.responseDefault(c, resp, nil)
}

type GetUserRequest struct {
	UID uint32 `uri:"uid" binding:"required"`
}

type GetUserResponse struct {
	UID        uint32 `json:"uid"`
	Username   string `json:"username"`
	CreateTime string `json:"create_time"`
}

func (handler *handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  GetUserRequest
		resp GetUserResponse
	)

	err := c.ShouldBindUri(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	user, err := handler.userService.Get(ctx, userservice.GetCommand{UID: req.UID})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, err)
		return
	}

	resp = GetUserResponse{
		UID:        user.UID,
		Username:   user.UserCenterUser.Username,
		CreateTime: user.CreateTime.Format("2006-01-02 15:04:05"),
	}

	handler.responseDefault(c, resp, nil)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	UID uint32 `json:"uid"`
}

func (handler *handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  CreateUserRequest
		resp CreateUserResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	uid, err := handler.userService.Create(ctx, userservice.CreateCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.responseDefault(c, nil, err)
		return
	}

	resp = CreateUserResponse{UID: uid}

	handler.responseDefault(c, resp, nil)
}
