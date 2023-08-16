package cmd

import (
	"context"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/liangjunmo/goproject/internal/app/server/serverapi"
	"github.com/liangjunmo/goproject/internal/app/server/serverconfig"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "goproject-server api cli tool",
	Long:  "goproject-server api cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()

		release := serverapi.Build(router)
		defer release()

		server := &http.Server{
			Addr:    serverconfig.Config.Api.Addr,
			Handler: router,
		}

		go func() {
			err := server.ListenAndServe()
			if err == http.ErrServerClosed {
				golog.Println("http server closed")
			} else {
				golog.Fatal(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			golog.Fatal(err)
		}
	},
}
