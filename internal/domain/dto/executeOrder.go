package dto

import (
	"strings"
)

type ExecuteOrder struct {
	UserId   int64
	OrderId  int64
	Symbol   string
	Exchange string `json:"exchange"`
}

func NewExecuteOrder(
	userId, orderId int64,
	exchange, symbol string,
) *ExecuteOrder {
	return &ExecuteOrder{
		Symbol:   strings.TrimSpace(strings.ToUpper(symbol)),
		OrderId:  orderId,
		UserId:   userId,
		Exchange: exchange,
	}
}
