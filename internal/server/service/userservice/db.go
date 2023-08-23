package userservice

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/server/servercode"
)

func DbGetUserByUid(ctx context.Context, db *gorm.DB, uid uint32) (User, bool, error) {
	var user User

	err := db.WithContext(ctx).Model(&User{}).Where("id = ?", uid).Limit(1).Scan(&user).Error
	if err != nil {
		return User{}, false, fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	if user.Id == 0 {
		return User{}, false, nil
	}

	return user, true, nil
}

func DbGetUserByUsername(ctx context.Context, db *gorm.DB, username string) (User, bool, error) {
	var user User

	err := db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Limit(1).Scan(&user).Error
	if err != nil {
		return User{}, false, fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	if user.Id == 0 {
		return User{}, false, nil
	}

	return user, true, nil
}

func DbCreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	return nil
}

func DbUpdateUserByUid(ctx context.Context, db *gorm.DB, uid uint32, user User) error {
	err := db.WithContext(ctx).Model(&User{}).Where("id = ?", uid).Limit(1).Updates(&user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", servercode.InternalServerError, err)
	}

	return nil
}
