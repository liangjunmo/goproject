package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/liangjunmo/goproject/internal/goproject"
)

func init() {
	rootCmd.SetVersionTemplate("goproject-server:" + goproject.VersionInfo())

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")

	cobra.OnInitialize(func() {
		loadConfig()
	})
}

var rootCmd = &cobra.Command{
	Use:     "goproject",
	Version: goproject.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

type ConfigTemplate struct {
	Environment             string `mapstructure:"environment"`
	Debug                   bool   `mapstructure:"debug"`
	JWTKey                  string `mapstructure:"jwtKey"`
	APIAServerAddr          string `mapstructure:"apiServerAddr"`
	UserCenterRPCServerAddr string `mapstructure:"userCenterRPCServerAddr"`
	DB                      struct {
		Addr     string `mapstructure:"addr"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	} `mapstructure:"db"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
}

var (
	configFile string
	config     ConfigTemplate
)

func loadConfig() {
	if configFile == "" {
		configFile = os.Getenv("GOPROJECT_CONFIG_FILE")
	}

	if configFile == "" {
		log.Fatal("config file is required")
	}

	log.Printf("use config file: %s", configFile)

	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Debug {
		log.Printf("config: %+v", config)
	}
}
