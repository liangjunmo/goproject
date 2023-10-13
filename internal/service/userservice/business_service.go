package userservice

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/dbdata"
	"github.com/liangjunmo/goproject/internal/redismutex"
	"github.com/liangjunmo/goproject/internal/types"
)

type BusinessService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (types.User, error)
}

type businessService struct {
	db        *gorm.DB
	redisSync *redsync.Redsync
}

func newBusinessService(db *gorm.DB, redisSync *redsync.Redsync) BusinessService {
	return &businessService{
		db:        db,
		redisSync: redisSync,
	}
}

func (service *businessService) CreateUser(ctx context.Context, req CreateUserRequest) (types.User, error) {
	mutex := redismutex.NewCreateUserMutex(service.redisSync, req.UID)

	err := mutex.Lock()
	if err != nil {
		return types.User{}, codes.Timeout
	}
	defer mutex.Unlock()

	user, ok, err := dbdata.GetUserByUID(ctx, service.db, req.UID)
	if err != nil {
		return types.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if ok {
		return types.User{}, codes.UserAlreadyExists
	}

	user = types.User{
		UID: req.UID,
	}

	err = dbdata.CreateUser(ctx, service.db, &user)
	if err != nil {
		return types.User{}, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return user, nil
}
