package excahnger

import (
	"context"
	binanceCli "github.com/adshao/go-binance/v2"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
)

type BinanceAdapter struct{}

func NewBinanceAdapter() *BinanceAdapter {
	return &BinanceAdapter{}
}

func (b *BinanceAdapter) CreateOrder(pubKey, secKey, passPhrase string, order *dto.Order) (int64, error) {
	result, err := binanceCli.NewClient(pubKey, secKey).
		NewCreateOrderService().
		Symbol(order.Symbol).
		Side(binanceCli.SideType(order.Side)).
		Type(binanceCli.OrderType(order.OrderType)).
		TimeInForce(binanceCli.TimeInForceType(order.TimeInForce)).
		Quantity(order.Quantity).
		Price(order.Price).Do(context.Background())
	if err != nil {
		return 0, err
	}

	return result.OrderID, nil
}

func (b *BinanceAdapter) CancelOrder(pubKey, secKey, passPhrase string, order *dto.CancelOrder) error {
	_, err := binanceCli.NewClient(pubKey, secKey).
		NewCancelOrderService().Symbol(order.Symbol).
		OrderID(order.OrderId).Do(context.Background())

	return err
}
