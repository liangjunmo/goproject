package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/app/server/serverworker1"
)

func init() {
	rootCmd.AddCommand(worker1Cmd)
}

var worker1Cmd = &cobra.Command{
	Use:   "worker1",
	Short: "goproject-server worker1 cli tool",
	Long:  "goproject-server worker1 cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}

		release := serverworker1.Build(ctx, wg)
		defer release()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c

		cancel()
		wg.Wait()
	},
}
