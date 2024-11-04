package tpsl

import (
	"context"
	"fmt"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	redisClient *redis.Client
	stream      string
	group       string
	consumer    string
	blockTime   time.Duration
}

func NewRedisQueue(redisClient *redis.Client,
	stream, group, consumer string,
	blockTime time.Duration,
) *RedisQueue {
	return &RedisQueue{
		redisClient: redisClient,
		stream:      stream,
		group:       group,
		consumer:    consumer,
		blockTime:   blockTime,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context,
	order *domain.Order,
	priority int,
) error {
	orderData := map[string]interface{}{
		"orderId":   order.Id,
		"symbol":    order.Symbol,
		"userId":    order.UserId,
		"price":     order.Price,
		"side":      order.Side,
		"orderType": order.OrderType,
		"priority":  priority,
	}

	err := q.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: orderData,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to enqueue order: %w", err)
	}

	return nil
}

func (q *RedisQueue) Dequeue(ctx context.Context) (redis.XMessage, error) {
	streams, err := q.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    q.group,
		Consumer: q.consumer,
		Streams:  []string{q.stream, ">"},
		Count:    1,
		Block:    q.blockTime,
	}).Result()

	if err != nil {
		return redis.XMessage{}, fmt.Errorf("failed to read from stream: %w", err)
	}

	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return redis.XMessage{}, fmt.Errorf("no messages in stream")
	}

	return streams[0].Messages[0], nil
}

func (q *RedisQueue) Ack(ctx context.Context, messageID string) error {
	_, err := q.redisClient.XAck(ctx, q.stream, q.group, messageID).Result()
	if err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	return nil
}
