package keydb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var conn *redis.Client

func NewConnection(cfg config.CacheStorage, logger *log.Logger) *redis.Client {
	conn = redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddress(),
		Password: cfg.GetPassword(),
		DB:       cfg.GetDB(),
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	logger.InfoLog.Println("key-postgresDb connection establish")

	return conn
}

func CloseConnection(conn *redis.Client, logger *log.Logger) {
	if err := conn.Close(); err != nil {
		logger.ErrorLog.Panic(err)
	}

	logger.InfoLog.Println("close key-postgresDb connection")
}
