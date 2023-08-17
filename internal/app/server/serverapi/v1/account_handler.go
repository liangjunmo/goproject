package v1

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/server/servercode"
)

type AccountHandler struct {
	*BaseHandler
	accountUseCase *AccountUseCase
}

func NewAccountHandler(accountUseCase *AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
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
	l := log.WithContext(ctx)

	var (
		req  LoginRequest
		resp LoginResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	resp, err = handler.accountUseCase.Login(ctx, req)
	if err != nil {
		if errors.Is(err, servercode.InternalServerError) {
			l.Error(err)
		} else {
			l.Warn(err)
		}

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
	l := log.WithContext(ctx)

	var (
		req  CreateTokenRequest
		resp CreateTokenResponse
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		l.Error(err)
		handler.Response(c, nil, fmt.Errorf("%w: %v", servercode.InvalidRequest, err))
		return
	}

	resp, err = handler.accountUseCase.CreateToken(ctx, req)
	if err != nil {
		if errors.Is(err, servercode.InternalServerError) {
			l.Error(err)
		} else {
			l.Warn(err)
		}

		handler.Response(c, nil, err)
		return
	}

	handler.Response(c, resp, nil)
}

func (handler *AccountHandler) AuthMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	l := log.WithContext(ctx)

	claims, err := handler.accountUseCase.Auth(ctx, c.GetHeader("Authorization"))
	if err != nil {
		if errors.Is(err, servercode.InternalServerError) {
			l.Error(err)
		} else {
			l.Warn(err)
		}

		c.Abort()
		handler.Response(c, nil, err)
		return
	}

	c.Set("user_claims", claims)
	c.Next()
}
