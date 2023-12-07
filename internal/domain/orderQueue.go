package domain

import (
	"sync"
)

type OrdersQueue struct {
	Exchange string
	Orders   map[string]map[int64]*Order
	sync.RWMutex
}

func NewOrderQueue(exchange string) *OrdersQueue {
	return &OrdersQueue{
		Exchange: exchange,
		Orders:   make(map[string]map[int64]*Order),
	}
}

func (ol *OrdersQueue) Add(order *Order) bool {
	ol.Lock()
	defer ol.Unlock()
	_, ok := ol.Orders[order.Symbol]
	if !ok {
		ol.Orders[order.Symbol] = make(map[int64]*Order)
	}

	ol.Orders[order.Symbol][order.Id] = order

	return true
}

func (ol *OrdersQueue) Remove(symbol string, id int64) bool {
	_, ok := ol.Orders[symbol]
	if !ok {
		return false
	}
	ol.Lock()
	defer ol.Unlock()
	delete(ol.Orders[symbol], id)
	return true
}

func (ol *OrdersQueue) UpdatePrice(symbol string, id int64, price string) bool {
	_, ok := ol.Orders[symbol]
	if !ok {
		return false
	}
	ol.Lock()
	defer ol.Unlock()
	order, ok := ol.Orders[symbol][id]
	if !ok {
		return false
	}
	order.Price = price
	delete(ol.Orders[symbol], id)
	ol.Orders[symbol][id] = order

	return true
}

func (ol *OrdersQueue) Len() int {
	return len(ol.Orders)
}

func (ol *OrdersQueue) Exist(symbol string, id int64) bool {
	val, ok := ol.Orders[symbol]
	if !ok {
		return false
	}
	_, ok = val[id]

	return ok
}
