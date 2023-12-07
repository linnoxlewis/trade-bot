package service

import (
	"context"
	"errors"
	"github.com/linnoxlewis/trade-bot/internal/pkg/excahnger"
	"math/big"

	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/pkg/log"
)

var errInvalidFormat = errors.New("err invalid format")
var errApiKeysNotFound = errors.New("err api keys not found")

type ApiKeysSrv interface {
	ClearApiKey(ctx context.Context, userId uuid.UUID, exchange string) error
	DeleteApiKey(ctx context.Context, userId uuid.UUID, exchange string) error
	AddApiKeys(ctx context.Context, userId uuid.UUID, apiKeys *dto.ApiKeys) error
	GetApiKeyByExchangeAndTgId(ctx context.Context, tgId int64, exchange string) (*domain.ApiKeys, error)
}

type Order struct {
	cfg        *config.Config
	exchanger  excahnger.Exchanger
	apiKeyRepo ApiKeyRepo
	logger     *log.Logger
}

func NewOrder(cfg *config.Config,
	exchanger excahnger.Exchanger,
	apiKeyRepo ApiKeyRepo,
	logger *log.Logger) *Order {
	return &Order{cfg: cfg,
		exchanger:  exchanger,
		apiKeyRepo: apiKeyRepo,
		logger:     logger,
	}
}

func (o *Order) CreateOrder(ctx context.Context, order *dto.Order, tgUserId int64) (*domain.Order, error) {
	keys, err := o.getApiKeys(ctx, tgUserId, order.Exchange)
	if err != nil {
		return nil, errApiKeysNotFound
	}

	if err := o.checkStopPercent(order.StopPercent); order.OrderType == consts.OrderTypeStopLossLimit && err != nil {
		return nil, err
	}

	var orderId int64
	var iceberg string

	if order.OrderType == consts.OrderTypeStopLossLimit {
		price, err := o.getPrice(order.StopPercent, order.StopPrice, order.Side)
		if err != nil {
			return nil, err
		}
		/*
			price, err = o.validatePrice(price, order.Symbol)
			if err != nil {
				return nil, err
			}
		*/
		order.Price = price

	}
	orderId, err = o.exchanger.CreateOrder(keys, order)
	if err != nil {
		return nil, err
	}

	result := &domain.Order{
		Id:          orderId,
		Symbol:      order.Symbol,
		Side:        order.Side,
		OrderType:   order.OrderType,
		Quantity:    order.Quantity,
		Price:       order.Price,
		TimeInForce: order.TimeInForce,
		StopPrice:   order.StopPrice,
	}

	if iceberg != "" {
		result.IcebergQty = iceberg
	}

	if err != nil {
		o.logger.ErrorLog.Println("err create order: " + err.Error())
		return nil, err
	}

	return result, err
}

func (o *Order) CancelOrder(ctx context.Context, order *dto.CancelOrder, tgUserId int64) error {
	keys, err := o.getApiKeys(ctx, tgUserId, order.Exchange)
	if err != nil {
		return err
	}

	return o.exchanger.CancelOrder(keys, order)
}

func (o *Order) checkStopPercent(stopPercent string) error {
	bigPer, ok := new(big.Float).SetString(stopPercent)
	if !ok {
		return errInvalidFormat
	}

	stopPer, _ := bigPer.Float64()
	if stopPer < 0 || stopPer > 100 {
		return errInvalidFormat
	}

	return nil
}

func (o *Order) getApiKeys(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	keys, err := o.apiKeyRepo.GetByApiKeyByTgUserIdAndExchange(ctx, int(userId), exchange)
	if err != nil {
		o.logger.ErrorLog.Println("err get api keys: ", err)
	}
	if keys == nil {
		return nil, errApiKeysNotFound
	}
	keys.DecodePassKey(o.cfg.GetSecretKey())
	keys.DecodePrivKey(o.cfg.GetSecretKey())

	return keys, nil
}

func (o *Order) getPrice(stopPercent string, stopPrice string, side string) (string, error) {
	bigPrice, ok := new(big.Float).SetString(stopPrice)
	if !ok {
		return "", errInvalidFormat
	}

	bigPercent, ok := new(big.Float).SetString(stopPercent)
	if !ok {
		return "", errInvalidFormat
	}

	percent, ok := new(big.Float).SetString("100")
	if !ok {
		return "", errInvalidFormat
	}

	result, ok := new(big.Float).SetString("0")
	if !ok {
		return "", errInvalidFormat
	}

	result.Mul(bigPercent, bigPrice)
	result.Quo(result, percent)

	if side == consts.OrderSideBuy {
		result.Add(bigPrice, result)
	} else {
		result.Sub(bigPrice, result)
	}

	return result.String(), nil
}
