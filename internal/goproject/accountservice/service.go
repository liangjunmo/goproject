package accountservice

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/goproject/userservice"
	"github.com/liangjunmo/goproject/internal/pkg/hashutil"
)

type Service interface {
	Login(ctx context.Context, cmd LoginCommand) (ticket string, failedCount uint32, err error)
	CreateToken(ctx context.Context, cmd CreateTokenCommand) (token string, err error)
	Authorize(ctx context.Context, cmd AuthorizeCommand) (*UserJwtClaims, error)
}

func ProvideService(config Config, redisClient *redis.Client, userService userservice.Service) Service {
	return newDefaultService(
		config,
		newDefaultRedisManager(redisClient),
		userService,
	)
}

type defaultService struct {
	log          *logrus.Entry
	config       Config
	redisManager redisManager
	userService  userservice.Service
}

func newDefaultService(config Config, redisManager redisManager, userService userservice.Service) Service {
	return &defaultService{
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

type LoginResult struct {
	Ticket      string
	FailedCount uint32
}

func (service *defaultService) Login(ctx context.Context, cmd LoginCommand) (ticket string, failedCount uint32, err error) {
	count, err := service.redisManager.GetLoginFailedCount(ctx, cmd.Username)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if count >= 5 {
		return "", count, codes.LoginFailedReachLimit
	}

	err = service.userService.ValidatePassword(ctx, userservice.ValidatePasswordCommand{
		Username: cmd.Username,
		Password: cmd.Password,
	})
	if err != nil {
		e := service.redisManager.SetLoginFailedCount(ctx, cmd.Username)
		if e != nil {
			return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, e)
		}

		return "", count + 1, err
	}

	err = service.redisManager.DelLoginFailedCount(ctx, cmd.Username)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	user, err := service.userService.GetByUsername(ctx, userservice.GetByUsernameCommand{Username: cmd.Username})
	if err != nil {
		return "", 0, err
	}

	ticket = service.generateLoginTicket(user.UID)

	err = service.redisManager.SetLoginTicket(ctx, ticket, user.UID, time.Minute)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return ticket, 0, nil
}

type CreateTokenCommand struct {
	Ticket string
}

func (service *defaultService) CreateToken(ctx context.Context, req CreateTokenCommand) (token string, err error) {
	uid, ok, err := service.redisManager.GetLoginTicket(ctx, req.Ticket)
	if err != nil {
		return "", fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if !ok {
		return "", codes.AuthorizeFailedInvalidTicket
	}

	user, err := service.userService.Get(ctx, userservice.GetCommand{UID: uid})
	if err != nil {
		return "", err
	}

	claims := UserJwtClaims{
		UID: user.UID,
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	token, err = service.generateJwtToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

type AuthorizeCommand struct {
	Token string
}

func (service *defaultService) Authorize(ctx context.Context, cmd AuthorizeCommand) (*UserJwtClaims, error) {
	if cmd.Token == "" {
		return nil, codes.AuthorizeFailedInvalidToken
	}

	jwtClaims, err := service.parseJwtToken(cmd.Token, &UserJwtClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtClaims.(*UserJwtClaims)
	if !ok {
		return nil, fmt.Errorf("%w: jwt claims can not trans to *UserJwtClaims", codes.AuthorizeFailed)
	}

	_, err = service.userService.Get(ctx, userservice.GetCommand{UID: claims.UID})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (service *defaultService) generateLoginTicket(uid uint32) string {
	t := fmt.Sprintf("%d%d", uid, time.Now().Unix())

	b := hashutil.SHA1StringToByte(t)

	return base64.URLEncoding.EncodeToString(b)
}

func (service *defaultService) generateJwtToken(claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(service.config.JWTKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return token, nil
}

func (service *defaultService) parseJwtToken(token string, claims jwt.Claims) (jwt.Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.config.JWTKey), nil
	})
	if err != nil {
		return jwt.Claims(nil), fmt.Errorf("%w: %v", codes.AuthorizeFailedInvalidToken, err)
	}

	if jwtToken != nil && jwtToken.Valid {
		return jwtToken.Claims, nil
	}

	return jwt.Claims(nil), fmt.Errorf("%w: invalid jwt token", codes.AuthorizeFailedInvalidToken)
}
