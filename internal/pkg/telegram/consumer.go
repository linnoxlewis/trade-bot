package telegram

import (
	"context"
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

type Consumer struct {
	fetcher   FetcherInterface
	processor ProcessorInterface
	logger    *log.Logger
	batchSize int
}

func New(tg *telegram.Client,
	userSrv UserSrv,
	orderSrv OrderSrv,
	accountSrv AccountSrv,
	i18n *i18n.I18n,
	admins []int,
	logger *log.Logger,
	batchSize int) Consumer {
	return Consumer{
		fetcher:   NewFetcher(tg),
		processor: NewProcessor(tg, userSrv, orderSrv, accountSrv, i18n, logger, admins, 0),
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c Consumer) Start(ctx context.Context) error {
	c.logger.InfoLog.Println("start tg consumer")

	for {
		gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
		if err != nil {
			c.logger.ErrorLog.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(ctx, gotEvents); err != nil {
			c.logger.ErrorLog.Print(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []Event) error {
	for _, event := range events {
		c.logger.InfoLog.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			c.logger.ErrorLog.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
