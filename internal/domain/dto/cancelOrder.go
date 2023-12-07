package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type CancelOrder struct {
	Id       int64  `json:"id"`
	Symbol   string `json:"ccy"`
	Exchange string `json:"exchange"`
}

func NewCancelOrder(
	orderId int64,
	symbol,
	exchange string) *CancelOrder {
	return &CancelOrder{
		Symbol:   strings.TrimSpace(strings.ToUpper(symbol)),
		Id:       orderId,
		Exchange: exchange,
	}
}

func (o *CancelOrder) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Symbol, validation.Required,
			validation.Length(6, 10),
			validation.Match(symbolRegexp)),
		validation.Field(&o.Exchange, validation.Required),
		validation.Field(&o.Id, validation.Required,
			validation.Min(1)),
	)
}
