package exchanger

import (
	"context"
	"errors"
	binanceCli "github.com/adshao/go-binance/v2"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
	"strings"
	"sync"
)

type BinanceAdapter struct {
	useTestnet bool
}

func NewBinanceAdapter(useTestnet bool) *BinanceAdapter {
	return &BinanceAdapter{
		useTestnet: useTestnet,
	}
}

func (b *BinanceAdapter) CreateOrder(pubKey, secKey, passPhrase string, order *dto.Order) (int64, error) {
	binanceCli.UseTestnet = b.useTestnet
	ordType := binanceCli.OrderType(strings.ToUpper(order.OrderType))
	cli := binanceCli.NewClient(pubKey, secKey).
		NewCreateOrderService().
		Symbol(strings.ToUpper(order.Symbol)).
		Side(binanceCli.SideType(strings.ToUpper(order.Side))).
		Type(ordType).
		Quantity(order.Quantity)

	if ordType != binanceCli.OrderTypeMarket {
		cli.Price(order.Price)
		cli.TimeInForce(binanceCli.TimeInForceTypeGTC)
	}

	result, err := cli.Do(context.Background())
	if err != nil {
		return 0, err
	}

	return result.OrderID, nil
}

func (b *BinanceAdapter) CancelOrder(pubKey, secKey, passPhrase string, order *dto.CancelOrder) error {
	binanceCli.UseTestnet = b.useTestnet
	_, err := binanceCli.NewClient(pubKey, secKey).
		NewCancelOrderService().Symbol(order.Symbol).
		OrderID(order.Id).Do(context.Background())

	return err
}

func (b *BinanceAdapter) UpdateOrder(pubKey, secKey, passPhrase string, order *dto.UpdateOrder) (int64, error) {
	binanceCli.UseTestnet = b.useTestnet
	if err := b.CancelOrder(pubKey, secKey, passPhrase, &dto.CancelOrder{
		Id:       order.OrderId,
		Symbol:   order.Symbol,
		Exchange: order.Exchange,
	}); err != nil {
		return 0, err
	}

	return b.CreateOrder(pubKey, secKey, passPhrase, &dto.Order{
		Symbol:      order.Symbol,
		Side:        order.Side,
		OrderType:   order.OrderType,
		TimeInForce: order.TimeInForce,
		Quantity:    order.Quantity,
		Price:       order.Price,
	})
}

func (b *BinanceAdapter) GetBalance(pubKey, secKey, passPhrase string) (domain.Balance, error) {
	binanceCli.UseTestnet = b.useTestnet
	result, err := binanceCli.NewClient(pubKey, secKey).
		NewGetAccountService().Do(context.Background())

	if err != nil {
		return nil, err
	}
	var balance domain.Balance

	if result.Balances == nil {
		return nil, errors.New("empty balances")
	}
	for _, v := range result.Balances {
		balSymb := domain.NewBalanceSymbol(v.Asset, v.Free)
		balance = append(balance, balSymb)
	}

	return balance, nil
}

func (b *BinanceAdapter) GetOpenOrders(pubKey, secKey, passPhrase, symbol string) ([]domain.Order, error) {
	binanceCli.UseTestnet = b.useTestnet
	cli := binanceCli.NewClient(pubKey, secKey).NewListOpenOrdersService()
	if symbol != "" {
		cli.Symbol(symbol)
	}
	openOrders, err := cli.Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(openOrders) == 0 {
		return nil, nil
	}

	orderList := make([]domain.Order, len(openOrders))

	for i := 0; i < len(openOrders); i++ {
		var order domain.Order

		/*if openOrders[i].ClientOrderID != "" {
			id, _ := strconv.Atoi(openOrders[i].ClientOrderID)
			order.Id = int64(id)
		}*/

		order.ExecOrderId = openOrders[i].OrderID
		order.OrderType = string(openOrders[i].Type)
		order.Status = b.getStatus(string(openOrders[i].Status))
		order.Quantity = openOrders[i].OrigQuantity
		order.Symbol = openOrders[i].Symbol
		order.Exchange = BinanceType
		order.Price = openOrders[i].Price
		order.Side = string(openOrders[i].Side)
		order.TimeInForce = string(openOrders[i].TimeInForce)
		order.StopPrice = openOrders[i].StopPrice
		order.IcebergQty = openOrders[i].IcebergQuantity
		orderList[i] = order
	}

	return orderList, nil
}

