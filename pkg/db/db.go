package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"tct_backend/auth/config"
	"tct_backend/auth/pkg/log"
)

var db *gorm.DB

const dbUri = "host=%s user=%s dbname=%s port=%s sslmode=disable password=%s"

func newDb(cfg *config.Config, logger *log.Logger) *gorm.DB {
	url := fmt.Sprintf(dbUri,
		cfg.GetDbHost(),
		cfg.GetDbUsername(),
		cfg.GetDbName(),
		cfg.GetDbPort(),
		cfg.GetDbPassword(),
	)

	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	return database
}

func StartDB(cfg *config.Config, logger *log.Logger) *gorm.DB {
	if db == nil {
		db = newDb(cfg, logger)
	}
	logger.InfoLog.Println("Connecting to database...")

	return db
}

func CloseDB(db *gorm.DB, logger *log.Logger) {
	logger.InfoLog.Println("Close database Connection")

	database, err := db.DB()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	if err = database.Close(); err != nil {
		logger.ErrorLog.Panic(err)
	}
}
