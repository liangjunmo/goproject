package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type ListService interface {
	ListUser(ctx context.Context, cmd ListUserParams) (pagination.Pagination, []User, error)
}

type listService struct {
	db *gorm.DB
}

func NewListService(db *gorm.DB) ListService {
	return &listService{
		db: db,
	}
}

func (service *listService) ListUser(ctx context.Context, cmd ListUserParams) (pagination.Pagination, []User, error) {
	db := service.db.WithContext(ctx).Model(&User{})

	var count int64

	err := db.Count(&count).Error
	if err != nil {
		return pagination.Pagination{}, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	p := cmd.PaginationRequest.Paginate(count)

	if count == 0 {
		return p, nil, nil
	}

	var users []User

	err = db.Offset(p.Offset).Limit(p.Limit).Order("id desc").Find(&users).Error
	if err != nil {
		return pagination.Pagination{}, nil, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return p, users, nil
}
