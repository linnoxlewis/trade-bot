package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"regexp"
)

var symbolRegexp = regexp.MustCompile("^[A-Z0-9]{1,20}$")
var intRegexp = regexp.MustCompile("^([0-9]{1,20})(\\.[0-9]{1,20})?$")
var intNegativeRegexp = regexp.MustCompile("^([-0-9]{1,20})(\\.[-0-9]{1,20})?$")
var errInvalidFormat = errors.New("err invalid format")

var minStopPricePer = 0
var maxStopPricePer = 100.1

var timeInForce = []interface{}{
	consts.TimeInForceGTC,
	consts.TimeInForceIOK,
	consts.TimeInForceFOK,
}

type Order struct {
	Exchange    string `json:"exchange"`
	Symbol      string `json:"ccy"`
	OrderType   string `json:"type"`
	Side        string `json:"side"`
	Quantity    string `json:"qty"`
	Price       string `json:"price"`
	TpPercent   string `json:"tp_percent"`
	SlPercent   string `json:"sl_percent"`
	TpPrice     string `json:"tp_price"`
	SlPrice     string `json:"sl_price"`
	TpType      string `json:"tp_type"`
	SlType      string `json:"sl_type"`
	Ts          string `json:"ts"`
	TimeInForce string `json:"tif"`
	StopPercent string `json:"stopPercent"`
	StopPrice   string `json:"stopPrice"`
	IcebergQty  string `json:"icebergQty"`
}

func (o *Order) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Symbol, validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),
		validation.Field(&o.Exchange, validation.Required),
		validation.Field(&o.Side, validation.Required,
			validation.In(consts.OrderSideBuy, consts.OrderSideSell)),
		validation.Field(&o.OrderType, validation.Required,
			validation.In(consts.OrderTypeLimit, consts.OrderTypeMarket, consts.OrderTypeStopLossLimit)),
		validation.Field(&o.Price,
			validation.When(o.OrderType == consts.OrderTypeLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&o.Quantity, validation.Required,
			validation.Match(intRegexp),
			validation.By(zeroString),
		),
		validation.Field(&o.TimeInForce, validation.When(
			o.OrderType == consts.OrderTypeLimit, validation.Required),
			validation.In(timeInForce...),
		),
		validation.Field(&o.StopPercent, validation.When(
			o.OrderType == consts.OrderTypeStopLossLimit, validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&o.StopPrice,
			validation.When(o.OrderType == consts.OrderTypeStopLossLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&o.TpPercent,
			validation.When(o.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&o.SlPercent,
			validation.When(o.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&o.TpPrice,
			validation.When(o.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&o.SlPrice,
			validation.When(o.SlPercent == "",
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

func zeroString(value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return errInvalidFormat
	}

	matchString, err := regexp.MatchString("^[0]+$", val)
	if err != nil || matchString {
		return errInvalidFormat
	}

	return nil
}

func (o *Order) IsEmptyTpSl() bool {
	return o.TpPercent == "" &&
		o.SlPercent == "" &&
		o.TpPrice == "" &&
		o.SlPrice == ""
}
