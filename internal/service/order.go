package service

import (
	"context"
	"github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"

	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/helper/conversion"
	"github.com/linnoxlewis/trade-bot/internal/pkg/exchanger"
	"math/big"
	"strings"
)

var (
	errOrderNotFound   = "orderNotFound"
	errInvalidFormat   = "invalidFormat"
	errApiKeysNotFound = "apiKeysNotFound"
)

type OrderRepo interface {
	CreateOrderWithSettings(ctx context.Context, order *domain.Order, settings *domain.Settings) (id int64, err error)
	CreateOrder(ctx context.Context, order *domain.Order) (id int64, err error)
	UpdateTpSl(ctx context.Context, id int64, price string, settings *domain.Settings) error
	CancelOrder(ctx context.Context, id int64, symbol, exchange string) error
	ExecuteOrder(ctx context.Context, id int64) error
	ActivateOrder(ctx context.Context, id int64) error

	GetOrder(ctx context.Context, orderId int64, symbol, exchange string) (*domain.Order, error)
	GetTpSlOrderByBaseOrder(ctx context.Context, id int64, symbol, exchange, tpsl string) (*domain.Order, error)
	GetTpSlOrdersByBaseOrder(ctx context.Context, id int64) ([]*domain.Order, error)
	GetOpposingTpSlOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetActiveTpSlOrders(ctx context.Context, exchange string) ([]*domain.Order, error)
	GetActiveSymbols(ctx context.Context) (domain.SymbolList, error)
	GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error)

	Atomic(ctx context.Context, fn func(ctx context.Context, orderRepo OrderRepo) error) (err error)
}

type Order struct {
	cfg               *config.Config
	exchanger         exchanger.Exchanger
	apiKeyRepo        ApiKeyRepo
	orderRepo         OrderRepo
	i18n              *i18n.I18n
	logger            *log.Logger
	binanceQueue      *domain.OrdersQueue
	limitBinanceQueue *domain.OrdersQueue
	keyDbCli          *redis.Client
}

func NewOrder(cfg *config.Config,
	exchanger exchanger.Exchanger,
	apiKeyRepo ApiKeyRepo,
	orderRepo OrderRepo,
	binanceQueue *domain.OrdersQueue,
	limitBinanceQueue *domain.OrdersQueue,
	keyDbCli *redis.Client,
	i18n *i18n.I18n,
	logger *log.Logger) *Order {
	return &Order{cfg: cfg,
		exchanger:         exchanger,
		apiKeyRepo:        apiKeyRepo,
		orderRepo:         orderRepo,
		binanceQueue:      binanceQueue,
		limitBinanceQueue: limitBinanceQueue,
		keyDbCli:          keyDbCli,
		i18n:              i18n,
		logger:            logger,
	}
}

