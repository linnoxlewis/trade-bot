package excahnger

import (
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
)

const (
	BinanceType = "binance"
	KucoinType  = "kucoin"
)

type Exchanger interface {
	CreateOrder(keys *domain.ApiKeys, order *dto.Order) (int64, error)
	CancelOrder(keys *domain.ApiKeys, order *dto.CancelOrder) error
	//TODO::UpdateOrders(order []*dto.UpdateOrder, userId uuid.UUID) ([]*domain.Order, error)
}

type ExchangerCli interface {
	CreateOrder(pubKey, secKey, passPhrase string, order *dto.Order) (int64, error)
	CancelOrder(pubKey, secKey, passPhrase string, order *dto.CancelOrder) error
	//TODO::UpdateOrders(order []*dto.UpdateOrder, userId uuid.UUID) ([]*domain.Order, error)
}

type ExchangeCli struct {
	binanceCli *BinanceAdapter
	exType     string
}

func NewExchanger() *ExchangeCli {
	return &ExchangeCli{
		binanceCli: NewBinanceAdapter(),
	}
}

func (e *ExchangeCli) getType(exType string) ExchangerCli {
	switch exType {
	case BinanceType:
		return e.binanceCli
	default:
		return nil
	}
}

func (e *ExchangeCli) CreateOrder(keys *domain.ApiKeys, order *dto.Order) (int64, error) {
	return e.getType(order.Exchange).CreateOrder(keys.PubKey, keys.PrivKey, keys.Passphrase, order)
}

func (e *ExchangeCli) CancelOrder(keys *domain.ApiKeys, order *dto.CancelOrder) error {
	return e.getType(order.Exchange).CancelOrder(keys.PubKey, keys.PrivKey, keys.Passphrase, order)
}
