package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/goproject/manager"
	"github.com/liangjunmo/goproject/internal/goproject/model"
	"github.com/liangjunmo/goproject/internal/goproject/service"
	"github.com/liangjunmo/goproject/internal/pkg/hashutil"
)

type AccountUseCaseConfig struct {
	JWTKey string
}

type AccountUseCase interface {
	Login(ctx context.Context, cmd LoginCommand) (ticket string, failedCount uint32, err error)
	CreateToken(ctx context.Context, cmd CreateTokenCommand) (token string, err error)
	Authorize(ctx context.Context, cmd AuthorizeCommand) (*model.UserJwtClaims, error)
}

func NewAccountUseCase(config AccountUseCaseConfig, redisManager manager.RedisManager, userService service.UserService) AccountUseCase {
	return newAccountUseCase(config, redisManager, userService)
}

type accountUseCase struct {
	log          *logrus.Entry
	config       AccountUseCaseConfig
	redisManager manager.RedisManager
	userService  service.UserService
}

func newAccountUseCase(config AccountUseCaseConfig, redisManager manager.RedisManager, userService service.UserService) AccountUseCase {
	return &accountUseCase{
		log:          logrus.WithField("tag", "goproject.accountservice.service"),
		config:       config,
		redisManager: redisManager,
		userService:  userService,
	}
}

type LoginCommand struct {
	Username string
	Password string
}

func (usecase *accountUseCase) Login(ctx context.Context, cmd LoginCommand) (ticket string, failedCount uint32, err error) {
	count, err := usecase.redisManager.GetLoginFailedCount(ctx, cmd.Username)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if count >= 5 {
		return "", count, codes.LoginFailedReachLimit
	}

	err = usecase.userService.ValidatePassword(ctx, service.ValidatePasswordCommand{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		e := usecase.redisManager.SetLoginFailedCount(ctx, cmd.Username)
		if e != nil {
			return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, e)
		}

		return "", count + 1, err
	}

	err = usecase.redisManager.DelLoginFailedCount(ctx, cmd.Username)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	user, err := usecase.userService.GetByUsername(ctx, service.GetByUsernameCommand{Username: cmd.Username})
	if err != nil {
		return "", 0, err
	}

	ticket = usecase.generateLoginTicket(user.UID)

	err = usecase.redisManager.SetLoginTicket(ctx, ticket, user.UID, time.Minute)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return ticket, 0, nil
}

type CreateTokenCommand struct {
	Ticket string
}

func (usecase *accountUseCase) CreateToken(ctx context.Context, req CreateTokenCommand) (token string, err error) {
	uid, ok, err := usecase.redisManager.GetLoginTicket(ctx, req.Ticket)
	if err != nil {
		return "", fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return "", codes.AuthorizeFailedInvalidTicket
	}

	user, err := usecase.userService.Get(ctx, service.GetCommand{UID: uid})
	if err != nil {
		return "", err
	}

	claims := model.UserJwtClaims{
		UID: user.UID,
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	token, err = usecase.generateJwtToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

type AuthorizeCommand struct {
	Token string
}

func (usecase *accountUseCase) Authorize(ctx context.Context, cmd AuthorizeCommand) (*model.UserJwtClaims, error) {
	if cmd.Token == "" {
		return nil, codes.AuthorizeFailedInvalidToken
	}

	jwtClaims, err := usecase.parseJwtToken(cmd.Token, &model.UserJwtClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtClaims.(*model.UserJwtClaims)
	if !ok {
		return nil, fmt.Errorf("%w: jwt claims can not trans to *UserJwtClaims", codes.AuthorizeFailed)
	}

	_, err = usecase.userService.Get(ctx, service.GetCommand{UID: claims.UID})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (usecase *accountUseCase) generateLoginTicket(uid uint32) string {
	t := fmt.Sprintf("%d%d", uid, time.Now().Unix())

	b := hashutil.SHA1StringToByte(t)

	return base64.URLEncoding.EncodeToString(b)
}

func (usecase *accountUseCase) generateJwtToken(claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(usecase.config.JWTKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return token, nil
}

func (usecase *accountUseCase) parseJwtToken(token string, claims jwt.Claims) (jwt.Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(usecase.config.JWTKey), nil
	})
	if err != nil {
		return jwt.Claims(nil), fmt.Errorf("%w: %v", codes.AuthorizeFailedInvalidToken, err)
	}

	if jwtToken != nil && jwtToken.Valid {
		return jwtToken.Claims, nil
	}

	return jwt.Claims(nil), fmt.Errorf("%w: invalid jwt token", codes.AuthorizeFailedInvalidToken)
}
