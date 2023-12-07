package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	telegramCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
	"log"
	"net/url"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

/*
	{
	    "command": "create",
	    "exchange": "binance",
	    "ccy": "btcusdt",
	    "type":"limit",
	    "side": "market",
	    "qty":"0.001",
	    "price": "25000",
	    "tp": 0.01,
	    "sl": 0.03,
	    "ts": 0
	}
*/
var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type UserSrv interface {
	CreateUser(ctx context.Context, username string, tgId int64) (err error)
}

type OrderSrv interface {
	CreateOrder(ctx context.Context, order *dto.Order, tgUserId int64) (*domain.Order, error)
	CancelOrder(ctx context.Context, order *dto.CancelOrder, tgUserId int64) error
}

type Processor struct {
	tg       *telegramCli.Client
	userSrv  UserSrv
	orderSrv OrderSrv
	offset   int
}

func NewProcessor(tg *telegramCli.Client,
	userSrv UserSrv,
	orderSrv OrderSrv,
	offset int) *Processor {
	return &Processor{tg, userSrv, orderSrv, offset}
}

func (p *Processor) Process(ctx context.Context, event Event) error {
	switch event.Type {
	case Message:
		return p.processMessage(ctx, event)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) meta(event Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("can't get meta: " + ErrUnknownMetaType.Error())
	}

	return res, nil
}

func (p *Processor) processMessage(ctx context.Context, event Event) error {
	meta, err := p.meta(event)
	if err != nil {
		return errors.New("can't process message: " + err.Error())
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return errors.New("can't process message: " + err.Error())
	}

	return nil
}

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		if err := p.userSrv.CreateUser(ctx, username, int64(chatID)); err != nil {
			log.Println("can not save user: ", err)

			return err
		}

		return p.sendHello(ctx, chatID)
	default:
		var order *dto.Order
		if err := json.Unmarshal([]byte(text), &order); err != nil {
			p.tg.SendMessage(ctx, chatID, msgUnknownCommand)

			return err
		}

		orderMdl, err := p.orderSrv.CreateOrder(ctx, order, int64(chatID))
		if err != nil {
			p.tg.SendMessage(ctx, chatID, msgNoCreateOrder)

			return err
		}

		return p.tg.SendMessage(ctx, chatID, fmt.Sprintf(msgOrderCreated, order.Exchange, orderMdl.Symbol, orderMdl.Id))
	}
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
