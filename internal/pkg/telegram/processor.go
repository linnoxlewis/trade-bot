package telegram

import (
	"context"
	"encoding/json"
	errs "errors"
	"fmt"
	"strings"

	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	telegramCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
)

// Команды
const (
	helpCmd    = "/help"
	startCmd   = "/start"
	balanceCmd = "/balance"
	ordersCmd  = "/orders"
)

var (
	errUnknownEventType = errs.New("unknown event type")
	errUnknownMetaType  = errs.New("unknown meta type")
)

type Processor struct {
	tg       *telegramCli.TgClient
	executor *Executor
	logger   *log.Logger
	admins   []int
}

func NewProcessor(tg *telegramCli.TgClient,
	executor *Executor,
	logger *log.Logger,
	admins []int) *Processor {
	return &Processor{
		tg:       tg,
		executor: executor,
		admins:   admins,
		logger:   logger,
	}
}

func (p *Processor) Process(ctx context.Context, event Event) error {
	switch event.Type {
	case Message:
		return p.processMessage(ctx, event)
	default:
		return p.handleError(ctx, nil, errUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event Event) error {
	meta, err := p.extractMeta(event)
	if err != nil {
		return p.handleError(ctx, nil, err)
	}

	lang := helper.GetLang(meta.Lang)

	if err = p.executeCommand(ctx, strings.TrimSpace(event.Text), meta.ChatID, lang); err != nil {
		return p.handleError(ctx, &meta.ChatID, err)
	}

	return nil
}

func (p *Processor) extractMeta(event Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errUnknownMetaType
	}

	return meta, nil
}

func (p *Processor) executeCommand(ctx context.Context, text string, chatID int, lang string) error {
	switch text {
	case helpCmd:
		return p.executor.sendHelp(ctx, chatID, lang)
	case balanceCmd:
		return p.executor.sendChooseBalanceExchange(ctx, chatID, lang)
	case startCmd:
		return p.executor.sendStart(ctx, chatID, lang)
	case ordersCmd:
		return p.executor.sendChooseActiveOrdersExchange(ctx, chatID, lang)
	default:
		return p.handleJsonCommand(ctx, chatID, text, lang)
	}
}

func (p *Processor) handleJsonCommand(ctx context.Context, chatID int, text, lang string) error {
	var order *dto.TgOrder
	if err := json.Unmarshal([]byte(text), &order); err != nil {
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand, "")
	}

	switch order.Command {
	case consts.TgCreateOrderCommand:
		return p.executor.handleCreateOrder(ctx, chatID, order, lang)
	case consts.TgCancelOrderCommand:
		return p.executor.handleCancelOrder(ctx, chatID, order)
	case consts.TgUpdateTpSLCommand:
		return p.executor.handleUpdateTpSL(ctx, chatID, order)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand, "")
	}
}

func (p *Processor) handleError(ctx context.Context, chatID *int, err error) error {
	if err == nil {
		return nil
	}

	p.logger.ErrorLog.Printf("Error: %v", err)

	if chatID != nil {
		p.tg.SendMessage(ctx, *chatID, msgSomethingWrong, "")
	}

	for _, admin := range p.admins {
		p.tg.SendMessage(ctx, admin, fmt.Sprintf("Error: %v", err), "")
	}

	return err
}
