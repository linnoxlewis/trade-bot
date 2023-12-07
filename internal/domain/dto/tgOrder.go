package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
)

type TgOrder struct {
	Id        int64  `json:"id"`
	Command   string `json:"command"`
	Exchange  string `json:"exchange"`
	Symbol    string `json:"ccy"`
	OrderType string `json:"type"`
	Side      string `json:"side"`
	Quantity  string `json:"qty"`
	Price     string `json:"price"`
	TpPercent string `json:"tp_percent"`
	SlPercent string `json:"sl_percent"`
	TpPrice   string `json:"tp_price"`
	SlPrice   string `json:"sl_price"`
	TpType    string `json:"tp_type"`
	SlType    string `json:"sl_type"`
	Ts        string `json:"ts"`

	TimeInForce string `json:"tif"`
	StopPercent string `json:"stopPercent"`
	StopPrice   string `json:"stopPrice"`
	IcebergQty  string `json:"icebergQty"`
}

func (o *TgOrder) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Command,
			validation.Required),

		validation.Field(&o.Symbol,
			validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),

		validation.Field(&o.Exchange,
			validation.Required),

		validation.Field(&o.Side,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.OrderType == consts.OrderTypeLimit,
				validation.Required,
				validation.In(consts.OrderSideBuy, consts.OrderSideSell))),

		validation.Field(&o.OrderType,
			validation.When(o.Command == consts.TgCreateOrderCommand, validation.Required,
				validation.In(consts.OrderTypeLimit, consts.OrderTypeMarket, consts.OrderTypeStopLossLimit))),

		validation.Field(&o.Price,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.OrderType == consts.OrderTypeLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString))),

		validation.Field(&o.Quantity,
			validation.When(o.Command == consts.TgCreateOrderCommand,
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString))),

		validation.Field(&o.TimeInForce,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.OrderType == consts.OrderTypeLimit,
				validation.Required),
			validation.In(timeInForce...),
		),

		validation.Field(&o.StopPercent, validation.When(
			o.Command == consts.TgCreateOrderCommand && o.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),

		validation.Field(&o.StopPrice,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.OrderType == consts.OrderTypeStopLossLimit,
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&o.TpPercent,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&o.SlPercent,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&o.TpPrice,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&o.SlPrice,
			validation.When(o.Command == consts.TgCreateOrderCommand && o.SlPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&o.Ts,
			validation.Match(intRegexp),
			validation.By(zeroString),
		),
	)
}
