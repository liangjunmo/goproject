package v1

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"

	"github.com/liangjunmo/goproject/internal/app/server/servercode"
	"github.com/liangjunmo/goproject/internal/app/server/serverconfig"
	"github.com/liangjunmo/goproject/internal/app/server/service/userservice"
)

type AccountUseCase struct {
	redisClient    *redis.Client
	userHubService userservice.HubService
}

func NewAccountUseCase(redisClient *redis.Client, userHubService userservice.HubService) *AccountUseCase {
	return &AccountUseCase{
		redisClient:    redisClient,
		userHubService: userHubService,
	}
}

func (uc *AccountUseCase) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	count, err := RedisGetLoginFailedCount(ctx, uc.redisClient, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}

	if count >= 5 {
		return LoginResponse{
			FailedCount: count,
		}, servercode.LoginFailedReachLimit
	}

	err = uc.userHubService.ValidatePassword(ctx, userservice.ValidatePasswordCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		e := RedisSetLoginFailedCount(ctx, uc.redisClient, req.Username)
		if e != nil {
			return LoginResponse{}, e
		}

		return LoginResponse{
			FailedCount: count + 1,
		}, err
	}

	err = RedisDelLoginFailedCount(ctx, uc.redisClient, req.Username)
	if err != nil {
		return LoginResponse{}, err
	}

	user, err := uc.userHubService.GetUser(ctx, userservice.GetUserCommand{
		Username: req.Username,
	})
	if err != nil {
		return LoginResponse{}, err
	}

	ticket := uc.generateLoginTicket(user.Id)

	err = RedisSetLoginTicket(ctx, uc.redisClient, ticket, user.Id, time.Minute)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Ticket: ticket,
	}, nil
}

func (uc *AccountUseCase) CreateToken(ctx context.Context, req CreateTokenRequest) (CreateTokenResponse, error) {
	uid, ok, err := RedisGetLoginTicket(ctx, uc.redisClient, req.Ticket)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	if !ok {
		return CreateTokenResponse{}, servercode.AuthorizeInvalidTicket
	}

	user, err := uc.userHubService.GetUser(ctx, userservice.GetUserCommand{
		Uid: uid,
	})
	if err != nil {
		return CreateTokenResponse{}, err
	}

	claims := UserJwtClaims{
		Uid: user.Id,
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	token, err := uc.generateJwtToken(claims)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	return CreateTokenResponse{
		Token: token,
	}, nil
}

func (uc *AccountUseCase) Auth(ctx context.Context, token string) (*UserJwtClaims, error) {
	if token == "" {
		return nil, servercode.AuthorizeInvalidTicket
	}

	jwtClaims, err := uc.parseJwtToken(token, &UserJwtClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtClaims.(*UserJwtClaims)
	if !ok {
		return nil, fmt.Errorf("%w: jwt claims can not trans to *UserJwtClaims", servercode.AuthorizeFailed)
	}

	_, err = uc.userHubService.GetUser(ctx, userservice.GetUserCommand{
		Uid: claims.Uid,
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (uc *AccountUseCase) generateLoginTicket(uid uint32) string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%d%d", uid, time.Now().Unix())))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func (uc *AccountUseCase) generateJwtToken(claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(serverconfig.Config.Api.JwtKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	return token, nil
}

func (uc *AccountUseCase) parseJwtToken(token string, claims jwt.Claims) (jwt.Claims, error) {
	var jwtToken *jwt.Token

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(serverconfig.Config.Api.JwtKey), nil
	})
	if err != nil {
		return jwt.Claims(nil), fmt.Errorf("%w: %v", servercode.AuthorizeInvalidToken, err)
	}

	if jwtToken != nil && jwtToken.Valid {
		return jwtToken.Claims, nil
	}

	return jwt.Claims(nil), fmt.Errorf("%w: invalid jwt token", servercode.AuthorizeInvalidToken)
}
