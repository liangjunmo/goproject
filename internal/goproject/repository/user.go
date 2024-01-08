package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/goproject/model"
	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type UserRepository interface {
	Begin() (UserRepository, error)
	Commit() error
	Rollback() error
	List(ctx context.Context, criteria UserCriteria) (pagination.Pagination, []model.User, error)
	Search(ctx context.Context, criteria UserCriteria) (map[uint32]model.User, error)
	Get(ctx context.Context, uid uint32) (user model.User, exist bool, err error)
	Create(ctx context.Context, user *model.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return newUserRepository(db)
}

var (
	userSortFields = map[string]string{
		"id":  "id",
		"uid": "uid",
	}

	userSortDirections = map[string]string{
		"asc":  "asc",
		"desc": "desc",
	}
)

type UserCriteria struct {
	pagination.Request
	Sorts map[string]string
	Uids  []uint32
}

type userRepository struct {
	db *gorm.DB
}

func newUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repository *userRepository) Begin() (UserRepository, error) {
	tx := repository.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return newUserRepository(tx), nil
}

func (repository *userRepository) Commit() error {
	return repository.db.Commit().Error
}

func (repository *userRepository) Rollback() error {
	return repository.db.Rollback().Error
}

func (repository *userRepository) List(ctx context.Context, criteria UserCriteria) (pagination.Pagination, []model.User, error) {
	db := repository.db.WithContext(ctx).Model(&model.User{})

	if len(criteria.Uids) != 0 {
		db = db.Where("id in (?)", criteria.Uids)
	}

	var (
		count int64
		users []model.User
	)

	err := db.Count(&count).Error
	if err != nil {
		return nil, nil, err
	}

	p := criteria.Request.Paginate(count)

	if count == 0 {
		return p, nil, nil
	}

	db = db.Limit(p.GetLimit()).Offset(p.GetOffset())

	db = repository.buildOrderBy(db, criteria.Sorts)

	err = db.Find(&users).Error
	if err != nil {
		return nil, nil, err
	}

	return p, users, nil
}

func (repository *userRepository) Search(ctx context.Context, criteria UserCriteria) (map[uint32]model.User, error) {
	db := repository.db.WithContext(ctx).Model(&model.User{})

	if len(criteria.Uids) != 0 {
		db = db.Where("id in (?)", criteria.Uids)
	}

	var users []model.User

	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]model.User)

	for _, u := range users {
		m[u.UID] = u
	}

	return m, nil
}

func (repository *userRepository) Get(ctx context.Context, uid uint32) (user model.User, exist bool, err error) {
	err = repository.db.WithContext(ctx).Take(&user, "uid = ?", uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, false, nil
		}

		return model.User{}, false, err
	}

	return user, true, nil
}

func (repository *userRepository) Create(ctx context.Context, user *model.User) error {
	return repository.db.WithContext(ctx).Create(user).Error
}

func (repository *userRepository) buildOrderBy(db *gorm.DB, sorts map[string]string) *gorm.DB {
	if len(sorts) == 0 {
		return db.Order("id desc")
	}

	for field, direction := range sorts {
		if realField, ok := userSortFields[field]; !ok {
			continue
		} else {
			field = realField
		}

		if _, ok := userSortDirections[direction]; !ok {
			direction = "desc"
		}

		db = db.Order(field + " " + direction)
	}

	return db
}
