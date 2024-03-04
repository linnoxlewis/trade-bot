package heartbeat

import (
	"context"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/pkg/exchanger"
	"github.com/linnoxlewis/trade-bot/internal/pkg/telegram"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"os"
	"time"
)

type LimitOrderTicker struct {
	cfg              *config.Config
	orderSrv         telegram.OrderSrv
	exchanger        exchanger.Exchanger
	logger           *log.Logger
	ordersQueue      *domain.OrdersQueue
	limitOrdersQueue *domain.OrdersQueue
	ticker           *time.Ticker
	exchange         string
}

func NewLimitOrderTicker(cfg *config.Config,
	orderSrv telegram.OrderSrv,
	limitOrdersQueue *domain.OrdersQueue,
	logger *log.Logger,
	exchange string,
	heartbeatPeriod time.Duration) *LimitOrderTicker {
	result := &LimitOrderTicker{
		cfg:              cfg,
		ticker:           time.NewTicker(heartbeatPeriod),
		orderSrv:         orderSrv,
		logger:           logger,
		exchange:         exchange,
		limitOrdersQueue: limitOrdersQueue,
	}

	result.checkLimitQueue()

	return result
}

func (l *LimitOrderTicker) Tick(ctx context.Context, interrupt chan os.Signal) {
	l.logger.InfoLog.Printf("Start check limit orders in %s ", l.exchange)
	for {
		select {
		case <-interrupt:
			l.logger.InfoLog.Println("Limit orders Ticker stop")
		case <-ctx.Done():
			l.logger.InfoLog.Println("Limit orders Ticker stop")
			return
		case <-l.ticker.C:
			l.checkOrders(ctx)
		}
	}
}

func (l *LimitOrderTicker) checkOrders(ctx context.Context) {
	for _, val := range l.limitOrdersQueue.Orders {
		val := val
		for _, v := range val {
			if v.Inwork == true {
				continue
			}
			go l.checkOrder(ctx, v)
		}
	}
}

func (l *LimitOrderTicker) checkOrder(ctx context.Context, order *domain.Order) {
	order.Inwork = true
	order.Lock()
	defer func() {
		order.Inwork = false
		order.Unlock()
	}()

	if order.Status != consts.OrderStatusActive {
		return
	}

	excOrder, err := l.orderSrv.GetOrder(ctx, order.ExecOrderId,
		order.UserId,
		order.Symbol,
		order.Exchange,
		true)
	if err != nil {
		l.logger.ErrorLog.Println("err execute order:", order.Id, err)

		return
	}

	if excOrder.Status == consts.OrderStatusFilled {
		err = l.orderSrv.SetFilledLimitOrder(ctx, order)
		if err != nil {
			l.logger.ErrorLog.Println("err set status filled order:", order.Id, err)

			return
		}
	}
}

func (l *LimitOrderTicker) checkLimitQueue() {
	if l.limitOrdersQueue.Len() == 0 {
		orders, err := l.orderSrv.GetLimitOrders(context.Background(), l.exchange)
		if err != nil {
			l.logger.ErrorLog.Println("Error get Orders from db: ", err)

			return
		}
		for _, v := range orders {
			l.limitOrdersQueue.Add(v)
		}
	}
}
