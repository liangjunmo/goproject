package dbutil

import (
	"time"

	"github.com/liangjunmo/gotraceutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(dsn string, config *gorm.Config) (*gorm.DB, error) {
	if config == nil {
		config = &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: gotraceutil.NewGormLogger(
				gormlogger.Config{
					SlowThreshold:             time.Millisecond * 100,
					Colorful:                  false,
					IgnoreRecordNotFoundError: true,
					ParameterizedQueries:      false,
					LogLevel:                  gormlogger.Info,
				},
			),
		}
	}

	return gorm.Open(mysql.Open(dsn), config)
}
