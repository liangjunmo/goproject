package types

import (
	"time"

	"gorm.io/gorm"
)

type UserCenterUser struct {
	UID        uint32         `gorm:"column:id;type:int unsigned;not null;auto_increment;primary key;" json:"-"`
	CreateTime time.Time      `gorm:"column:create_time;type:datetime;not null;autoCreateTime;" json:"-"`
	UpdateTime time.Time      `gorm:"column:update_time;type:datetime;not null;autoUpdateTime;" json:"-"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;type:datetime;default:null;index:idx_delete_time;" json:"-"`
	Username   string         `gorm:"column:username;type:varchar(32);not null;index:idx_username,unique;" json:"-"`
	Password   string         `gorm:"column:password;type:varchar(100);not null;" json:"-"`
}

func (*UserCenterUser) TableName() string {
	return "usercenter_user"
}
