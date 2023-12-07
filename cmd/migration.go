package cmd

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/pkg/db"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

const path = "schema"

var mgrCmd = &cobra.Command{
	Use: "migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		logger := log.NewLogger()
		database := db.StartDB(cfg, logger)

		defer db.CloseDB(context.Background(), database, logger)

		if err := goose.Run(args[0], database, path, args[1:]...); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mgrCmd)
}
