package telegram

import (
	"context"
	"encoding/json"
	errs "errors"
	"fmt"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/helper/conversion"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	telegramCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
	"strings"
)

const (
	HelpCmd    = "/help"
	StartCmd   = "/start"
	BalanceCmd = "/balance"
	OrdersCmd  = "/orders"

	activeOrdersExchangeBinanceCmd = "active_orders_exchange_binance"
	activeOrdersExchangeKucoinCmd  = "active_orders_exchange_kucoin"
	activeOrdersExchangeOkxCmd     = "active_orders_exchange_okx"
	balanceExchangeBinanceCmd      = "balance_exchange_binance"
	balanceExchangeKucoinCmd       = "balance_exchange_kucoin"
	balanceExchangeOkxCmd          = "balance_exchange_okx"
)

var (
	ErrUnknownEventType = errs.New("unknown event type")
	ErrUnknownMetaType  = errs.New("unknown meta type")
)

type UserSrv interface {
	CreateUser(ctx context.Context, username string, tgId int64) (err error)
}

type OrderSrv interface {
	CreateOrder(ctx context.Context, order *dto.Order, tgUserId int64) (*domain.Order, error)
	CancelOrder(ctx context.Context, order *dto.CancelOrder, tgUserId int64) error
	GetActiveTpSlOrders(ctx context.Context, exchange string) ([]*domain.Order, error)
	ExecuteTpSlOrder(ctx context.Context, userId int64, order *domain.Order) (int64, error)
	UpdateTpslOrder(ctx context.Context, orderDto *dto.UpdateTpSl) error
	GetUserActiveOrders(ctx context.Context, userId int64, exchange string) ([]domain.Order, error)
	GetOrder(ctx context.Context, orderId int64, tgUserId int64, symbol, exchange string, inExchange bool) (*domain.Order, error)
	SetFilledLimitOrder(ctx context.Context, order *domain.Order) error
	GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error)
}

type AccountSrv interface {
	GetBalance(ctx context.Context, userId int64, exchange string) (balance domain.Balance, err error)
}

type Processor struct {
	tg         *telegramCli.Client
	userSrv    UserSrv
	orderSrv   OrderSrv
	accountSrv AccountSrv
	i18n       *i18n.I18n
	logger     *log.Logger
	clbrd      ClickBoard
	admins     []int
	offset     int
}

func NewProcessor(tg *telegramCli.Client,
	userSrv UserSrv,
	orderSrv OrderSrv,
	accountSrv AccountSrv,
	i18n *i18n.I18n,
	logger *log.Logger,
	admins []int,
	offset int) *Processor {
	return &Processor{
		tg,
		userSrv,
		orderSrv,
		accountSrv,
		i18n,
		logger,
		ClickBoard{},
		admins,
		offset}
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
		return Meta{}, errs.New("can't get meta: " + ErrUnknownMetaType.Error())
	}

	return res, nil
}

func (p *Processor) processMessage(ctx context.Context, event Event) error {
	meta, err := p.meta(event)
	if err != nil {
		return errs.New("can't process message: " + err.Error())
	}

	/*defer func() {
		if pn := recover(); pn != nil {
			p.tg.DeleteMessage(ctx, meta.ChatID, meta.MessageId)
		}
	}()*/

	lang := helper.GetLang(meta.Lang)

	if err = p.doCmd(ctx, event.Text, meta.ChatID, lang); err != nil {
		return errs.New("can't process message: " + err.Error())
	}

	return nil
}

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, lang string) error {
	var err error
	defer func() {
		if pn := recover(); pn != nil {
			p.logger.ErrorLog.Printf("tg panic: ", pn)

			go p.tg.SendMessage(ctx, chatID, msgSomethingWrong, "")
			for _, v := range p.admins {
				p.tg.SendMessage(ctx, v, fmt.Sprintf("panic detected: %v", pn), "")
			}
			panic(pn)
		}
		if err != nil {
			go p.sendError(ctx, chatID, err)
		}
	}()

	text = strings.TrimSpace(text)
	p.logger.InfoLog.Printf("got new command '%s' from '%s", text, chatID)

	switch text {
	case HelpCmd:
		err = p.sendHelp(ctx, chatID, lang)
		break
	case BalanceCmd:
		err = p.sendChooseBalanceExchange(ctx, chatID, lang)
		break
	case StartCmd:
		err = p.sendStart(ctx, chatID, lang)
		break
	case OrdersCmd:
		err = p.sendChooseActiveOrdersExchange(ctx, chatID, lang)
		break
	case activeOrdersExchangeBinanceCmd:
		err = p.sendActiveOrders(ctx, consts.Binance, chatID, lang)
		break
	case activeOrdersExchangeKucoinCmd:
		err = p.sendActiveOrders(ctx, consts.Kucoin, chatID, lang)
		break
	case activeOrdersExchangeOkxCmd:
		err = p.sendActiveOrders(ctx, consts.Okx, chatID, lang)
		break
	case balanceExchangeBinanceCmd:
		err = p.sendBalance(ctx, chatID, consts.Binance, lang)
		break
	case balanceExchangeKucoinCmd:
		err = p.sendBalance(ctx, chatID, consts.Kucoin, lang)
		break
	case balanceExchangeOkxCmd:
		err = p.sendBalance(ctx, chatID, consts.Okx, lang)
		break

	default:
		var order *dto.TgOrder
		if err = json.Unmarshal([]byte(text), &order); err != nil {
			p.tg.SendMessage(ctx, chatID, msgUnknownCommand, "")

			return nil
		}
		err = p.doJsonCmd(ctx, chatID, order, lang)
		break
	}

	return err
}

