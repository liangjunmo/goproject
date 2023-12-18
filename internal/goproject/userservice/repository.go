package userservice

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/pkg/pagination"
)

type repository interface {
	Begin() (repository, error)
	Commit() error
	Rollback() error
	List(ctx context.Context, criteria criteria) (pagination.Pagination, []User, error)
	Search(ctx context.Context, criteria criteria) (map[uint32]User, error)
	Get(ctx context.Context, uid uint32) (user User, exist bool, err error)
	Create(ctx context.Context, user *User) error
}

var (
	sortFields = map[string]string{
		"id":  "id",
		"uid": "uid",
	}

	sortDirections = map[string]string{
		"asc":  "asc",
		"desc": "desc",
	}
)

type criteria struct {
	pagination.Request
	sorts map[string]string
	uids  []uint32
}

type defaultRepository struct {
	db *gorm.DB
}

func newDefaultRepository(db *gorm.DB) *defaultRepository {
	return &defaultRepository{
		db: db,
	}
}

func (repository *defaultRepository) Begin() (repository, error) {
	tx := repository.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return newDefaultRepository(tx), nil
}

func (repository *defaultRepository) Commit() error {
	return repository.db.Commit().Error
}

func (repository *defaultRepository) Rollback() error {
	return repository.db.Rollback().Error
}

func (repository *defaultRepository) List(ctx context.Context, criteria criteria) (pagination.Pagination, []User, error) {
	db := repository.db.WithContext(ctx).Model(&User{})

	if len(criteria.uids) != 0 {
		db = db.Where("id in (?)", criteria.uids)
	}

	var (
		count int64
		users []User
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

	db = repository.buildOrderBy(db, criteria.sorts)

	err = db.Find(&users).Error
	if err != nil {
		return nil, nil, err
	}

	return p, users, nil
}

func (repository *defaultRepository) Search(ctx context.Context, criteria criteria) (map[uint32]User, error) {
	db := repository.db.WithContext(ctx).Model(&User{})

	if len(criteria.uids) != 0 {
		db = db.Where("id in (?)", criteria.uids)
	}

	var users []User

	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]User)

	for _, u := range users {
		m[u.UID] = u
	}

	return m, nil
}

func (repository *defaultRepository) Get(ctx context.Context, uid uint32) (user User, exist bool, err error) {
	err = repository.db.WithContext(ctx).Take(&user, "uid = ?", uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, false, nil
		}

		return User{}, false, err
	}

	return user, true, nil
}

func (repository *defaultRepository) Create(ctx context.Context, user *User) error {
	return repository.db.WithContext(ctx).Create(user).Error
}

func (repository *defaultRepository) buildOrderBy(db *gorm.DB, sorts map[string]string) *gorm.DB {
	if len(sorts) == 0 {
		return db.Order("id desc")
	}

	for field, direction := range sorts {
		if realField, ok := sortFields[field]; !ok {
			continue
		} else {
			field = realField
		}

		if _, ok := sortDirections[direction]; !ok {
			direction = "desc"
		}

		db = db.Order(field + " " + direction)
	}

	return db
}
