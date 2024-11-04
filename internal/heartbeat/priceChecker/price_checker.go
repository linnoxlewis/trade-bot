package priceChecker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/service"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"net/url"
	"strings"
	"sync"
	"time"
)

var urlSocket = url.URL{
	Scheme: "wss",
	Host:   "stream.binance.com:9443",
	Path:   "/ws",
}

type PriceStorage interface {
	SetPrice(ctx context.Context, pair string, price string) error
	GetPrice(ctx context.Context, pair string) (string, error)
	Close() error
}

type CurrencyPairListener struct {
	mu           sync.RWMutex
	symbolSrv    *service.SymbolService
	priceStorage PriceStorage
	newPairsChan <-chan string
	done         chan struct{}
	wg           sync.WaitGroup
	connections  map[string]*websocket.Conn
	logger       *log.Logger
}

func NewCurrencyPairListener(symbolSrv *service.SymbolService,
	priceStorage PriceStorage,
	newPairsChan <-chan string,
	logger *log.Logger,
) *CurrencyPairListener {
	return &CurrencyPairListener{
		symbolSrv:    symbolSrv,
		priceStorage: priceStorage,
		newPairsChan: newPairsChan,
		done:         make(chan struct{}),
		connections:  make(map[string]*websocket.Conn),
		logger:       logger,
	}
}

func (l *CurrencyPairListener) Start(ctx context.Context) error {
	l.subscribeInitialPairs(ctx)

	l.wg.Add(1)
	go l.listenNewPairs(ctx)

	return nil
}

func (l *CurrencyPairListener) Stop() {
	close(l.done)
	l.wg.Wait()

	l.mu.Lock()
	for pair, conn := range l.connections {
		l.logger.InfoLog.Printf("Close WebSocket pair %s", pair)
		conn.Close()
	}
	l.mu.Unlock()

	if err := l.priceStorage.Close(); err != nil {
		l.logger.ErrorLog.Printf("Error close storage: %v", err)
	}

	l.logger.InfoLog.Println("Completed all connections")
}

// subscribeInitialPairs подписывается на начальные валютные пары через WebSocket.
func (l *CurrencyPairListener) subscribeInitialPairs(ctx context.Context) {
	initialSymbols := l.getInitialPairs(ctx)

	l.logger.InfoLog.Println("Initial symbols: ", initialSymbols)

	for _, pair := range initialSymbols {
		formattedPair := formatPairForExchange(string(pair))
		l.initialPrices(ctx, formattedPair)

		l.wg.Add(1)
		go l.listenPair(ctx, formattedPair)

		l.logger.InfoLog.Printf("Start tracking %s", formattedPair)
	}
}

func (l *CurrencyPairListener) connectToWebSocketForPair(pair string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(urlSocket.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("err connect to socket for %s: %w", pair, err)
	}

	l.logger.InfoLog.Printf("Connected webSocket for pair %s", pair)

	return conn, nil
}

// listenNewPairs слушает канал для получения новых валютных пар.
func (l *CurrencyPairListener) listenNewPairs(ctx context.Context) {
	defer l.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case pair := <-l.newPairsChan:
			formattedPair := formatPairForExchange(pair)

			if _, err := l.priceStorage.GetPrice(ctx, formattedPair); err != nil {
				l.priceStorage.SetPrice(ctx, formattedPair, "0.0")

				// Для каждой новой пары создаем отдельное соединение и подписываемся на нее
				l.wg.Add(1)
				go l.listenPair(ctx, formattedPair)

				l.logger.InfoLog.Printf("Добавлена новая пара: %s", formattedPair)
			}
		}
	}
}

