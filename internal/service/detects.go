package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/pkg/excahnger"
)

type DetectsSrv struct {
	cfg       *config.Config
	exchanger excahnger.Exchanger
	keyDbCli  *redis.Client
	exchange  string
}

// TODO
func (d *DetectsSrv) DetectPump(symbol, exchange string) {}