func (o *Order) CreateOrder(ctx context.Context, orderDto *dto.Order, tgUserId int64) (order *domain.Order, err error) {
	var execOrderId int64
	var orderId int64
	var iceberg string

	keys, err := o.getApiKeys(ctx, tgUserId, orderDto.Exchange)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil && execOrderId != 0 && order.OrderType == consts.OrderTypeLimit {
			if err := o.exchanger.CancelOrder(keys,
				&dto.CancelOrder{Id: orderId,
					Exchange: orderDto.Exchange,
					Symbol:   orderDto.Symbol}); err != nil {
				o.logger.ErrorLog.Println("err cancel order in defer: ", err)
			}
		}
	}()

	if err = o.checkStopPercent(orderDto.StopPercent); orderDto.OrderType == consts.OrderTypeStopLossLimit && err != nil {
		return nil, err
	}

	if orderDto.OrderType == consts.OrderTypeStopLossLimit {
		price, err := o.getPrice(orderDto.StopPercent, orderDto.StopPrice, orderDto.Side)
		if err != nil {
			return nil, err
		}
		/*
			price, err = o.validatePrice(price, order.Symbol)
			if err != nil {
				return nil, err
			}
		*/
		orderDto.Price = price
	}

	if err := o.orderRepo.Atomic(ctx, func(ctx context.Context, orderRepo OrderRepo) error {
		execOrderId, err = o.exchanger.CreateOrder(keys, orderDto)
		if err != nil {
			o.logger.ErrorLog.Println("err create order: " + err.Error())

			return errors.BadRequestError(err.Error())
		}

		order = &domain.Order{
			ExecOrderId: execOrderId,
			Symbol:      orderDto.Symbol,
			Side:        strings.ToUpper(orderDto.Side),
			OrderType:   strings.ToUpper(orderDto.OrderType),
			Quantity:    orderDto.Quantity,
			Price:       orderDto.Price,
			TimeInForce: orderDto.TimeInForce,
			StopPrice:   orderDto.StopPrice,
			Exchange:    orderDto.Exchange,
			UserId:      tgUserId,
			Status:      consts.OrderStatusActive,
			TpSl:        consts.BaseOrderType,
		}
		if order.OrderType == consts.OrderTypeMarket {
			var price string
			if err := o.keyDbCli.Get(context.Background(),
				consts.TradePriceCacheKey+orderDto.Exchange+"_"+orderDto.Symbol).Scan(&price); err != nil {
				o.logger.ErrorLog.Println("err can get market price:", err)

				return errors.InternalServerError(err)
			}
			order.Price = price
			order.Status = consts.OrderStatusFilled
		}

		if iceberg != "" {
			order.IcebergQty = iceberg
		}

		var settings *domain.Settings
		if !orderDto.IsEmptyTpSl() {
			settings = &domain.Settings{
				TpPercent: orderDto.TpPercent,
				SlPercent: orderDto.SlPercent,
				TpPrice:   orderDto.TpPrice,
				SlPrice:   orderDto.SlPrice,
				Ts:        orderDto.Ts,
				TpType:    orderDto.TpType,
				SlType:    orderDto.SlType,
			}
		}

		orderId, err = orderRepo.CreateOrderWithSettings(ctx, order, settings)
		if err != nil {
			o.logger.ErrorLog.Println("Can`t save order:", err)

			return errors.InternalServerError(err)
		}
		order.Id = orderId

		if (settings.TpPercent != "" || settings.TpPrice != "") && order.TpSl == consts.BaseOrderType {
			if err = o.createTpSlOrder(ctx, order, consts.TpOrderType, settings); err != nil {
				return err
			}
		}

		if (settings.SlPercent != "" || settings.SlPrice != "") && order.TpSl == consts.BaseOrderType {
			if err = o.createTpSlOrder(ctx, order, consts.SlOrderType, settings); err != nil {
				return err
			}
		}
		if order.OrderType == consts.OrderTypeLimit {
			o.limitBinanceQueue.Add(order)
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return order, err
}

func (o *Order) UpdateTpslOrder(ctx context.Context, orderDto *dto.UpdateTpSl) error {
	updateTpsl := func(baseOrder *domain.Order,
		settings *domain.Settings,
		queue *domain.OrdersQueue,
		orderType string) error {
		tpSlOrder, err := o.orderRepo.GetTpSlOrderByBaseOrder(ctx,
			baseOrder.Id,
			baseOrder.Symbol,
			baseOrder.Exchange,
			orderType)
		if err != nil {
			o.logger.ErrorLog.Println("err get tp/sl order:", err)

			return errors.InternalServerError(err)
		}
		if tpSlOrder == nil {
			return errors.BadRequestError(o.i18n.T(errOrderNotFound, nil, "ru"))
		}
		price := o.getTpSlPrice(baseOrder, orderType, settings)

		if err := o.orderRepo.UpdateTpSl(ctx, tpSlOrder.Id, price, settings); err != nil {
			o.logger.ErrorLog.Println("err update tpsl order:", err)

			return errors.InternalServerError(err)
		}
		queue.UpdatePrice(tpSlOrder.Symbol, tpSlOrder.Id, price)

		return nil
	}

	settings := &domain.Settings{
		OrderId:   orderDto.Id,
		TpPercent: orderDto.TpPercent,
		SlPercent: orderDto.SlPercent,
		TpPrice:   orderDto.TpPrice,
		SlPrice:   orderDto.SlPrice,
	}

	baseOrder, err := o.orderRepo.GetOrder(ctx, orderDto.Id, orderDto.Symbol, orderDto.Exchange)
	if err != nil {
		o.logger.ErrorLog.Println("err get base order:", err)

		return errors.InternalServerError(err)
	}
	if baseOrder == nil {
		return errors.BadRequestError(o.i18n.T(errOrderNotFound, nil, "ru"))
	}

	queue := o.getExchangeQueue(baseOrder.Exchange)

	if !settings.IsTpEmpty() {
		if err = updateTpsl(baseOrder, settings, queue, consts.TpOrderType); err != nil {
			return err
		}
	}

	if !settings.IsSlEmpty() {
		if err = updateTpsl(baseOrder, settings, queue, consts.SlOrderType); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) CancelOrder(ctx context.Context, order *dto.CancelOrder, tgUserId int64) error {
	keys, err := o.getApiKeys(ctx, tgUserId, order.Exchange)
	if err != nil {
		return err
	}

	if err = o.exchanger.CancelOrder(keys, order); err != nil {
		o.logger.ErrorLog.Println("err cancel order: ", err)

		return errors.BadRequestError(err.Error())
	}

	if err = o.orderRepo.CancelOrder(ctx, order.Id, order.Symbol, order.Exchange); err != nil {
		o.logger.ErrorLog.Println("err cancel order in database: ", err)

		return errors.InternalServerError(err)
	}

	return nil
}

func (o *Order) GetActiveTpSlOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	result, err := o.orderRepo.GetActiveTpSlOrders(ctx, exchange)
	if err != nil {
		o.logger.ErrorLog.Println("err get active tpsl orders", err)

		return nil, errors.InternalServerError(err)
	}

	return result, nil
}

func (o *Order) GetUserActiveOrders(ctx context.Context, userId int64, exchange string) ([]domain.Order, error) {
	keys, err := o.getApiKeys(ctx, userId, exchange)
	if err != nil {
		return nil, err
	}

	result, err := o.exchanger.GetOpenOrders(keys, exchange, "")
	if err != nil {
		return nil, errors.BadRequestError(err.Error())
	}

	return result, nil
}

func (o *Order) ExecuteTpSlOrder(ctx context.Context, userId int64, order *domain.Order) (int64, error) {
	queue := o.getExchangeQueue(order.Exchange)
	if !queue.Exist(order.Symbol, order.Id) {
		return 0, nil
	}

	keys, err := o.getApiKeys(ctx, userId, order.Exchange)
	if err != nil {
		return 0, err
	}

	execOrder := new(dto.Order)
	conversion.FromDomainOrderToDtoOrder(order, execOrder)
	if execOrder.OrderType == consts.OrderTypeMarket {
		order.Status = consts.OrderStatusFilled
	}
	excId, err := o.exchanger.CreateOrder(keys, execOrder)
	if err != nil {
		o.logger.ErrorLog.Println("err execute order: " + err.Error())

		return 0, errors.BadRequestError(err.Error())
	}
	queue.Remove(order.Symbol, order.Id)

	o.orderRepo.Atomic(ctx, func(ctx context.Context, orderRepo OrderRepo) error {
		opposingOrder, err := orderRepo.GetOpposingTpSlOrder(ctx, order)
		if err != nil {
			o.logger.ErrorLog.Println("err get opposing order in db: " + err.Error())

			return errors.InternalServerError(err)
		}

		if opposingOrder != nil {
			if err := orderRepo.CancelOrder(ctx, opposingOrder.Id, opposingOrder.Symbol, opposingOrder.Exchange); err != nil {
				o.logger.ErrorLog.Println("err cancel opposing order in db: " + err.Error())

				return errors.InternalServerError(err)
			}
			if queue != nil {
				queue.Remove(opposingOrder.Symbol, opposingOrder.Id)
			}
		}

		order.ExecOrderId = excId
		if err := orderRepo.ExecuteOrder(ctx, order.Id); err != nil {
			o.logger.ErrorLog.Println("err exec order in db: " + err.Error())

			return errors.InternalServerError(err)
		}

		return nil
	})

	return excId, nil
}

func (o *Order) GetOrder(ctx context.Context, orderId int64, tgUserId int64, symbol, exchange string, inExchange bool) (*domain.Order, error) {
	if inExchange {
		keys, err := o.getApiKeys(ctx, tgUserId, exchange)
		if err != nil {
			return nil, err
		}
		res, err := o.exchanger.GetOrder(keys, exchange, symbol, orderId)
		if err != nil {
			o.logger.ErrorLog.Println("err get order: " + err.Error())

			return nil, errors.BadRequestError(err.Error())
		}
		return res, nil
	} else {
		res, err := o.orderRepo.GetOrder(ctx, orderId, exchange, symbol)
		if err != nil {
			o.logger.ErrorLog.Println("err get order: " + err.Error())

			return nil, errors.BadRequestError(err.Error())
		}
		return res, nil
	}
}

func (o *Order) SetFilledLimitOrder(ctx context.Context, order *domain.Order) error {
	limitQueue := o.getLimitExchangeQueue(order.Exchange)
	if !limitQueue.Exist(order.Symbol, order.Id) {
		o.logger.InfoLog.Println("not found in queue")

		return nil

	}
	limitQueue.Remove(order.Symbol, order.Id)
	err := o.orderRepo.Atomic(ctx, func(ctx context.Context, orderRepo OrderRepo) error {
		if err := orderRepo.ExecuteOrder(ctx, order.Id); err != nil {
			o.logger.ErrorLog.Println("err exec order in db: " + err.Error())

			return errors.InternalServerError(err)
		}

		tpSlOrders, err := orderRepo.GetTpSlOrdersByBaseOrder(ctx, order.Id)
		if err != nil {
			o.logger.ErrorLog.Println("err get tpsl orders in db: " + err.Error())

			return errors.InternalServerError(err)
		}

		tpSlQueue := o.getExchangeQueue(order.Exchange)
		for _, v := range tpSlOrders {
			if err = orderRepo.ActivateOrder(ctx, v.Id); err != nil {
				o.logger.ErrorLog.Println("err activate orders in db: " + err.Error())

				return errors.InternalServerError(err)
			}

			tpSlQueue.Add(v)
		}

		return nil
	})

	return err
}

func (o *Order) GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	orders, err := o.orderRepo.GetLimitOrders(ctx, exchange)
	if err != nil {
		o.logger.ErrorLog.Println("err get limit orders: " + err.Error())

		return nil, errors.BadRequestError(err.Error())
	}

	return orders, nil
}

func (o *Order) createTpSlOrder(ctx context.Context, order *domain.Order, tpSlType string, settings *domain.Settings) error {
	price := o.getTpSlPrice(order, tpSlType, settings)
	/*
		TODO::
			if err := o.validateInputParams(newPrice, order.Quantity, order.Symbol); err != nil {
				return nil, err
			}*/

	ordertype := consts.OrderTypeMarket
	if tpSlType == consts.SlOrderType && settings.SlType != "" {
		ordertype = settings.SlType
	} else if tpSlType == consts.TpOrderType && settings.TpType != "" {
		ordertype = settings.TpType
	}

	status := consts.OrderStatusActive
	if order.OrderType == consts.OrderTypeLimit {
		status = consts.OrderStatusTpSlInactive
	}

	tpSlOrder := &domain.Order{
		UserId:      order.UserId,
		Price:       price,
		Symbol:      strings.ToUpper(order.Symbol),
		Status:      status,
		OrderType:   ordertype,
		TimeInForce: consts.TimeInForceGTC,
		Quantity:    order.Quantity,
		Side:        o.getTpSlSide(order.Side),
		TpSl:        tpSlType,
		ExecOrderId: order.Id,
		Exchange:    order.Exchange,
	}

	newOrdId, err := o.orderRepo.CreateOrder(ctx, tpSlOrder)
	if err != nil {
		o.logger.ErrorLog.Println("Can`t save tpsl order:", err)

		return errors.InternalServerError(err)
	}
	tpSlOrder.Id = newOrdId

	if order.OrderType == consts.OrderTypeMarket {
		queue := o.getExchangeQueue(order.Exchange)
		if queue != nil {
			queue.Add(tpSlOrder)
		}
	}

	return nil
}

func (o *Order) checkStopPercent(stopPercent string) error {
	bigPer, ok := new(big.Float).SetString(stopPercent)
	if !ok {
		return errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
	}

	stopPer, _ := bigPer.Float64()
	if stopPer < 0 || stopPer > 100 {
		return errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
	}

	return nil
}

func (o *Order) getApiKeys(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	keys, err := o.apiKeyRepo.GetApiKeysByUserIdAndExchange(ctx, userId, exchange)
	if err != nil {
		o.logger.ErrorLog.Println("err get api keys: ", err)
	}
	if keys == nil {
		return nil, errors.BadRequestError(o.i18n.T(errApiKeysNotFound, nil, "ru"))
	}
	if keys.PrivKey != "" {
		keys.DecodePrivKey(o.cfg.GetApiSecret())
	}
	if keys.Passphrase != "" {
		keys.DecodePassKey(o.cfg.GetApiSecret())
	}

	return keys, nil
}

func (o *Order) getPrice(stopPercent string, stopPrice string, side string) (string, error) {
	bigPrice, ok := new(big.Float).SetString(stopPrice)
	if !ok {
		return "", errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
	}

	bigPercent, ok := new(big.Float).SetString(stopPercent)
	if !ok {
		return "", errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
	}

	percent, ok := new(big.Float).SetString("100")
	if !ok {
		return "", errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
	}

	result, ok := new(big.Float).SetString("0")
	if !ok {
		return "", errors.BadRequestError(o.i18n.T(errInvalidFormat, nil, "ru"))
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

func (o *Order) getTpSlPrice(order *domain.Order, tpSlType string, settings *domain.Settings) string {
	percentPrice := func(price, prct, ordSide string) string {
		ordPrice := helper.StringToBigFloat(price)
		percent := helper.StringToBigFloat(prct)
		totalPrice := new(big.Float)
		ordSide = strings.ToUpper(ordSide)
		tpSlType = strings.ToLower(tpSlType)

		if ordSide == consts.OrderSideBuy && tpSlType == consts.TpOrderType {
			totalPrice = helper.BigSumWithPercent(ordPrice, percent)
		} else if ordSide == consts.OrderSideSell && tpSlType == consts.TpOrderType {
			totalPrice = helper.BigDiffWithPercent(ordPrice, percent)
		} else if ordSide == consts.OrderSideBuy && tpSlType == consts.SlOrderType {
			totalPrice = helper.BigDiffWithPercent(ordPrice, percent)
		} else if ordSide == consts.OrderSideSell && tpSlType == consts.SlOrderType {
			totalPrice = helper.BigSumWithPercent(ordPrice, percent)
		}
		/*
			result, err := b.excInfo.GetPriceBySizeRules(totalPrice.String(), symbol)
			if err != nil {
				return err
			}*/

		return totalPrice.String()
	}

	if tpSlType == consts.TpOrderType {
		if settings.TpPrice == "" && settings.TpPercent != "" {
			return percentPrice(order.Price, settings.TpPercent, order.Side)
		} else {
			return settings.TpPrice
		}
	} else if tpSlType == consts.SlOrderType {
		if settings.SlPrice == "" && settings.SlPercent != "" {
			return percentPrice(order.Price, settings.SlPercent, order.Side)
		} else {
			return settings.SlPrice
		}
	}

	return ""
}

func (o *Order) getTpSlSide(side string) string {
	if side == consts.OrderSideSell {
		return consts.OrderSideBuy
	}

	return consts.OrderSideSell
}

func (o *Order) getExchangeQueue(exchange string) *domain.OrdersQueue {
	switch exchange {
	case consts.Binance:
		return o.binanceQueue
	default:
		return nil
	}
}

func (o *Order) getLimitExchangeQueue(exchange string) *domain.OrdersQueue {
	switch exchange {
	case consts.Binance:
		return o.limitBinanceQueue
	default:
		return nil
	}
}
