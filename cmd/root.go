package cmd

import (
	"context"
	"github.com/linnoxlewis/trade-bot/app"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "Init trade-bot",
	Short: "Init trade-bot application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			sgn := make(chan os.Signal, 1)
			signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-ctx.Done():
			case <-sgn:
			}
			cancel()
		}()

		app.Run(ctx)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
