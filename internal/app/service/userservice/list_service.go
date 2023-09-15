package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/codes"
	"github.com/liangjunmo/goproject/internal/app/types"
	"github.com/liangjunmo/goproject/internal/pkg/pageutil"
)

type ListService interface {
	ListUser(ctx context.Context, req ListUserRequest) (pageutil.Pagination, []types.User, error)
}

type listService struct {
	db *gorm.DB
}

func newListService(db *gorm.DB) ListService {
	return &listService{
		db: db,
	}
}

func (service *listService) ListUser(ctx context.Context, req ListUserRequest) (pageutil.Pagination, []types.User, error) {
	db := service.db.WithContext(ctx).Model(&types.User{})

	var count int64

	err := db.Count(&count).Error
	if err != nil {
		return pageutil.Pagination{}, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	p := req.PaginationRequest.Paginate(count)

	if count == 0 {
		return p, nil, nil
	}

	var users []types.User

	err = db.Offset(p.Offset).Limit(p.Limit).Order("id desc").Find(&users).Error
	if err != nil {
		return pageutil.Pagination{}, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return p, users, nil
}