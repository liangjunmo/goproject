package dbdata

import (
	"context"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/types"
)

func GetUserCenterUserByUID(ctx context.Context, db *gorm.DB, uid uint32) (types.UserCenterUser, bool, error) {
	var user types.UserCenterUser

	err := db.WithContext(ctx).Model(&types.UserCenterUser{}).Where("id = ?", uid).Limit(1).Scan(&user).Error
	if err != nil {
		return types.UserCenterUser{}, false, err
	}

	if user.UID == 0 {
		return types.UserCenterUser{}, false, nil
	}

	return user, true, nil
}

func GetUserCenterUserByUsername(ctx context.Context, db *gorm.DB, username string) (types.UserCenterUser, bool, error) {
	var user types.UserCenterUser

	err := db.WithContext(ctx).Model(&types.UserCenterUser{}).Where("username = ?", username).Limit(1).Scan(&user).Error
	if err != nil {
		return types.UserCenterUser{}, false, err
	}

	if user.UID == 0 {
		return types.UserCenterUser{}, false, nil
	}

	return user, true, nil
}

func CreateUserCenterUser(ctx context.Context, db *gorm.DB, user *types.UserCenterUser) error {
	err := db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserCenterUserByUID(ctx context.Context, db *gorm.DB, uid uint32, user types.UserCenterUser) error {
	err := db.WithContext(ctx).Model(&types.UserCenterUser{}).Where("id = ?", uid).Limit(1).Updates(&user).Error
	if err != nil {
		return err
	}

	return nil
}
