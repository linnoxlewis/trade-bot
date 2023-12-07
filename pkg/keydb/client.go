package keydb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var conn *redis.Client

func NewConnection(address string, password string, dbNum int, logger *log.Logger) *redis.Client {
	conn = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbNum,
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	logger.InfoLog.Println("key-db connection establish")

	return conn
}

func CloseConnection(conn *redis.Client, logger *log.Logger) {
	logger.InfoLog.Println("close key-db connection")
	if err := conn.Close(); err != nil {
		logger.ErrorLog.Panic(err)
	}
}
