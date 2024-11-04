package app

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

func Run(ctx context.Context) {
	sgn := make(chan os.Signal, 1)
	cfg := config.NewConfig()
	logger := log.NewLogger()

	i18nLocale := i18n.NewI18n(consts.LocaleDataPath, helper.GetLanguageList())

	keyDb := keydb.NewConnection(
		cfg.GetGlobalCacheAddress(),
		cfg.GetGlobalCachePassword(),
		consts.GlobalCacheDb0,
		logger)
	defer keydb.CloseConnection(keyDb, logger)

	database := postgresDb.StartDB(cfg, logger)
	defer postgresDb.CloseDB(ctx, database, logger)

	exchangePkg := exchanger.NewExchanger(cfg.GetTestnet())

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

	// orderSrv := service.NewOrder()

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

	tg := tgCli.NewTgClient(cfg.GetTgToken())
	tgBot := telegram.New(tg,
		userSrv,
		nil,
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
	httpServer := httpServer.NewHttpServer(cfg, externalMux)
	go func() {
		logger.InfoLog.Println("external server is listening on " + cfg.GetApiPort())

		if err := httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()
	defer func(httpServer *http.Server, ctx context.Context) {
		err := httpServer.Shutdown(ctx)
		if err != nil {

		}
	}(httpServer, ctx)

	newPairsChan := make(chan string)
	priceStorage, err := storage.NewRedisPriceStorage("localhost:6379")
	if err != nil {
		logger.ErrorLog.Fatalf("Error initializing RedisPriceStorage: %v", err)
	}
	priceChecker := priceChecker.NewCurrencyPairListener(symbolSrv,
		priceStorage,
		newPairsChan,
		logger)
	go func() {
		if err = priceChecker.Start(ctx); err != nil {
			panic(err)
		}
	}()
	defer priceChecker.Stop()

	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sgn:
	}
}
