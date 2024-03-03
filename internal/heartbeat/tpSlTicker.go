package heartbeat

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	ordSrv "github.com/linnoxlewis/trade-bot/internal/pkg/telegram"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/linnoxlewis/trade-bot/pkg/telegram"
	"math/big"
	"os"
	"time"
)

var (
	successExecuted = "successExecuted"
	executeFailed   = "executeFailed"
)

type TpSlTicker struct {
	cfg         *config.Config
	orderSrv    ordSrv.OrderSrv
	keyDbCli    *redis.Client
	telegramCli *telegram.Client
	logger      *log.Logger
	ordersQueue *domain.OrdersQueue
	ticker      *time.Ticker
	i18n        *i18n.I18n
	exchange    string
	debugMode   bool
}

func NewTpSlTicker(cfg *config.Config,
	orderSrv ordSrv.OrderSrv,
	cache *redis.Client,
	logger *log.Logger,
	telegramCli *telegram.Client,
	heartbeatPeriod time.Duration,
	ordersQueue *domain.OrdersQueue,
	exchange string,
	i18n *i18n.I18n,
	debugMode bool,
) *TpSlTicker {
	result := &TpSlTicker{
		cfg:         cfg,
		ticker:      time.NewTicker(heartbeatPeriod),
		ordersQueue: ordersQueue,
		orderSrv:    orderSrv,
		logger:      logger,
		keyDbCli:    cache,
		exchange:    exchange,
		telegramCli: telegramCli,
		i18n:        i18n,
		debugMode:   debugMode,
	}

	result.checkQueue()

	return result
}

func (t *TpSlTicker) Tick(ctx context.Context, interrupt chan os.Signal) {
	t.logger.InfoLog.Printf("Start check tpsl in %s ", t.exchange)
	for {
		select {
		case <-interrupt:
			t.logger.InfoLog.Println("Tp/Sl Ticker stop")

		case <-ctx.Done():
			t.logger.InfoLog.Println("Tp/Sl Ticker stop")
			return
		case <-t.ticker.C:
			t.checkOrders()
		}
	}
}

func (t *TpSlTicker) checkOrders() {
	for _, val := range t.ordersQueue.Orders {
		val := val
		for _, v := range val {
			if v.Inwork == true {
				continue
			}
			go t.checkTpSl(v)
		}
	}
}

func (t *TpSlTicker) checkTpSl(order *domain.Order) {
	order.Inwork = true
	order.Lock()
	defer func() {
		order.Inwork = false
		order.Unlock()
	}()
	if order.Status != consts.OrderStatusActive {
		return
	}
	price := t.getPriceFromCache(order.Symbol)
	if t.debugMode {
		t.logger.InfoLog.Printf("current trade price %s", price)
		t.logger.InfoLog.Printf("start checking %s %s order %d with price %s ", order.Side, order.TpSl, order.Id, order.Price)
	}
	currentPrice := helper.StringToBigFloat(order.Price)
	tradePrice := helper.StringToBigFloat(price)

	if tradePrice == nil {
		t.logger.ErrorLog.Println("TRADE PRICE NIL")

		return
	}

	if order.TpSl == consts.TpOrderType {
		t.checkTakeProfit(order, currentPrice, tradePrice)
	}
	if order.TpSl == consts.SlOrderType {
		t.checkStopLoss(order, currentPrice, tradePrice)
	}
}

func (t *TpSlTicker) checkTakeProfit(order *domain.Order, virtualPrice, tradePrice *big.Float) {
	if t.debugMode {
		t.logger.InfoLog.Println("----------------------------------------")
		t.logger.InfoLog.Println("tradePrice TP:", tradePrice, "  ")
		t.logger.InfoLog.Println("order price TP:", order.Price, " ")
		t.logger.InfoLog.Println("cmp TP :", tradePrice.Cmp(virtualPrice), " ")
		t.logger.InfoLog.Println("----------------------------------------")
	}

	if (order.Side == consts.OrderSideBuy && (tradePrice.Cmp(virtualPrice) == 0 ||
		tradePrice.Cmp(virtualPrice) == -1)) ||
		(order.Side == consts.OrderSideSell && (tradePrice.Cmp(virtualPrice) == 1 ||
			tradePrice.Cmp(virtualPrice) == 0)) {
		if t.debugMode {
			t.logger.InfoLog.Println("execute TP order:", order.Id,
				tradePrice,
				order.Side,
				virtualPrice,
				tradePrice.Cmp(virtualPrice))
		}
		order.Price = tradePrice.String()
		id, err := t.orderSrv.ExecuteTpSlOrder(context.Background(), order.UserId, order)
		if err != nil {
			t.logger.ErrorLog.Println("err execute order:", order.Id, err)

			t.telegramCli.SendMessage(context.Background(),
				int(order.UserId),
				t.i18n.T(executeFailed, map[string]interface{}{
					"TpSl":  order.TpSl,
					"Side":  order.Side,
					"Price": order.Price,
				}, "ru")+":"+err.Error(), "")

			return
		}

		if err := t.telegramCli.SendMessage(context.Background(), int(order.UserId),
			t.i18n.T(successExecuted, map[string]interface{}{
				"TpSl":     order.TpSl,
				"Side":     order.Side,
				"Price":    order.Price,
				"ExecId":   id,
				"Id":       order.Id,
				"Exchange": order.Exchange,
				"Symbol":   order.Symbol,
				"Quantity": order.Quantity,
			}, "ru"), ""); err != nil {
			t.logger.ErrorLog.Println("cant`t send tg message: ", err)
		}

		return
	}
}

