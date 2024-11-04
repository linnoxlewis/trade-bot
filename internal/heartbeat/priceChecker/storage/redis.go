package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisPriceStorage struct {
	client *redis.Client
}

// NewRedisPriceStorage создает новое хранилище на основе Redis.
func NewRedisPriceStorage(redisAddr string) (*RedisPriceStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisPriceStorage{client: client}, nil
}

// SetPrice устанавливает цену в Redis.
func (r *RedisPriceStorage) SetPrice(ctx context.Context, pair string, price string) error {
	return r.client.Set(ctx, pair, price, 0).Err()
}

// GetPrice возвращает цену из Redis.
func (r *RedisPriceStorage) GetPrice(ctx context.Context, pair string) (string, error) {
	price, err := r.client.Get(ctx, pair).Result()
	if errors.Is(err, redis.Nil) {
		return "0.0", fmt.Errorf("price for pair %s not found", pair)
	} else if err != nil {
		return "0.0", err
	}

	return price, nil
}

// Close закрывает соединение с Redis.
func (r *RedisPriceStorage) Close() error {
	return r.client.Close()
}
