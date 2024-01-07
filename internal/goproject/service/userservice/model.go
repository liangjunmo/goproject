package userservice

import (
	"time"

	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/api/usercenterproto"
)

type User struct {
	ID         uint32         `gorm:"column:id;type:int unsigned;not null;auto_increment;primary key;" json:"-"`
	CreateTime time.Time      `gorm:"column:create_time;type:datetime;not null;autoCreateTime;" json:"-"`
	UpdateTime time.Time      `gorm:"column:update_time;type:datetime;not null;autoUpdateTime;" json:"-"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;type:datetime;default:null;index:idx_delete_time;" json:"-"`
	UID        uint32         `gorm:"column:uid;type:int unsigned;not null;index:idx_uid;" json:"-"`

	UserCenterUser *usercenterproto.User `gorm:"-" json:"-"`
}

func (*User) TableName() string {
	return "user"
}