func (t *TpSlTicker) checkStopLoss(order *domain.Order, virtualPrice, tradePrice *big.Float) {
	if t.debugMode {
		t.logger.InfoLog.Println("----------------------------------------")
		t.logger.InfoLog.Println("tradePrice SL:", tradePrice, "  ")
		t.logger.InfoLog.Println("order price SL:", order.Price, " ")
		t.logger.InfoLog.Println("cmp SL:", tradePrice.Cmp(virtualPrice), " ")
		t.logger.InfoLog.Println("----------------------------------------")
	}
	if (order.Side == consts.OrderSideSell && (tradePrice.Cmp(virtualPrice) == -1 || tradePrice.Cmp(virtualPrice) == 0)) ||
		(order.Side == consts.OrderSideBuy && (tradePrice.Cmp(virtualPrice) == +1 || tradePrice.Cmp(virtualPrice) == 0)) {
		if t.debugMode {
			t.logger.InfoLog.Println("----------------------------------------")
			t.logger.InfoLog.Println("order id::", order.Id, "  ")
			t.logger.InfoLog.Println("tradePrice SL:", tradePrice, "  ")
			t.logger.InfoLog.Println("order price SL:", order.Price, " ")
			t.logger.InfoLog.Println("cmp SL:", tradePrice.Cmp(virtualPrice), " ")
			t.logger.InfoLog.Println("----------------------------------------")
		}
		if t.debugMode {
			t.logger.InfoLog.Println("execute Sl deal:", order.Id,
				tradePrice,
				order.Side,
				virtualPrice,
				tradePrice.Cmp(virtualPrice))
		}
		order.Price = tradePrice.String()
		id, err := t.orderSrv.ExecuteTpSlOrder(context.Background(), order.UserId, order)
		if err != nil {
			t.logger.InfoLog.Println("err execute order:", order.Id, err)
			t.telegramCli.SendMessage(context.Background(),
				int(order.UserId),
				t.i18n.T(executeFailed, map[string]interface{}{
					"TpSl":  order.TpSl,
					"Side":  order.Side,
					"Price": order.Price,
				}, "ru")+":"+err.Error(), "")

			return
		}

		if err := t.telegramCli.SendMessage(context.Background(), int(order.UserId),
			t.i18n.T(successExecuted, map[string]interface{}{
				"TpSl":        order.TpSl,
				"Side":        order.Side,
				"Price":       order.Price,
				"ExecOrderId": id,
				"Id":          order.Id,
				"Exchange":    order.Exchange,
				"Symbol":      order.Symbol,
				"Quantity":    order.Quantity,
			}, "ru"), "",
		); err != nil {
			t.logger.ErrorLog.Println("cant`t send tg message: ", err)
		}

		return
	}
}

func (t *TpSlTicker) checkQueue() {
	if t.ordersQueue.Len() == 0 {
		orders, err := t.orderSrv.GetActiveTpSlOrders(context.Background(), t.exchange)
		if err != nil {
			t.logger.ErrorLog.Println("Error get Orders from db: ", err)
		}
		for _, v := range orders {
			t.ordersQueue.Add(v)
		}
	}
}

func (t *TpSlTicker) getPriceFromCache(symbol string) string {
	var price string
	if err := t.keyDbCli.Get(context.Background(),
		consts.TradePriceCacheKey+t.exchange+"_"+symbol).Scan(&price); err != nil {
		return ""
	}

	return price
}
