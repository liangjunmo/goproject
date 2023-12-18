package cmd

import (
	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/goproject"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run api server",
	Run: func(cmd *cobra.Command, args []string) {
		goproject.RunAPIServer(goproject.APIServerConfig{
			Environment:             config.Environment,
			Debug:                   config.Debug,
			JWTKey:                  config.JWTKey,
			Addr:                    config.APIAServerAddr,
			UserCenterRPCServerAddr: config.UserCenterRPCServerAddr,
			DB: goproject.DBConfig{
				Addr:     config.DB.Addr,
				User:     config.DB.User,
				Password: config.DB.Password,
				Database: config.DB.Database,
			},
			Redis: goproject.RedisConfig{
				Addr:     config.Redis.Addr,
				Password: config.Redis.Password,
			},
		})
	},
}
