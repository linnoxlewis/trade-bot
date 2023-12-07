package telegram

import (
	"context"
	"github.com/linnoxlewis/trade-bot/pkg/telegram"
	"log"
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
	batchSize int
}

func New(tg *telegram.Client, userSrv UserSrv, orderSrv OrderSrv, batchSize int) Consumer {
	return Consumer{
		fetcher:   NewFetcher(tg),
		processor: NewProcessor(tg, userSrv, orderSrv, 0),
		batchSize: batchSize,
	}
}

func (c Consumer) Start(ctx context.Context) error {
	log.Println("start tg consumer")

	for {
		gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(ctx, gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
