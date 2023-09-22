package v1

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/pkg/hashutil"
	"github.com/liangjunmo/goproject/internal/redisdata"
	"github.com/liangjunmo/goproject/internal/server/config"
	"github.com/liangjunmo/goproject/internal/service/userservice"
)

type AccountComponent struct {
	redisClient *redis.Client
	userService userservice.Service
}

func NewAccountComponent(redisClient *redis.Client, userService userservice.Service) *AccountComponent {
	return &AccountComponent{
		redisClient: redisClient,
		userService: userService,
	}
}

func (component *AccountComponent) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	count, err := redisdata.GetLoginFailedCount(ctx, component.redisClient, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}

	if count >= 5 {
		return LoginResponse{
			FailedCount: count,
		}, codes.LoginFailedReachLimit
	}

	err = component.userService.ValidatePassword(ctx, userservice.ValidatePasswordRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		e := redisdata.SetLoginFailedCount(ctx, component.redisClient, req.Username)
		if e != nil {
			return LoginResponse{}, e
		}

		return LoginResponse{
			FailedCount: count + 1,
		}, err
	}

	err = redisdata.DelLoginFailedCount(ctx, component.redisClient, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}

	user, err := component.userService.GetUser(ctx, userservice.GetUserRequest{
		Username: req.Username,
	})
	if err != nil {
		return LoginResponse{}, err
	}

	ticket := component.generateLoginTicket(user.UID)

	err = redisdata.SetLoginTicket(ctx, component.redisClient, ticket, user.UID, time.Minute)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Ticket: ticket,
	}, nil
}

func (component *AccountComponent) CreateToken(ctx context.Context, req CreateTokenRequest) (CreateTokenResponse, error) {
	uid, ok, err := redisdata.GetLoginTicket(ctx, component.redisClient, req.Ticket)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	if !ok {
		return CreateTokenResponse{}, codes.AuthorizeInvalidTicket
	}

	user, err := component.userService.GetUser(ctx, userservice.GetUserRequest{
		UID: uid,
	})
	if err != nil {
		return CreateTokenResponse{}, err
	}

	claims := UserJwtClaims{
		UID: user.UID,
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	token, err := component.generateJwtToken(claims)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	return CreateTokenResponse{
		Token: token,
	}, nil
}

func (component *AccountComponent) Auth(ctx context.Context, token string) (*UserJwtClaims, error) {
	if token == "" {
		return nil, codes.AuthorizeInvalidTicket
	}

	jwtClaims, err := component.parseJwtToken(token, &UserJwtClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtClaims.(*UserJwtClaims)
	if !ok {
		return nil, fmt.Errorf("%w: jwt claims can not trans to *UserJwtClaims", codes.AuthorizeFailed)
	}

	_, err = component.userService.GetUser(ctx, userservice.GetUserRequest{
		UID: claims.UID,
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (component *AccountComponent) generateLoginTicket(uid uint32) string {
	s := fmt.Sprintf("%d%d", uid, time.Now().Unix())
	b := hashutil.Sha1StringToByte(s)
	return base64.URLEncoding.EncodeToString(b)
}

func (component *AccountComponent) generateJwtToken(claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(config.Config.API.JwtKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return token, nil
}

func (component *AccountComponent) parseJwtToken(token string, claims jwt.Claims) (jwt.Claims, error) {
	var jwtToken *jwt.Token

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.API.JwtKey), nil
	})
	if err != nil {
		return jwt.Claims(nil), fmt.Errorf("%w: %v", codes.AuthorizeInvalidToken, err)
	}

	if jwtToken != nil && jwtToken.Valid {
		return jwtToken.Claims, nil
	}

	return jwt.Claims(nil), fmt.Errorf("%w: invalid jwt token", codes.AuthorizeInvalidToken)
}
