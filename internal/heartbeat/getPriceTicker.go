package heartbeat

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"strings"
	"sync"
)

type tradeUpdate struct {
	EventType     string `json:"e"`
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	TradeID       int64  `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	BuyerOrderID  int64  `json:"b"`
	SellerOrderID int64  `json:"a"`
	TradeTime     int64  `json:"T"`
	BuyerIsMarket bool   `json:"m"`
	Ignore        bool   `json:"M"`
}

var (
	errUknownExchange   = errors.New("err unknown exchange ")
	binanceTradeChannel = "@trade"
	binanceSocketUrl    = "wss://stream.binance.com:9443/ws/"
)

type TradePriceTicker struct {
	cfg      *config.Config
	conn     *websocket.Conn
	keyDbCli *redis.Client
	logger   *log.Logger
	symbol   string
	exchange string
	sync.Mutex
}

func NewTradePriceTicker(cfg *config.Config,
	keyDbCli *redis.Client,
	exchange string,
	symbol string,
	logger *log.Logger) (*TradePriceTicker, error) {
	var socketUrl string
	switch exchange {
	case consts.Binance:
		socketUrl = binanceSocketUrl + strings.ToLower(symbol) + binanceTradeChannel
		break
	default:
		return nil, errUknownExchange
	}
	conn, err := establishWebSocketConnection(socketUrl)
	if err != nil {
		logger.ErrorLog.Println("err socket connection:", err)

		return nil, err
	}

	return &TradePriceTicker{
		cfg:      cfg,
		conn:     conn,
		logger:   logger,
		keyDbCli: keyDbCli,
		symbol:   symbol,
		exchange: exchange,
	}, nil
}

func establishWebSocketConnection(endpoint string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (t *TradePriceTicker) HandleWebSocketData(ctx context.Context) {
	t.logger.InfoLog.Println("Start check price " + t.exchange + " " + t.symbol)
	for {
		select {
		case <-ctx.Done():
			t.logger.InfoLog.Println("program done")

			return
		default:
			_, message, err := t.conn.ReadMessage()
			if err != nil {
				t.logger.ErrorLog.Printf("err read message: %v\n", err)

				return
			}

			var tradeUpdate tradeUpdate
			if err = json.Unmarshal(message, &tradeUpdate); err != nil {
				t.logger.ErrorLog.Printf("err unmarshal message: %v\n", err)

				continue
			}

			t.SetPriceToCache(ctx, tradeUpdate.Price)
		}
	}
}

func (t *TradePriceTicker) SetPriceToCache(ctx context.Context, price string) {
	t.Lock()
	t.keyDbCli.Set(ctx, consts.TradePriceCacheKey+t.exchange+"_"+t.symbol,
		price,
		redis.KeepTTL)
	t.Unlock()
}