func (b *BinanceAdapter) GetOrder(pubKey, secKey, passPhrase, symbol string, orderId int64) (*domain.Order, error) {
	binanceCli.UseTestnet = b.useTestnet
	result, err := binanceCli.NewClient(pubKey, secKey).
		NewGetOrderService().
		OrderID(orderId).
		Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		Id:          result.OrderID,
		OrderType:   string(result.Type),
		Status:      b.getStatus(string(result.Status)),
		Quantity:    result.OrigQuantity,
		Symbol:      result.Symbol,
		Exchange:    BinanceType,
		Price:       result.Price,
		Side:        string(result.Side),
		TimeInForce: string(result.TimeInForce),
		StopPrice:   result.StopPrice,
		IcebergQty:  result.IcebergQuantity,
	}

	return order, nil
}

func (b *BinanceAdapter) GetSymbols(pubKey, secKey, passPhrase string) ([]string, error) {
	binanceCli.UseTestnet = b.useTestnet
	result, err := binanceCli.NewClient(pubKey, secKey).
		NewExchangeInfoService().
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	symbols := make([]string, 0)
	for _, v := range result.Symbols {
		symbols = append(symbols, v.Symbol)
	}

	return symbols, err
}

func (b *BinanceAdapter) getStatus(status string) string {
	switch status {
	case string(binanceCli.OrderStatusTypeNew):
		return consts.OrderStatusActive
	case string(binanceCli.OrderStatusTypePartiallyFilled):
		return consts.OrderStatusPartFilled
	case string(binanceCli.OrderStatusTypeFilled):
		return consts.OrderStatusFilled
	case string(binanceCli.OrderStatusTypeCanceled):
		return consts.OrderStatusCanceled
	default:
		return ""
	}
}

func (b *BinanceAdapter) Depth(pubKey, secKey, passPhrase, symbol string) (*domain.Depth, error) {
	binanceCli.UseTestnet = b.useTestnet
	result, err := binanceCli.
		NewClient(pubKey, secKey).
		NewDepthService().
		Symbol(symbol).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	ask := make([]domain.PriceLevel, len(result.Asks), len(result.Asks))
	bid := make([]domain.PriceLevel, len(result.Bids), len(result.Bids))
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for k, v := range result.Asks {
			ask[k] = domain.PriceLevel{
				Price:    v.Price,
				Quantity: v.Quantity,
			}
		}
	}()
	go func() {
		defer wg.Done()
		for k, v := range result.Bids {
			bid[k] = domain.PriceLevel{
				Price:    v.Price,
				Quantity: v.Quantity,
			}
		}
	}()
	wg.Wait()

	return &domain.Depth{
		LastUpdateID: result.LastUpdateID,
		Asks:         ask,
		Bids:         bid,
	}, nil
}

func (b *BinanceAdapter) Trades(pubKey, secKey, passPhrase, symbol string) (domain.Balance, error) {
	binanceCli.UseTestnet = true
	result, err := binanceCli.NewClient(pubKey, secKey).
		NewGetAccountService().Do(context.Background())

	if err != nil {
		return nil, err
	}
	var balance domain.Balance

	if result.Balances == nil {
		return nil, errors.New("empty balances")
	}
	for _, v := range result.Balances {
		balSymb := domain.NewBalanceSymbol(v.Asset, v.Free)
		balance = append(balance, balSymb)
	}

	return balance, nil
}
