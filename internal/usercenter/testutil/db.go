package testutil

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	// GOPROJECT_USERCENTER_TEST_DB="user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("GOPROJECT_USERCENTER_TEST_DB")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gormlogger.Info,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
