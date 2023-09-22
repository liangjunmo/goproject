package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/service/userservice"
	"github.com/liangjunmo/goproject/internal/worker/userworker"
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

		release := buildWorker1(ctx, wg)
		defer release()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c

		cancel()
		wg.Wait()
	},
}

func buildWorker1(ctx context.Context, wg *sync.WaitGroup) (release func()) {
	db := connectDb(true)

	redisClient := connectRedis()

	release = func() {
		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}

	redisSync := newRedisSync(redisClient)

	userService := userservice.NewService(db, redisSync)

	wg.Add(1)
	go userworker.NewListUserWorker(db, userService).Run(ctx, wg)

	return
}
