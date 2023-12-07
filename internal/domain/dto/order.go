package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"regexp"
	"strings"
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
	Command   string  `json:"command"`
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"ccy"`
	OrderType string  `json:"type"`
	Side      string  `json:"side"`
	Quantity  string  `json:"qty"`
	Price     string  `json:"price"`
	Tp        float32 `json:"tp"`
	Sl        float32 `json:"sl"`
	Ts        float32 `json:"ts"`

	TimeInForce string `json:"tif"`
	StopPercent string `json:"stopPercent"`
	StopPrice   string `json:"stopPrice"`
	IcebergQty  string `json:"icebergQty"`
}

func NewOrder(
	command string,
	exchange string,
	symbol string,
	side string,
	orderType string,
	quantity string,
	price string,
	tp, sl float32,
) *Order {
	return &Order{
		Command:   command,
		Exchange:  exchange,
		Symbol:    strings.ToUpper(symbol),
		Side:      side,
		OrderType: orderType,
		Quantity:  quantity,
		Price:     price,
		Tp:        tp,
		Sl:        sl,
	}
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
