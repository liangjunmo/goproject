package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint32         `gorm:"column:id;type:int unsigned;not null;auto_increment;primary key;" json:"-"`
	CreateTime time.Time      `gorm:"column:create_time;type:datetime;not null;autoCreateTime;" json:"-"`
	UpdateTime time.Time      `gorm:"column:update_time;type:datetime;not null;autoUpdateTime;" json:"-"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time;type:datetime;default:null;index:idx_delete_time;" json:"-"`
	UID        uint32         `gorm:"column:uid;type:int unsigned;not null;index:idx_uid;" json:"-"`
}

func (*User) TableName() string {
	return "user"
}

type UserDetail struct {
	UID        uint32
	Username   string
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime gorm.DeletedAt
}
