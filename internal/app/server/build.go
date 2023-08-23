package server

import (
	"context"
	"fmt"
	golog "log"
	"os"
	"strings"
	"time"

	"github.com/liangjunmo/gocode"
	"github.com/liangjunmo/gotraceutil"
	"github.com/liangjunmo/logrushook"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/liangjunmo/goproject/internal/app/server/servercode"
	"github.com/liangjunmo/goproject/internal/app/server/serverconfig"
	"github.com/liangjunmo/goproject/internal/pkg/dbutil"
	"github.com/liangjunmo/goproject/internal/pkg/timeutil"
)

func BuildConfig(configFile string) error {
	if configFile == "" {
		return fmt.Errorf("serverconfig file is required")
	}

	golog.Printf("use serverconfig file: %s", configFile)

	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&serverconfig.Config)
	if err != nil {
		return err
	}

	golog.Printf("serverconfig: %+v", serverconfig.Config)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	serverconfig.ProjectDir = dir

	return nil
}

func BuildTrace() error {
	gotraceutil.SetTraceIdGenerator(gotraceutil.DefaultTraceIdGenerator)

	gotraceutil.SetTraceIdKey(serverconfig.TraceIdKey)

	return nil
}

func BuildLog() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: timeutil.LayoutTime,
	})

	logrus.AddHook(gotraceutil.NewLogrusHook())

	logrus.AddHook(
		logrushook.NewReportCallerLogrusHook(
			[]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel},
			func(path string) string {
				return strings.Replace(path, serverconfig.ProjectDir+"/", "", -1)
			},
		),
	)

	logrus.AddHook(
		logrushook.NewTransErrorToWarningLogrusHook(
			[]gocode.Code{servercode.InternalServerError},
			false,
		),
	)

	return nil
}

func BuildDb(debug bool) (*gorm.DB, error) {
	level := gormlogger.Warn

	if debug {
		level = gormlogger.Info
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		serverconfig.Config.Db.User,
		serverconfig.Config.Db.Password,
		serverconfig.Config.Db.Addr,
		serverconfig.Config.Db.Database,
	)

	db, err := dbutil.Connect(dsn, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: gotraceutil.NewGormLogger(
			gormlogger.Config{
				SlowThreshold:             time.Millisecond * 100,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				LogLevel:                  level,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func BuildRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     serverconfig.Config.Redis.Addr,
		Password: serverconfig.Config.Redis.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
