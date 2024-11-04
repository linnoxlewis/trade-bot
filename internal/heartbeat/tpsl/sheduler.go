package tpsl

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"log"
	"math"
	"strconv"
)

type DynamicScheduler struct {
	redisClient  *redis.Client
	orderService OrderService
	queue        *RedisQueue
}

func NewDynamicScheduler(redisClient *redis.Client,
	orderService OrderService,
	queue *RedisQueue) *DynamicScheduler {
	return &DynamicScheduler{
		redisClient:  redisClient,
		orderService: orderService,
		queue:        queue,
	}
}

// Функция для планирования ордеров
func (s *DynamicScheduler) scheduleOrders(ctx context.Context, symbol string) {
	// Получаем активные ордера из базы данных
	orders, err := s.orderService.GetActiveOrders(ctx, symbol)
	if err != nil {
		log.Println("Error fetching active orders:", err)

		return
	}

	// Получение текущей цены актива из кеша
	currentPrice, err := s.getPriceFromCache(symbol)
	if err != nil {
		log.Println("Error getting current price:", err)
		return
	}

	// Расчет приоритета для каждого ордера
	for _, order := range orders {
		priority := s.calculatePriority(order, currentPrice)
		log.Printf("Order ID: %d, Priority: %d\n", order.Id, priority)

		// Добавляем ордер в Redis Streams с учетом приоритета
		err := s.queue.Enqueue(ctx, order, priority)
		if err != nil {
			log.Printf("Failed to enqueue order ID %d: %v\n", order.Id, err)
		} else {
			log.Printf("Order ID %d enqueued with priority %d\n", order.Id, priority)
		}
	}
}

// Функция для расчета приоритета ордера
func (s *DynamicScheduler) calculatePriority(order *domain.Order,
	currentPrice float64) int {
	// Преобразуем TP и SL цены к float
	tpPrice, _ := strconv.ParseFloat(order.Settings.TpPrice, 64)
	slPrice, _ := strconv.ParseFloat(order.Settings.SlPrice, 64)

	// Рассчитываем разницу между текущей ценой и TP/SL
	tpDiff := math.Abs(currentPrice - tpPrice)
	slDiff := math.Abs(currentPrice - slPrice)

	// Используем минимальную разницу для приоритета
	minDiff := math.Min(tpDiff, slDiff)

	// Чем меньше разница, тем выше приоритет
	priority := int((1 / minDiff) * 1000)

	// Если цены слишком далеки, приоритет будет низким
	if math.IsInf(float64(priority), 0) || math.IsNaN(float64(priority)) {
		priority = 0
	}

	return priority
}

// Получение текущей цены из кеша
func (s *DynamicScheduler) getPriceFromCache(symbol string) (float64, error) {
	priceStr, err := s.redisClient.Get(context.Background(), "price:"+symbol).Result()
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
