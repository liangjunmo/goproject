package cmd

import (
	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/goproject"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker server",
	Run: func(cmd *cobra.Command, args []string) {
		goproject.RunWorkerServer(goproject.WorkerServerConfig{
			Environment:             config.Environment,
			Debug:                   config.Debug,
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
