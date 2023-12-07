package conversion

import (
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"strings"
)

func FromTgOrderToDtoOrder(dtoOrd *dto.TgOrder, ord *dto.Order) {
	ord.Symbol = dtoOrd.Symbol
	ord.OrderType = strings.ToUpper(dtoOrd.OrderType)
	ord.Exchange = strings.ToLower(dtoOrd.Exchange)
	ord.Price = dtoOrd.Price
	ord.Quantity = dtoOrd.Quantity
	ord.TimeInForce = strings.ToUpper(dtoOrd.TimeInForce)
	ord.Side = strings.ToUpper(dtoOrd.Side)
	ord.TpPercent = dtoOrd.TpPercent
	ord.SlPercent = dtoOrd.SlPercent
	ord.TpPrice = dtoOrd.TpPrice
	ord.SlPrice = dtoOrd.SlPrice
	ord.Ts = dtoOrd.Ts
	ord.StopPercent = dtoOrd.StopPercent
	ord.StopPrice = dtoOrd.StopPrice
	ord.IcebergQty = dtoOrd.StopPrice
	ord.TpType = strings.ToUpper(dtoOrd.TpType)
	ord.SlType = strings.ToUpper(dtoOrd.SlType)
}

func FromTgOrderToDtoCancelOrder(dtoOrd *dto.TgOrder, ord *dto.CancelOrder) {
	ord.Symbol = dtoOrd.Symbol
	ord.Exchange = strings.ToLower(dtoOrd.Exchange)
	ord.Id = dtoOrd.Id
}

func FromTgOrderToDtoUpdateTpSlOrder(dtoOrd *dto.TgOrder, ord *dto.UpdateTpSl) {
	ord.Symbol = dtoOrd.Symbol
	ord.Exchange = strings.ToLower(dtoOrd.Exchange)
	ord.Id = dtoOrd.Id
	ord.SlPercent = dtoOrd.SlPercent
	ord.TpPercent = dtoOrd.TpPercent
	ord.SlPrice = dtoOrd.SlPrice
	ord.TpPrice = dtoOrd.TpPrice
}

func FromDomainOrderToDtoOrder(dtoOrd *domain.Order, ord *dto.Order) {
	ord.Symbol = dtoOrd.Symbol
	ord.OrderType = dtoOrd.OrderType
	ord.Exchange = dtoOrd.Exchange
	ord.Price = dtoOrd.Price
	ord.Quantity = dtoOrd.Quantity
	ord.TimeInForce = dtoOrd.TimeInForce
	ord.Side = dtoOrd.Side
	ord.StopPrice = dtoOrd.StopPrice
	ord.IcebergQty = dtoOrd.StopPrice
}
