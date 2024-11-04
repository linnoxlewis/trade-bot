package main

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat/priceChecker"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat/priceChecker/storage"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"github.com/linnoxlewis/trade-bot/internal/service"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/linnoxlewis/trade-bot/pkg/postgresDb"
	_ "go.uber.org/automaxprocs"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sgn := make(chan os.Signal, 1)

	ctx := context.Background()
	cfg := config.NewConfig()

	logger := log.NewLogger()

	database := postgresDb.StartDB(cfg, logger)
	defer postgresDb.CloseDB(ctx, database, logger)

	orderRepo := repository.NewOrderRepository(database)

	symbolSrv := service.NewSymbolService(orderRepo, logger)
	ch := make(chan string)
	priceStorage := storage.NewLocalPriceStorage()
	newPriceChecker := priceChecker.NewCurrencyPairListener(
		symbolSrv,
		priceStorage,
		ch,
		logger)

	go newPriceChecker.Start(ctx)
	defer newPriceChecker.Stop()

	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-sgn:
	}
}