// listenPair подписывается на WebSocket поток для конкретной валютной пары.
func (l *CurrencyPairListener) listenPair(ctx context.Context, pair string) {
	defer l.wg.Done()

	conn, err := l.connectToWebSocketForPair(pair)
	if err != nil {
		l.logger.ErrorLog.Printf("err connect to WebSocket for pair %s: %v", pair, err)
		return
	}

	// Сохраняем соединение для этой пары
	l.mu.Lock()
	l.connections[pair] = conn
	l.mu.Unlock()

	tradingPair := strings.ToLower(pair)
	subscribeMessage := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{fmt.Sprintf("%s@trade", tradingPair)},
		"id":     time.Now().Unix(),
	}

	if err := conn.WriteJSON(subscribeMessage); err != nil {
		l.logger.ErrorLog.Printf("err subsricbe on pair %s: %v", pair, err)

		return
	}

	// Чтение сообщений из WebSocket
	go l.readMessagesForPair(ctx, conn, pair)

	<-ctx.Done()

	// Отписываемся и закрываем соединение
	unsubscribeMessage := map[string]interface{}{
		"method": "UNSUBSCRIBE",
		"params": []string{fmt.Sprintf("%s@trade", tradingPair)}, // Отписываемся от торгового потока
		"id":     time.Now().Unix(),
	}

	if err := conn.WriteJSON(unsubscribeMessage); err != nil {
		l.logger.ErrorLog.Printf("Ошибка отписки от пары %s: %v", pair, err)
	}

	err = conn.Close()
	if err != nil {
		l.logger.ErrorLog.Printf("err close connection: %v", err)
	}

	l.mu.Lock()
	delete(l.connections, pair)
	l.mu.Unlock()

	l.logger.InfoLog.Printf("Connection webSocket closed for pair %s", pair)
}

// readMessagesForPair читает сообщения из WebSocket для конкретной валютной пары.
func (l *CurrencyPairListener) readMessagesForPair(ctx context.Context, conn *websocket.Conn, pair string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				l.logger.ErrorLog.Printf("Ошибка чтения сообщения для %s: %v", pair, err)
				return
			}

			l.handleWebSocketMessage(message)
		}
	}
}

// handleWebSocketMessage обрабатывает полученное сообщение и обновляет цену.
func (l *CurrencyPairListener) handleWebSocketMessage(message []byte) {
	var data struct {
		Symbol string `json:"s"`
		Price  string `json:"p"`
	}

	if err := json.Unmarshal(message, &data); err != nil {
		l.logger.ErrorLog.Printf("Error unmarshalling message: %v", err)
		return
	}

	if data.Symbol == "" {
		l.logger.InfoLog.Printf("Received message without symbol, likely a subscription response: %s", string(message))
		return
	}

	pair := strings.ToUpper(data.Symbol)

	if err := l.priceStorage.SetPrice(context.Background(), pair, data.Price); err != nil {
		l.logger.ErrorLog.Printf("Error updating price for %s: %v", pair, err)
		return
	}

	l.logger.InfoLog.Printf("New price for %s: %s", pair, data.Price)
}

// GetPrice возвращает текущую цену валютной пары.
func (l *CurrencyPairListener) GetPrice(pair string) (string, bool) {
	ctx := context.Background()
	price, err := l.priceStorage.GetPrice(ctx, pair)
	if err != nil {
		l.logger.ErrorLog.Printf("Error fetching price for %s: %v", pair, err)

		return "0.0", false
	}
	return price, true
}

func (l *CurrencyPairListener) getInitialPairs(ctx context.Context) domain.SymbolList {
	initialSymbols, err := l.symbolSrv.GetActiveSymbols(ctx)
	if err != nil {
		l.logger.ErrorLog.Println("Ошибка получения активных символов, используем символы по умолчанию")
		initialSymbols, _ = l.symbolSrv.GetDefaultSymbols(ctx)
	} else if len(initialSymbols) == 0 {
		l.logger.ErrorLog.Println("Не найдены активные символы, используем символы по умолчанию")
		initialSymbols, _ = l.symbolSrv.GetDefaultSymbols(ctx)
	}

	return initialSymbols
}

func (l *CurrencyPairListener) initialPrices(ctx context.Context, formattedPair string) {
	if _, err := l.priceStorage.GetPrice(ctx, formattedPair); err != nil {
		if err := l.priceStorage.SetPrice(ctx, formattedPair, "0.0"); err != nil {
			l.logger.ErrorLog.Printf("Ошибка установки начальной цены для %s: %v", formattedPair, err)
		}
	}
}

// Преобразует пару в формат Binance (например, "BTC/USD" -> "btcusdt").
func formatPairForExchange(pair string) string {
	return pair
	return strings.ToLower(strings.ReplaceAll(pair, "/", ""))
}