func (p *Processor) doJsonCmd(ctx context.Context, chatID int, order *dto.TgOrder, lang string) error {
	switch order.Command {
	case consts.TgCreateOrderCommand:
		dtoCreateOrder := new(dto.Order)
		conversion.FromTgOrderToDtoOrder(order, dtoCreateOrder)
		if err := dtoCreateOrder.Validate(); err != nil {
			return errors.BadRequestError(err.Error())
		}
		var orderMdl *domain.Order
		orderMdl, err := p.orderSrv.CreateOrder(ctx, dtoCreateOrder, int64(chatID))
		if err != nil {
			return err
		}

		return p.tg.SendMessage(ctx, chatID, p.i18n.T(msgOrderCreated, map[string]interface{}{
			"ExecOrderId": orderMdl.ExecOrderId,
			"Id":          orderMdl.Id,
			"Side":        orderMdl.Side,
			"Exchange":    orderMdl.Exchange,
			"Symbol":      orderMdl.Symbol,
			"Price":       orderMdl.Price,
			"Quantity":    orderMdl.Quantity,
		}, lang), "")

	case consts.TgCancelOrderCommand:
		dtoCancelOrder := new(dto.CancelOrder)
		conversion.FromTgOrderToDtoCancelOrder(order, dtoCancelOrder)
		if err := dtoCancelOrder.Validate(); err != nil {
			return errors.BadRequestError(err.Error())
		}
		if err := p.orderSrv.CancelOrder(ctx, dtoCancelOrder, int64(chatID)); err != nil {
			p.tg.SendMessage(ctx, chatID, msgNoCanceledOrder, "")

			return err
		}
		go p.tg.SendMessage(ctx, chatID, fmt.Sprintf(msgCanceledOrder), "")

		return nil

	case consts.TgUpdateTpSLCommand:
		updateTpSL := new(dto.UpdateTpSl)
		conversion.FromTgOrderToDtoUpdateTpSlOrder(order, updateTpSL)
		if err := updateTpSL.Validate(); err != nil {
			return errors.BadRequestError(err.Error())
		}
		if err := p.orderSrv.UpdateTpslOrder(ctx, updateTpSL); err != nil {
			return err
		}
		go p.tg.SendMessage(ctx, chatID, fmt.Sprintf(msgUpdateTpsSlOrder), "")

		return nil
	default:
		return nil
	}
}

func (p *Processor) sendHelp(ctx context.Context, chatID int, lang string) error {
	return p.tg.SendMessage(ctx, chatID, p.i18n.T(msgHelp, nil, lang), "")
}

func (p *Processor) sendChooseActiveOrdersExchange(ctx context.Context, chatID int, lang string) error {
	return p.tg.SendMessage(ctx, chatID, p.i18n.T("chooseExchange", nil, lang), p.clbrd.MakeActiveOrdersExchangeKeyboard())
}

func (p *Processor) sendChooseBalanceExchange(ctx context.Context, chatID int, lang string) error {
	return p.tg.SendMessage(ctx, chatID, p.i18n.T("chooseExchange", nil, lang), p.clbrd.MakeBalanceExchangeKeyboard())
}

func (p *Processor) sendActiveOrders(ctx context.Context, exchange string, chatID int, lang string) error {
	orders, err := p.orderSrv.GetUserActiveOrders(ctx, int64(chatID), exchange)
	if err != nil {
		return err
	}
	msg := p.i18n.T("yourOrders", nil, lang) + "\n"
	for i := 0; i < len(orders); i++ {
		msg += fmt.Sprintf("ID:%v\nsymbol:%s\nside:%s\ntype:%s\nquantity:%s\nprice:%s\nstatus:%s\n",
			orders[i].ExecOrderId,
			orders[i].Symbol,
			orders[i].Side,
			orders[i].OrderType,
			orders[i].Quantity,
			orders[i].Price,
			orders[i].Status) + "\n"
		msg += "---------------------------"
	}
	go func() {
		if err := p.tg.SendMessage(ctx, chatID, msg, ""); err != nil {
			p.logger.ErrorLog.Println(err)
		}
	}()

	return nil
}

func (p *Processor) sendBalance(ctx context.Context, chatID int, exchange, lang string) error {
	balance, err := p.accountSrv.GetBalance(ctx, int64(chatID), exchange)
	if err != nil {
		return err
	}
	msg := p.i18n.T("yourBalance", nil, lang) + "\n"
	for _, v := range balance {
		msg += fmt.Sprintf("%s: %s \n", v.Symbol, v.Quantity)
	}
	go func() {
		if err := p.tg.SendMessage(ctx, chatID, msg, ""); err != nil {
			p.logger.ErrorLog.Println(err)
		}
	}()

	return nil
}

func (p *Processor) sendStart(ctx context.Context, chatID int, lang string) error {
	if err := p.userSrv.CreateUser(ctx, "", int64(chatID)); err != nil {
		return err
	}
	go func() {
		if err := p.tg.SendMessage(ctx,
			chatID,
			p.i18n.T(msgHelp, nil, lang), ""); err != nil {
			p.logger.ErrorLog.Println(err)
		}
	}()

	return nil
}

func (p *Processor) sendError(ctx context.Context, chatID int, err error) {
	if errs.As(err, &errors.Error{}) {
		err := err.(errors.Error)
		if err.IsInternalServerError() {
			go p.tg.SendMessage(ctx, chatID, msgSomethingWrong, "")
			for _, v := range p.admins {
				p.tg.SendMessage(ctx, v, err.Error(), "")
			}
		} else {
			p.tg.SendMessage(ctx, chatID, err.Error(), "")
		}
	} else {
		p.tg.SendMessage(ctx, chatID, msgSomethingWrong, "")
	}
}
