package dbdata

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/types"
)

func GetUserByUID(ctx context.Context, db *gorm.DB, uid uint32) (types.User, bool, error) {
	var user types.User

	err := db.WithContext(ctx).Model(&types.User{}).Where("id = ?", uid).Limit(1).Scan(&user).Error
	if err != nil {
		return types.User{}, false, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if user.UID == 0 {
		return types.User{}, false, nil
	}

	return user, true, nil
}

func GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (types.User, bool, error) {
	var user types.User

	err := db.WithContext(ctx).Model(&types.User{}).Where("username = ?", username).Limit(1).Scan(&user).Error
	if err != nil {
		return types.User{}, false, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if user.UID == 0 {
		return types.User{}, false, nil
	}

	return user, true, nil
}

func CreateUser(ctx context.Context, db *gorm.DB, user *types.User) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}

func UpdateUserByUID(ctx context.Context, db *gorm.DB, uid uint32, user types.User) error {
	err := db.WithContext(ctx).Model(&types.User{}).Where("id = ?", uid).Limit(1).Updates(&user).Error
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}