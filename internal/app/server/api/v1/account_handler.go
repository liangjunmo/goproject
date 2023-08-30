package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
)

type AccountHandler struct {
	*BaseHandler
	accountComponent *AccountComponent
}

func NewAccountHandler(accountUseCase *AccountComponent) *AccountHandler {
	return &AccountHandler{
		accountComponent: accountUseCase,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Ticket      string `json:"ticket"`
	FailedCount uint32 `json:"failed_count"`
}

func (handler *AccountHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  LoginRequest
		resp LoginResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	resp, err = handler.accountComponent.Login(ctx, req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, resp, err)
		return
	}

	handler.Response(c, resp, nil)
}

type CreateTokenRequest struct {
	Ticket string `json:"ticket"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (handler *AccountHandler) CreateToken(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		req  CreateTokenRequest
		resp CreateTokenResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", codes.InvalidRequest, err))
		return
	}

	resp, err = handler.accountComponent.CreateToken(ctx, req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		handler.Response(c, nil, err)
		return
	}

	handler.Response(c, resp, nil)
}

func (handler *AccountHandler) AuthMiddleware(c *gin.Context) {
	ctx := c.Request.Context()

	claims, err := handler.accountComponent.Auth(ctx, c.GetHeader("Authorization"))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(err)
		c.Abort()
		handler.Response(c, nil, err)
		return
	}

	c.Set("user_claims", claims)
	c.Next()
}
