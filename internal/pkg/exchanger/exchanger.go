package exchanger

import (
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/dto"
)

const (
	BinanceType = "binance"
	KucoinType  = "kucoin"
	OkxType     = "okx"
	HuobiType   = "huobi"
)

var ExchangeList = []string{BinanceType}

type Exchanger interface {
	CreateOrder(keys *domain.ApiKeys, order *dto.Order) (int64, error)
	CancelOrder(keys *domain.ApiKeys, order *dto.CancelOrder) error
	UpdateOrder(keys *domain.ApiKeys, order *dto.UpdateOrder) (int64, error)
	Balance(keys *domain.ApiKeys, exchange string) (domain.Balance, error)
	GetOpenOrders(keys *domain.ApiKeys, exchange, symbol string) ([]domain.Order, error)
	GetOrder(keys *domain.ApiKeys, exchange, symbol string, orderId int64) (*domain.Order, error)
	GetSymbols(keys *domain.ApiKeys, exchange string) ([]string, error)
}

type ExchangerCli interface {
	CreateOrder(pubKey, secKey, passPhrase string, order *dto.Order) (int64, error)
	CancelOrder(pubKey, secKey, passPhrase string, order *dto.CancelOrder) error
	UpdateOrder(pubKey, secKey, passPhrase string, order *dto.UpdateOrder) (int64, error)
	GetBalance(pubKey, secKey, passPhrase string) (domain.Balance, error)
	GetOpenOrders(pubKey, secKey, passPhrase, symbol string) ([]domain.Order, error)
	GetOrder(pubKey, secKey, passPhrase, symbol string, orderId int64) (*domain.Order, error)
	GetSymbols(pubKey, secKey, passPhrase string) ([]string, error)
}

type ExchangeCli struct {
	binanceCli *BinanceAdapter
	exType     string
	isTestnet  bool
}

func NewExchanger(isTestnet bool) *ExchangeCli {
	return &ExchangeCli{
		binanceCli: NewBinanceAdapter(isTestnet),
		isTestnet:  isTestnet,
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

func (e *ExchangeCli) UpdateOrder(keys *domain.ApiKeys, order *dto.UpdateOrder) (int64, error) {
	return e.getType(order.Exchange).UpdateOrder(keys.PubKey, keys.PrivKey, keys.Passphrase, order)
}

func (e *ExchangeCli) Balance(keys *domain.ApiKeys, exchange string) (domain.Balance, error) {
	return e.getType(exchange).GetBalance(keys.PubKey, keys.PrivKey, keys.Passphrase)
}

func (e *ExchangeCli) GetOpenOrders(keys *domain.ApiKeys, exchange, symbol string) ([]domain.Order, error) {
	return e.getType(exchange).GetOpenOrders(keys.PubKey, keys.PrivKey, keys.Passphrase, symbol)
}

func (e *ExchangeCli) GetOrder(keys *domain.ApiKeys, exchange, symbol string, orderId int64) (*domain.Order, error) {
	return e.getType(exchange).GetOrder(keys.PubKey, keys.PrivKey, keys.Passphrase, symbol, orderId)
}

func (e *ExchangeCli) GetSymbols(keys *domain.ApiKeys, exchange string) ([]string, error) {
	return e.getType(exchange).GetSymbols(keys.PubKey, keys.PrivKey, keys.Passphrase)
}
