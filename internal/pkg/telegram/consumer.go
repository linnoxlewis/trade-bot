package telegram

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/pkg/validation"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/linnoxlewis/trade-bot/pkg/telegram"
	"time"
)

type ConsumerInterface interface {
	Start() error
}

type FetcherInterface interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type ProcessorInterface interface {
	Process(ctx context.Context, event Event) error
}

type TgBot struct {
	fetcher   FetcherInterface
	processor ProcessorInterface
	logger    *log.Logger
	batchSize int
}

func NewTgBot(tg *telegram.TgClient,
	validator validation.Validator,
	userSrv UserSrv,
	orderSrv OrderSrv,
	accountSrv AccountSrv,
	i18n *i18n.I18n,
	admins []int,
	logger *log.Logger,
	batchSize int) *TgBot {
	executor := NewExecutor(tg,
		validator,
		userSrv,
		orderSrv,
		accountSrv,
		i18n)

	processor := NewProcessor(tg,
		executor,
		logger,
		admins)

	return &TgBot{
		fetcher:   NewFetcher(tg),
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (t *TgBot) Start(ctx context.Context) error {
	t.logger.InfoLog.Println("start tg consumer")

	for {
		select {
		case <-ctx.Done():
			t.logger.InfoLog.Println("tg bot stopped")

			return ctx.Err()
		default:

		}

		gotEvents, err := t.fetcher.Fetch(ctx, t.batchSize)
		if err != nil {
			t.logger.ErrorLog.Printf("err fetch tg boy: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := t.handleEvents(ctx, gotEvents); err != nil {
			t.logger.ErrorLog.Print(err)
			continue
		}
	}
}

func (t *TgBot) handleEvents(ctx context.Context, events []Event) error {
	for _, event := range events {
		t.logger.InfoLog.Printf("got new event: %s", event.Text)

		if err := t.processor.Process(ctx, event); err != nil {
			t.logger.ErrorLog.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
