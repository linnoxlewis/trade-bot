package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type CancelOrder struct {
	OrderId  int64
	Symbol   string
	Exchange string `json:"exchange"`
}

func NewCancelOrder(
	orderId int64,
	symbol string,
) *CancelOrder {
	return &CancelOrder{
		Symbol:  strings.TrimSpace(strings.ToUpper(symbol)),
		OrderId: orderId,
	}
}

func (o *CancelOrder) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Symbol, validation.Required,
			validation.Length(6, 10),
			validation.Match(symbolRegexp)),
		validation.Field(&o.Exchange, validation.Required),
		validation.Field(&o.OrderId, validation.Required,
			validation.Min(1)),
	)
}
