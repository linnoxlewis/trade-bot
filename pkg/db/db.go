package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

const dbUri = "host=%s user=%s dbname=%s port=%s sslmode=disable password=%s"

func newDb(cfg *config.Config, logger *log.Logger) *sql.DB {
	url := fmt.Sprintf(dbUri,
		cfg.GetDbHost(),
		cfg.GetDbUsername(),
		cfg.GetDbName(),
		cfg.GetDbPort(),
		cfg.GetDbPassword(),
	)

	db, err := sql.Open("pgx", url)
	if err != nil {
		logger.ErrorLog.Panic("Cant connect db: ", err)
		os.Exit(1)
	}

	return db
}

func StartDB(cfg *config.Config, logger *log.Logger) *sql.DB {
	var once sync.Once
	if db == nil {
		once.Do(func() {
			db = newDb(cfg, logger)
		})
	} else {
		logger.ErrorLog.Println("Single instance already created.")
		return db
	}

	logger.InfoLog.Println("Connecting to database...")

	return db
}

func CloseDB(ctx context.Context, db *sql.DB, logger *log.Logger) {
	logger.InfoLog.Println("Close database Connection")
	if err := db.Close(); err != nil {
		logger.ErrorLog.Panic(err)
	}
}
