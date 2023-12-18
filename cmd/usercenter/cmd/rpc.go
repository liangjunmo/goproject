package cmd

import (
	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/usercenter"
)

func init() {
	rootCmd.AddCommand(rpcCmd)
}

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Run usercenter rpc server",
	Run: func(cmd *cobra.Command, args []string) {
		usercenter.RunRPCServer(usercenter.RPCServerConfig{
			Environment: config.Environment,
			Debug:       config.Debug,
			Addr:        config.RPCServerAddr,
			DB: usercenter.DBConfig{
				Addr:     config.DB.Addr,
				User:     config.DB.User,
				Password: config.DB.Password,
				Database: config.DB.Database,
			},
			Redis: usercenter.RedisConfig{
				Addr:     config.Redis.Addr,
				Password: config.Redis.Password,
			},
		})
	},
}
