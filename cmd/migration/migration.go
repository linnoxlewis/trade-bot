package migration

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/linnoxlewis/trade-bot/pkg/postgresDb"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

const path = "schema"

var mgrCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run postgresDb migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		logger := log.NewLogger()
		database := postgresDb.StartDB(cfg, logger)

		defer postgresDb.CloseDB(context.Background(), database, logger)

		if err := goose.Run(args[0], database, path, args[1:]...); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mgrCmd)
}
