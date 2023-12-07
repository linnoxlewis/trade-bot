package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
)

type UpdateOrder struct {
	OrderId     int64   `json:"orderId"`
	Exchange    string  `json:"exchange"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	OrderType   string  `json:"orderType"`
	Quantity    string  `json:"quantity"`
	Price       string  `json:"price"`
	TimeInForce string  `json:"timeInForce"`
	StopPrice   string  `json:"stopPrice"`
	StopPercent string  `json:"stopPercent"`
	Tp          float32 `json:"tp"`
	Sl          float32 `json:"sl"`
	Ts          float32 `json:"ts"`
}

func NewUpdateOrder(
	orderId int64,
	symbol,
	side,
	orderType,
	quantity,
	price,
	timeInForce,
	stopPrice,
	stopPercent string,
	tp, sl, ts float32,
) *UpdateOrder {
	return &UpdateOrder{
		OrderId:     orderId,
		Symbol:      symbol,
		Side:        side,
		OrderType:   orderType,
		Quantity:    quantity,
		Price:       price,
		TimeInForce: timeInForce,
		StopPrice:   stopPrice,
		StopPercent: stopPercent,
		Tp:          tp,
		Sl:          sl,
		Ts:          ts,
	}
}

func (o *UpdateOrder) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Symbol, validation.Required,
			validation.Length(6, 10),
			validation.Match(symbolRegexp),
		),
		validation.Field(&o.Side, validation.Required,
			validation.In(consts.OrderSideBuy, consts.OrderSideSell),
		),
		validation.Field(&o.OrderType, validation.Required,
			validation.In(consts.OrderTypeLimit, consts.OrderTypeStopLossLimit),
		),
		validation.Field(&o.Price, validation.When(
			o.OrderType == consts.OrderTypeLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&o.Quantity, validation.Required,
			validation.Match(intRegexp),
		),
		validation.Field(&o.TimeInForce, validation.Required,
			validation.In(timeInForce...),
		),
		validation.Field(&o.OrderId, validation.Required,
			validation.Min(1),
		),
		validation.Field(&o.StopPercent, validation.When(
			o.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&o.StopPrice, validation.When(
			o.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required),
			validation.Match(intRegexp),
		),
	)
}
