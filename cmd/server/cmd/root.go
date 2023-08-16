package cmd

import (
	golog "log"
	"os"

	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/app/server"
	"github.com/liangjunmo/goproject/internal/app/server/serverenv"
	"github.com/liangjunmo/goproject/internal/version"
)

var (
	configFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	rootCmd.SetVersionTemplate("goproject-server:" + version.Describe())

	cobra.OnInitialize(func() {
		if configFile == "" {
			configFile = os.Getenv(serverenv.GOPROJECTServerConfigFile)
		}

		err := server.BuildConfig(configFile)
		if err != nil {
			golog.Fatal(err)
		}
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
		golog.Fatal(err)
	}
}
