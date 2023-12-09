package app

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/pkg/excahnger"
	"github.com/linnoxlewis/trade-bot/internal/pkg/telegram"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"github.com/linnoxlewis/trade-bot/internal/service"
	restServer "github.com/linnoxlewis/trade-bot/internal/transport/api/server"
	grpcServer "github.com/linnoxlewis/trade-bot/internal/transport/grpc/server"
	"github.com/linnoxlewis/trade-bot/pkg/db"
	i18n2 "github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/keydb"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	tgCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context) {
	sgn := make(chan os.Signal, 1)
	cfg := config.NewConfig()
	logger := log.NewLogger()

	i18n := i18n2.NewI18n(consts.LocaleDataPath, helper.GetLanguageList())

	keyDb := keydb.NewConnection(
		cfg.GetGlobalCacheAddress(),
		cfg.GetGlobalCachePassword(),
		consts.GlobalCacheDb0,
		logger)
	defer keydb.CloseConnection(keyDb, logger)

	database := db.StartDB(cfg, logger)
	defer db.CloseDB(ctx, database, logger)

	binanceQueue := domain.NewOrderQueue(consts.Binance)
	limitBinanceQueue := domain.NewOrderQueue(consts.Binance)

	exchangePkg := excahnger.NewExchanger(cfg.GetTestnet())

	apiKeysRepo := repository.NewApiKeyRepo(database)
	userRepo := repository.NewUserRepository(database)
	orderRepo := repository.NewOrderRepository(database)
	symbolsRepo := repository.NewSymbolRepository(database)

	userSrv := service.NewUserService(cfg, userRepo, i18n, logger)
	apiKeySrv := service.NewApiKeysService(cfg, apiKeysRepo, userRepo, i18n, logger)
	orderSrv := service.NewOrder(cfg, exchangePkg, apiKeysRepo, orderRepo, binanceQueue, limitBinanceQueue, keyDb, i18n, logger)
	symbolSrv := service.NewSymbolService(orderRepo, symbolsRepo, logger)
	accountSrv := service.NewAccountService(cfg, apiKeysRepo, exchangePkg, i18n, logger)

	admins, err := userSrv.GetAdmins(ctx)
	if err != nil {
		panic("cant get admin list")
	}

	tg := tgCli.New(cfg.GetTgToken())
	tgBot := telegram.New(tg, userSrv, orderSrv, accountSrv, i18n, admins, logger, 100)
	go func() {
		if err := tgBot.Start(ctx); err != nil {
			panic("service is stopped " + err.Error())
		}
	}()

	grpcSrv := grpcServer.NewGrpc(cfg.GetGrpcPort(), *logger)
	go grpcSrv.StartServer()
	defer grpcSrv.StopServer()

	restSrv := restServer.NewApiServer(cfg.GetApiPort(), cfg.GetApiServerMode(), apiKeySrv, logger)
	go restSrv.StartServer()
	defer restSrv.StopServer()

	symbols, _ := symbolSrv.GetActiveSymbols(ctx)
	if symbols.IsEmpty() {
		var err error
		symbols, err = symbolSrv.GetDefaultSymbols(ctx)
		if err != nil || symbols == nil {
			panic("cant get symbols for ticker")
		}
	}

	for _, v := range excahnger.ExchangeList {
		exchange := v
		for _, symbolVal := range symbols {
			priceTicker, err := heartbeat.NewTradePriceTicker(cfg,
				keyDb,
				exchange,
				string(symbolVal),
				logger)
			if err != nil {
				continue
			}
			go priceTicker.HandleWebSocketData(ctx)
		}

		tpSlTicker := heartbeat.NewTpSlTicker(cfg,
			orderSrv,
			keyDb,
			logger,
			tg,
			time.Second,
			binanceQueue,
			v,
			i18n,
			true)
		go tpSlTicker.Tick(ctx, sgn)

		checkLimitOrders := heartbeat.NewLimitOrderTicker(cfg,
			orderSrv,
			limitBinanceQueue,
			logger,
			v,
			time.Second*5)
		go checkLimitOrders.Tick(ctx, sgn)
	}

	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sgn:
	}
}
