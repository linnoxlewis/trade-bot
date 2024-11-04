package main

import (
	"context"
	"errors"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat/priceChecker"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat/priceChecker/storage"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/pkg/exchanger"
	"github.com/linnoxlewis/trade-bot/internal/pkg/telegram"
	"github.com/linnoxlewis/trade-bot/internal/pkg/validation"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"github.com/linnoxlewis/trade-bot/internal/service"
	"github.com/linnoxlewis/trade-bot/internal/transport/http/handlers"
	httpServer "github.com/linnoxlewis/trade-bot/internal/transport/http/server"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/keydb"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/linnoxlewis/trade-bot/pkg/postgresDb"
	tgCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Main(ctx context.Context) {
	sgn := make(chan os.Signal, 1)
	cfg := config.NewConfig()
	logger := log.NewLogger()

	validator := validation.NewValidator()

	i18nLocale := i18n.NewI18n(consts.LocaleDataPath, helper.GetLanguageList())

	cacheStorage := keydb.NewConnection(cfg.CacheStorage, logger)
	defer keydb.CloseConnection(cacheStorage, logger)

	database := postgresDb.StartDB(cfg.DBConfig, logger)
	defer postgresDb.CloseDB(database, logger)

	exchangePkg := exchanger.NewExchanger(cfg.BinanceConfig.GetTestnet())

	apiKeysRepo := repository.NewApiKeyRepo(database)
	userRepo := repository.NewUserRepository(database)
	orderRepo := repository.NewOrderRepository(database)

	userSrv := service.NewUserService(cfg,
		userRepo,
		i18nLocale,
		logger)

	apiKeySrv := service.NewApiKeysService(cfg,
		apiKeysRepo,
		userRepo,
		i18nLocale,
		logger)
	/*
		orderSrv := service.NewOrder(cfg,
			exchangePkg,
			apiKeysRepo,
			orderRepo,
			binanceQueue,
			limitBinanceQueue,
			keyDb,
			i18nLocale,
			logger)
	*/
	symbolSrv := service.NewSymbolService(orderRepo, logger)

	accountSrv := service.NewAccountService(cfg,
		apiKeysRepo,
		exchangePkg,
		i18nLocale,
		logger)

	admins, err := userSrv.GetAdmins(ctx)
	if err != nil {
		panic("cant get admin list")
	}

	tg := tgCli.NewTgClient(cfg.TelegramConfig.GetTgToken())
	tgBot := telegram.NewTgBot(tg,
		validator,
		userSrv,
		orderSrv,
		accountSrv,
		i18nLocale,
		admins,
		logger,
		100)
	go func() {
		if err := tgBot.Start(ctx); err != nil {
			panic("service is stopped " + err.Error())
		}
	}()

	externalHandlers := handlers.NewApiKeyHandler(apiKeySrv, logger)
	externalMux := httpServer.NewExternalMux(externalHandlers)
	httpExternalServer := httpServer.NewHttpServer(cfg, externalMux)
	go func() {
		logger.InfoLog.Println("external server is listening on " + cfg.ServerConfig.GetPort())

		if err := httpExternalServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	defer func(httpServer *http.Server, ctx context.Context) {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.ErrorLog.Fatalf("Error shutdown server: %v", err)
		}
	}(httpExternalServer, ctx)

	newPairsChan := make(chan string)

	priceStorage, err := storage.NewRedisPriceStorage("localhost:6379")
	if err != nil {
		logger.ErrorLog.Fatalf("Error initializing RedisPriceStorage: %v", err)
	}

	priceCheckerListener := priceChecker.NewCurrencyPairListener(symbolSrv,
		priceStorage,
		newPairsChan,
		logger)

	go func() {
		if err = priceCheckerListener.Start(ctx); err != nil {
			panic(err)
		}
	}()
	defer priceCheckerListener.Stop()

	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-sgn:
	}
}

/*
func startTickers(ctx context.Context,
	symbols domain.SymbolList,
	sgn chan os.Signal,
	cfg *config.Config,
	keyDb *redis.Client,
	orderSrv *service.Order,
	tg *tgCli.TgClient,
	binanceQueue *domain.OrdersQueue,
	limitBinanceQueue *domain.OrdersQueue,
	i18nLocale *i18n.I18n,
	logger *log.Logger) {
	for _, v := range exchanger.ExchangeList {
		exchange := v

		tpSlTicker := heartbeat.NewTpSlTicker(cfg,
			orderSrv,
			keyDb,
			logger,
			tg,
			time.Second,
			binanceQueue,
			v,
			i18nLocale,
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
}*/
