package datautil

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/app/server/codes"
	"github.com/liangjunmo/goproject/internal/app/server/types"
)

func DbGetUserByUid(ctx context.Context, db *gorm.DB, uid uint32) (types.User, bool, error) {
	var user types.User

	err := db.WithContext(ctx).Model(&types.User{}).Where("id = ?", uid).Limit(1).Scan(&user).Error
	if err != nil {
		return types.User{}, false, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if user.Uid == 0 {
		return types.User{}, false, nil
	}

	return user, true, nil
}

func DbGetUserByUsername(ctx context.Context, db *gorm.DB, username string) (types.User, bool, error) {
	var user types.User

	err := db.WithContext(ctx).Model(&types.User{}).Where("username = ?", username).Limit(1).Scan(&user).Error
	if err != nil {
		return types.User{}, false, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if user.Uid == 0 {
		return types.User{}, false, nil
	}

	return user, true, nil
}

func DbCreateUser(ctx context.Context, db *gorm.DB, user *types.User) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}

func DbUpdateUserByUid(ctx context.Context, db *gorm.DB, uid uint32, user types.User) error {
	err := db.WithContext(ctx).Model(&types.User{}).Where("id = ?", uid).Limit(1).Updates(&user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}
