package tpsl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"log"
	"strconv"
	"sync"
)

type OrderService interface {
	GetActiveOrders(ctx context.Context, symbol string) ([]*domain.Order, error)
	ExecuteOrder(ctx context.Context, order *domain.Order) error
	LockOrder(ctx context.Context, orderID int64) (*domain.Order, error)
}

type PriceService interface {
	GetPrice(symbol string) (float64, error)
}

type NotificationService interface {
	SendNotification(userID int64, message string) error
}

type OrderWorker struct {
	orderService        OrderService
	priceService        PriceService
	notificationService NotificationService
	queue               *RedisQueue
	mu                  sync.Mutex
}

func NewOrderWorker(orderService OrderService,
	priceService PriceService,
	notificationService NotificationService,
	queue *RedisQueue) *OrderWorker {
	return &OrderWorker{
		orderService:        orderService,
		priceService:        priceService,
		notificationService: notificationService,
		queue:               queue,
	}
}

// Start Запуск воркера в отдельной горутине
func (w *OrderWorker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping worker")
			return
		default:
			orderMsg, err := w.queue.Dequeue(ctx)
			if err == nil {
				order, err := w.parseOrderFromMessage(orderMsg)
				if err != nil {
					log.Println("Error parsing order from message:", err)
					continue
				}

				w.processOrder(ctx, order)
			}
		}
	}
}

// Обработка ордера
func (w *OrderWorker) processOrder(ctx context.Context, order *domain.Order) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Получаем текущую цену актива
	price, err := w.priceService.GetPrice(order.Symbol)
	if err != nil {
		log.Println("Error fetching price:", err)
		return
	}

	lockedOrder, err := w.orderService.LockOrder(ctx, order.Id)
	if err != nil {
		log.Println("Error locking order:", err)
		return
	}

	// Преобразуем строковые цены ордера в float64 для корректного сравнения
	orderPrice, err := strconv.ParseFloat(lockedOrder.Price, 64)
	if err != nil {
		log.Printf("Error parsing order price for order ID %d: %v", lockedOrder.Id, err)
		return
	}

	stopPrice, err := strconv.ParseFloat(lockedOrder.StopPrice, 64)
	if err != nil {
		log.Printf("Error parsing stop price for order ID %d: %v", lockedOrder.Id, err)
		return
	}

	// Логика исполнения TP или SL
	if lockedOrder.Side == consts.OrderSideBuy && price >= orderPrice {
		w.executeOrder(ctx, lockedOrder, "Take Profit")
	} else if lockedOrder.Side == consts.OrderSideSell && price <= orderPrice {
		w.executeOrder(ctx, lockedOrder, "Take Profit")
	} else if lockedOrder.Side == consts.OrderSideBuy && price <= stopPrice {
		w.executeOrder(ctx, lockedOrder, "Stop Loss")
	} else if lockedOrder.Side == consts.OrderSideSell && price >= stopPrice {
		w.executeOrder(ctx, lockedOrder, "Stop Loss")
	} else {
		log.Println("Order condition not met, re-enqueueing...")

		// Повторно добавляем в очередь
		err := w.queue.Enqueue(ctx, lockedOrder, 1)
		if err != nil {
			log.Printf("Error re-enqueuing order ID %d: %v", lockedOrder.Id, err)
		}
	}
}

// Исполнение ордера
func (w *OrderWorker) executeOrder(ctx context.Context, order *domain.Order, reason string) {
	err := w.orderService.ExecuteOrder(ctx, order)
	if err != nil {
		log.Println("Error executing order:", order.Id, err)

		return
	}

	msg := fmt.Sprintf("%s executed for order %d at price %s", reason, order.Id, order.Price)
	w.notificationService.SendNotification(order.UserId, msg)

	log.Printf("Order %d executed: %s\n", order.Id, reason)
}

// Функция для извлечения данных ордера из сообщения Redis Streams
func (w *OrderWorker) parseOrderFromMessage(message redis.XMessage) (*domain.Order, error) {
	orderID, err := strconv.ParseInt(message.Values["orderId"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse orderId: %w", err)
	}

	userID, err := strconv.ParseInt(message.Values["userId"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userId: %w", err)
	}

	price := message.Values["price"].(string)
	symbol := message.Values["symbol"].(string)
	side := message.Values["side"].(string)
	orderType := message.Values["orderType"].(string)

	order := &domain.Order{
		Id:        orderID,
		UserId:    userID,
		Price:     price,
		Symbol:    symbol,
		Side:      side,
		OrderType: orderType,
	}

	return order, nil
}
