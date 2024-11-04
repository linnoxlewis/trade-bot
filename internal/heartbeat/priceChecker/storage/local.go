package storage

import (
	"context"
	"fmt"
	"sync"
)

type LocalPriceStorage struct {
	mu     sync.RWMutex
	prices map[string]string
}

// NewLocalPriceStorage создает новое локальное хранилище.
func NewLocalPriceStorage() *LocalPriceStorage {
	return &LocalPriceStorage{
		prices: make(map[string]string),
	}
}

// SetPrice устанавливает цену в локальном хранилище.
func (l *LocalPriceStorage) SetPrice(ctx context.Context, pair string, price string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prices[pair] = price
	return nil
}

// GetPrice возвращает цену из локального хранилища.
func (l *LocalPriceStorage) GetPrice(ctx context.Context, pair string) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	price, exists := l.prices[pair]
	if !exists {
		return "0.0", fmt.Errorf("price for pair %s not found", pair)
	}
	return price, nil
}

// Close закрывает локальное хранилище (ничего не делает).
func (l *LocalPriceStorage) Close() error {
	return nil
}
