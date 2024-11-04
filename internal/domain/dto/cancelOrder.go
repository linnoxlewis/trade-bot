package dto

import (
	"strings"
)

type CancelOrder struct {
	Id       int64  `json:"id"`
	Symbol   string `json:"ccy"`
	Exchange string `json:"exchange"`
}

func NewCancelOrder(orderId int64,
	symbol, exchange string) *CancelOrder {
	return &CancelOrder{
		Symbol:   strings.TrimSpace(strings.ToUpper(symbol)),
		Id:       orderId,
		Exchange: exchange,
	}
}
