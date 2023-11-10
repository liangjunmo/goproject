package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/liangjunmo/gocode"
	"github.com/liangjunmo/gotraceutil"
	"github.com/liangjunmo/logrushook"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/liangjunmo/goproject/internal/app/server/config"
	"github.com/liangjunmo/goproject/internal/codes"
	"github.com/liangjunmo/goproject/internal/pkg/timeutil"
	"github.com/liangjunmo/goproject/internal/version"
)

var (
	configFile string
)

const (
	envKeyServerConfigFile = "GOPROJECT_SERVER_CONFIG_FILE"
)

func init() {
	rootCmd.SetVersionTemplate("goproject-server:" + version.Describe())

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")

	cobra.OnInitialize(func() {
		if configFile == "" {
			configFile = os.Getenv(envKeyServerConfigFile)
		}

		loadConfig(configFile)
		initTrace()
		initLog()
	})
}

var rootCmd = &cobra.Command{
	Use:     "goproject-server",
	Short:   "goproject-server cli tool",
	Long:    "goproject-server cli tool",
	Version: version.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig(configFile string) {
	if configFile == "" {
		log.Fatal("config file is required")
	}

	log.Printf("use config file: %s", configFile)

	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config.Config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v", config.Config)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config.ProjectDir = dir

	log.Printf("project dir: %s", config.ProjectDir)
}

func initTrace() {
	gotraceutil.SetTraceIDKey(config.TraceIDKey)
	gotraceutil.SetTraceIDGenerator(gotraceutil.DefaultTraceIDGenerator)
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: timeutil.LayoutTime,
	})

	logrus.AddHook(gotraceutil.NewLogrusHook())

	logrus.AddHook(
		logrushook.NewReportCallerLogrusHook(
			[]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel},
			"file",
			func(path string, line int) string {
				return fmt.Sprintf("%s:%d", strings.Replace(path, config.ProjectDir+"/", "", -1), line)
			},
		),
	)

	logrus.AddHook(
		logrushook.NewTransformErrorLevelLogrusHook(
			logrus.WarnLevel,
			[]gocode.Code{codes.InternalServerError},
			true,
		),
	)
}

func connectDB(debug bool) *gorm.DB {
	level := gormlogger.Warn

	if debug {
		level = gormlogger.Info
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.DB.User,
		config.Config.DB.Password,
		config.Config.DB.Addr,
		config.Config.DB.Database,
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

	return db
}

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func newRedisSync(redisClient *redis.Client) *redsync.Redsync {
	return redsync.New(goredis.NewPool(redisClient))
}
