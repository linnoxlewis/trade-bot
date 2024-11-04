package app

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/heartbeat"
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
	"time"
)

func Run(ctx context.Context) {
	sgn := make(chan os.Signal, 1)
	cfg := config.NewConfig()
	logger := log.NewLogger()

	i18nLocale := i18n.NewI18n(consts.LocaleDataPath, helper.GetLanguageList())

	keyDb := keydb.NewConnection(
		cfg.CacheStorage.GetAddress(),
		cfg.CacheStorage.GetPassword(),
		cfg.CacheStorage.GetDB(),
		logger)
	defer keydb.CloseConnection(keyDb, logger)

	database := postgresDb.StartDB(cfg, logger)
	defer postgresDb.CloseDB(ctx, database, logger)

	binanceQueue := domain.NewOrderQueue(consts.Binance)
	limitBinanceQueue := domain.NewOrderQueue(consts.Binance)

	exchangePkg := exchanger.NewExchanger(cfg.GetTestnet())

	apiKeysRepo := repository.NewApiKeyRepo(database)
	userRepo := repository.NewUserRepository(database)
	orderRepo := repository.NewOrderRepository(database)
	symbolsRepo := repository.NewSymbolRepository(database)

	userSrv := service.NewUserService(cfg,
		userRepo,
		i18nLocale,
		logger)

	apiKeySrv := service.NewApiKeysService(cfg,
		apiKeysRepo,
		userRepo,
		i18nLocale,
		logger)

	orderSrv := service.NewOrder(cfg,
		exchangePkg,
		apiKeysRepo,
		orderRepo,
		binanceQueue,
		limitBinanceQueue,
		keyDb,
		i18nLocale,
		logger)

	symbolSrv := service.NewSymbolService(orderRepo,
		symbolsRepo,
		logger)

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
	tgBot := telegram.New(tg, userSrv, orderSrv, accountSrv, i18nLocale, admins, logger, 100)
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

	symbols, _ := symbolSrv.GetActiveSymbols(ctx)
	if len(symbols) == 0 {
		symbols, err = symbolSrv.GetDefaultSymbols(ctx)
		if err != nil || symbols == nil {
			panic("cant get symbols for ticker")
		}
	}

	startTickers(ctx,
		symbols,
		sgn,
		cfg,
		keyDb,
		orderSrv,
		tg,
		binanceQueue,
		limitBinanceQueue,
		i18nLocale, logger)

	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-sgn:
	}
}

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
		for _, symbolVal := range symbols {
			priceTicker, err := heartbeat.NewTradePriceTicker(cfg,
				keyDb,
				exchange,
				string(symbolVal),
				logger)
			if err != nil {
				logger.ErrorLog.Println("heartbeat error:", err)

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
}
