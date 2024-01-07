package usercenter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/liangjunmo/gocode"
	"github.com/liangjunmo/gotraceutil"
	"github.com/liangjunmo/logrushook/reportcallerhook"
	"github.com/liangjunmo/logrushook/transformerrorhook"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice"
)

type DBConfig struct {
	Addr     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Addr     string
	Password string
}

func initTracing(tracingKeys []string) {
	gotraceutil.SetTracingKeys(tracingKeys)
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.AddHook(gotraceutil.NewLogrusHook())

	{
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		handler := func(path string, line int) string {
			return fmt.Sprintf("%s:%d", strings.Replace(path, dir+"/", "", -1), line)
		}

		hook := reportcallerhook.New([]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel})

		hook.SetKey("file")
		hook.SetLocationHandler(handler)

		logrus.AddHook(hook)
	}

	{
		hook := transformerrorhook.New(logrus.WarnLevel)

		hook.ExcludeCodes([]gocode.Code{codes.InternalServerError})
		hook.DeleteErrorKey()

		logrus.AddHook(hook)
	}
}

func initDB(config DBConfig, debug bool) *gorm.DB {
	level := gormlogger.Warn

	if debug {
		level = gormlogger.Info
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Addr,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gotraceutil.NewGORMLogger(
			gormlogger.Config{
				SlowThreshold:             time.Millisecond * 100,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  level,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&userservice.User{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func initRedis(config RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	return redisClient
}
