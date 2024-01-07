package userserviceimpl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice"
)

type repository interface {
	Begin() (repository, error)
	Commit() error
	Rollback() error
	Search(ctx context.Context, criteria criteria) (map[uint32]userservice.User, error)
	Get(ctx context.Context, uid uint32) (user userservice.User, exist bool, err error)
	GetByUsername(ctx context.Context, username string) (user userservice.User, exist bool, err error)
	Create(ctx context.Context, user *userservice.User) error
}

var (
	sortFields = map[string]string{
		"id":  "id",
		"uid": "id",
	}

	sortDirections = map[string]string{
		"asc":  "asc",
		"desc": "desc",
	}
)

type criteria struct {
	sorts    map[string]string
	uids     []uint32
	username string
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

func (repository *defaultRepository) Search(ctx context.Context, criteria criteria) (map[uint32]userservice.User, error) {
	db := repository.db.WithContext(ctx).Model(&userservice.User{})

	if len(criteria.uids) != 0 {
		db = db.Where("id in (?)", criteria.uids)
	}

	if len(criteria.username) != 0 {
		db = db.Where("username like ?", "%"+criteria.username+"%")
	}

	var users []userservice.User

	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]userservice.User)

	for _, u := range users {
		m[u.UID] = u
	}

	return m, nil
}

func (repository *defaultRepository) Get(ctx context.Context, uid uint32) (user userservice.User, exist bool, err error) {
	err = repository.db.WithContext(ctx).Take(&user, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userservice.User{}, false, nil
		}

		return userservice.User{}, false, err
	}

	return user, true, nil
}

func (repository *defaultRepository) GetByUsername(ctx context.Context, username string) (user userservice.User, exist bool, err error) {
	err = repository.db.WithContext(ctx).Take(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userservice.User{}, false, nil
		}

		return userservice.User{}, false, err
	}

	return user, true, nil
}

func (repository *defaultRepository) Create(ctx context.Context, user *userservice.User) error {
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
