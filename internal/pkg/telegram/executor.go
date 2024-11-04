package telegram

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper/conversion"
	"github.com/linnoxlewis/trade-bot/internal/pkg/validation"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	telegramCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
)

type UserSrv interface {
	CreateUser(ctx context.Context, username string, tgId int64) error
}

type OrderSrv interface {
	CreateOrder(ctx context.Context, order *dto.Order, tgUserId int64) (*domain.Order, error)
	CancelOrder(ctx context.Context, order *dto.CancelOrder, tgUserId int64) error
	GetUserActiveOrders(ctx context.Context, userId int64, exchange string) ([]domain.Order, error)
	UpdateTpslOrder(ctx context.Context, orderDto *dto.UpdateTpSl) error
}

type AccountSrv interface {
	GetBalance(ctx context.Context, userId int64, exchange string) (domain.Balances, error)
}

type Executor struct {
	tg         *telegramCli.TgClient
	validator  validation.Validator
	userSrv    UserSrv
	orderSrv   OrderSrv
	accountSrv AccountSrv
	clbrd      ClickBoard
	i18n       *i18n.I18n
}

func NewExecutor(tg *telegramCli.TgClient,
	validator validation.Validator,
	userSrv UserSrv,
	orderSrv OrderSrv,
	accountSrv AccountSrv,
	i18n *i18n.I18n) *Executor {
	return &Executor{
		tg:         tg,
		validator:  validator,
		userSrv:    userSrv,
		orderSrv:   orderSrv,
		accountSrv: accountSrv,
		clbrd:      ClickBoard{},
		i18n:       i18n,
	}
}

func (e *Executor) handleCreateOrder(ctx context.Context,
	chatID int,
	order *dto.TgOrder,
	lang string,
) error {
	dtoOrder := new(dto.Order)
	conversion.FromTgOrderToDtoOrder(order, dtoOrder)

	if err := e.validator.ValidateCreateOrder(dtoOrder); err != nil {
		return errors.BadRequestError(err.Error())
	}

	orderMdl, err := e.orderSrv.CreateOrder(ctx, dtoOrder, int64(chatID))
	if err != nil {
		return err
	}

	return e.tg.SendMessage(ctx, chatID, e.i18n.T(msgOrderCreated,
		map[string]interface{}{
			"ExecOrderId": orderMdl.ExecOrderId,
			"Id":          orderMdl.Id,
			"Side":        orderMdl.Side,
			"Exchange":    orderMdl.Exchange,
			"Symbol":      orderMdl.Symbol,
			"Price":       orderMdl.Price,
			"Quantity":    orderMdl.Quantity,
		},
		lang),
		"")
}

func (e *Executor) handleCancelOrder(ctx context.Context, chatID int, order *dto.TgOrder) error {
	dtoCancelOrder := new(dto.CancelOrder)
	conversion.FromTgOrderToDtoCancelOrder(order, dtoCancelOrder)

	if err := e.validator.ValidateCancelOrder(dtoCancelOrder); err != nil {
		return errors.BadRequestError(err.Error())
	}

	if err := e.orderSrv.CancelOrder(ctx, dtoCancelOrder, int64(chatID)); err != nil {
		return e.tg.SendMessage(ctx, chatID, msgNoCanceledOrder, "")
	}

	return e.tg.SendMessage(ctx, chatID, msgCanceledOrder, "")
}

func (e *Executor) handleUpdateTpSL(ctx context.Context, chatID int, order *dto.TgOrder) error {
	updateTpSL := new(dto.UpdateTpSl)
	conversion.FromTgOrderToDtoUpdateTpSlOrder(order, updateTpSL)

	if err := e.validator.ValidateTpSl(updateTpSL); err != nil {
		return errors.BadRequestError(err.Error())
	}

	if err := e.orderSrv.UpdateTpslOrder(ctx, updateTpSL); err != nil {
		return err
	}

	return e.tg.SendMessage(ctx, chatID, msgUpdateTpsSlOrder, "")
}

// sendHelp отправляет сообщение с помощью
func (e *Executor) sendHelp(ctx context.Context,
	chatID int,
	lang string,
) error {
	return e.tg.SendMessage(ctx,
		chatID, e.i18n.T(msgHelp, nil, lang), "")
}

// sendChooseBalanceExchange отправляет список бирж для выбора баланса
func (e *Executor) sendChooseBalanceExchange(ctx context.Context,
	chatID int,
	lang string,
) error {
	return e.tg.SendMessage(ctx, chatID, e.i18n.T("chooseExchange",
		nil,
		lang),
		e.clbrd.MakeBalanceExchangeKeyboard())
}

// sendChooseActiveOrdersExchange отправляет список бирж для выбора активных ордеров
func (e *Executor) sendChooseActiveOrdersExchange(ctx context.Context, chatID int, lang string) error {
	return e.tg.SendMessage(ctx, chatID, e.i18n.T("chooseExchange",
		nil,
		lang),
		e.clbrd.MakeActiveOrdersExchangeKeyboard())
}

// sendStart регистрирует пользователя и отправляет приветственное сообщение
func (e *Executor) sendStart(ctx context.Context, chatID int, lang string) error {
	if err := e.userSrv.CreateUser(ctx,
		"", //TODO: add username from tg
		int64(chatID)); err != nil {
		return err
	}

	return e.tg.SendMessage(ctx, chatID, e.i18n.T(msgHelp, nil, lang), "")
}
