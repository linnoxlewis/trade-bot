package validation

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"regexp"
)

var symbolRegexp = regexp.MustCompile("^[A-Z0-9]{1,20}$")
var intRegexp = regexp.MustCompile("^([0-9]{1,20})(\\.[0-9]{1,20})?$")
var intNegativeRegexp = regexp.MustCompile("^([-0-9]{1,20})(\\.[-0-9]{1,20})?$")
var errInvalidFormat = errors.New("err invalid format")

var minStopPricePer = 0

var timeInForce = []interface{}{
	consts.TimeInForceGTC,
	consts.TimeInForceIOK,
	consts.TimeInForceFOK,
}

type Validator struct{}

func NewValidator() Validator {
	return Validator{}
}

func (v Validator) ValidateAddApiKeys(inputData *dto.ApiKeys) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.PubKey, validation.Required.Error("empty_field")),
		validation.Field(&inputData.PrivKey, validation.Required.Error("empty_field")),
		validation.Field(&inputData.UserId, validation.Required.Error("empty_field")),
		validation.Field(&inputData.Exchange, validation.Required.Error("empty_field")))
}

func (v Validator) ValidateCancelOrder(inputData *dto.CancelOrder) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.Symbol, validation.Required,
			validation.Length(6, 10),
			validation.Match(symbolRegexp)),
		validation.Field(&inputData.Exchange, validation.Required),
		validation.Field(&inputData.Id, validation.Required,
			validation.Min(1)),
	)
}

func (v Validator) ValidateCreateOrder(inputData *dto.Order) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.Symbol, validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),
		validation.Field(&inputData.Exchange, validation.Required),
		validation.Field(&inputData.Side, validation.Required,
			validation.In(consts.OrderSideBuy, consts.OrderSideSell)),
		validation.Field(&inputData.OrderType, validation.Required,
			validation.In(consts.OrderTypeLimit, consts.OrderTypeMarket, consts.OrderTypeStopLossLimit)),
		validation.Field(&inputData.Price,
			validation.When(inputData.OrderType == consts.OrderTypeLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.Quantity, validation.Required,
			validation.Match(intRegexp),
			validation.By(zeroString),
		),
		validation.Field(&inputData.TimeInForce, validation.When(
			inputData.OrderType == consts.OrderTypeLimit, validation.Required),
			validation.In(timeInForce...),
		),
		validation.Field(&inputData.StopPercent, validation.When(
			inputData.OrderType == consts.OrderTypeStopLossLimit, validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&inputData.StopPrice,
			validation.When(inputData.OrderType == consts.OrderTypeStopLossLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.TpPercent,
			validation.When(inputData.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.SlPercent,
			validation.When(inputData.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.TpPrice,
			validation.When(inputData.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.SlPrice,
			validation.When(inputData.SlPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.Ts,
			validation.Match(intRegexp),
			validation.By(zeroString),
		),
	)
}

func (v Validator) Validate(inputData *dto.TgOrder) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.Command,
			validation.Required),

		validation.Field(&inputData.Symbol,
			validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),

		validation.Field(&inputData.Exchange,
			validation.Required),

		validation.Field(&inputData.Side,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.OrderType == consts.OrderTypeLimit,
				validation.Required,
				validation.In(consts.OrderSideBuy, consts.OrderSideSell))),

		validation.Field(&inputData.OrderType,
			validation.When(inputData.Command == consts.TgCreateOrderCommand, validation.Required,
				validation.In(consts.OrderTypeLimit, consts.OrderTypeMarket, consts.OrderTypeStopLossLimit))),

		validation.Field(&inputData.Price,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.OrderType == consts.OrderTypeLimit, validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString))),

		validation.Field(&inputData.Quantity,
			validation.When(inputData.Command == consts.TgCreateOrderCommand,
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString))),

		validation.Field(&inputData.TimeInForce,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.OrderType == consts.OrderTypeLimit,
				validation.Required),
			validation.In(timeInForce...),
		),

		validation.Field(&inputData.StopPercent, validation.When(
			inputData.Command == consts.TgCreateOrderCommand &&
				inputData.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),

		validation.Field(&inputData.StopPrice,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.OrderType == consts.OrderTypeStopLossLimit,
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&inputData.TpPercent,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&inputData.SlPercent,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&inputData.TpPrice,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&inputData.SlPrice,
			validation.When(inputData.Command == consts.TgCreateOrderCommand &&
				inputData.SlPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),

		validation.Field(&inputData.Ts,
			validation.Match(intRegexp),
			validation.By(zeroString),
		),
	)
}

func (v Validator) ValidateUpdateOrder(inputData *dto.UpdateOrder) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.Symbol, validation.Required,
			validation.Length(6, 10),
			validation.Match(symbolRegexp),
		),
		validation.Field(&inputData.Side, validation.Required,
			validation.In(consts.OrderSideBuy, consts.OrderSideSell),
		),
		validation.Field(&inputData.OrderType, validation.Required,
			validation.In(consts.OrderTypeLimit, consts.OrderTypeStopLossLimit),
		),
		validation.Field(&inputData.Price, validation.When(
			inputData.OrderType == consts.OrderTypeLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&inputData.Quantity, validation.Required,
			validation.Match(intRegexp),
		),
		validation.Field(&inputData.TimeInForce, validation.Required,
			validation.In(timeInForce...),
		),
		validation.Field(&inputData.OrderId, validation.Required,
			validation.Min(1),
		),
		validation.Field(&inputData.StopPercent, validation.When(
			inputData.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required,
			validation.Match(intRegexp),
		)),
		validation.Field(&inputData.StopPrice, validation.When(
			inputData.OrderType == consts.OrderTypeStopLossLimit,
			validation.Required),
			validation.Match(intRegexp),
		),
	)
}

func (v Validator) ValidateTpSl(inputData *dto.UpdateTpSl) error {
	return validation.ValidateStruct(inputData,
		validation.Field(&inputData.Symbol, validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),

		validation.Field(&inputData.Id, validation.Required),

		validation.Field(&inputData.Exchange, validation.Required),

		validation.Field(&inputData.TpPercent,
			validation.When(inputData.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.SlPercent,
			validation.When(inputData.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.TpPrice,
			validation.When(inputData.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&inputData.SlPrice,
			validation.When(inputData.SlPercent == "",
				validation.Required,
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
