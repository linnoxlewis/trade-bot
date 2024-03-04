package tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type OrderRepo interface {
	UpdateTpSl(ctx context.Context, id int64, price string, settings *domain.Settings) error
	ExecuteOrder(ctx context.Context, id int64) error
	ActivateOrder(ctx context.Context, id int64) error

	GetOrder(ctx context.Context, orderId int64, symbol, exchange string) (*domain.Order, error)
	GetTpSlOrderByBaseOrder(ctx context.Context, id int64, symbol, exchange, tpsl string) (*domain.Order, error)
	GetTpSlOrdersByBaseOrder(ctx context.Context, id int64) ([]*domain.Order, error)
	GetOpposingTpSlOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetActiveSymbols(ctx context.Context) (domain.SymbolList, error)
	GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error)

	Atomic(ctx context.Context, fn func(ctx context.Context, orderRepo OrderRepo) error) (err error)
}

func TestOrderRepository_CreateOrder(t *testing.T) {
	tests := []struct {
		name    string
		order   domain.Order
		prepare func(mock sqlmock.Sqlmock)
		check   func(t *testing.T, id int64, err error)
	}{
		{
			name: "success",
			order: domain.Order{
				UserId:      123,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "none",
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("^INSERT INTO orders").
					WithArgs(
						"user123", "exchangeA", "BTC/USD", "active", "buy", "market", "1", "50000", "exec123", "GTC", "none", sqlmock.AnyArg(),
					).WillReturnRows(rows) // Вернуть строку с id
			},
			check: func(t *testing.T, id int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id) // Проверить, что id равен 1
			},
		},
		{
			name: "database_error",
			order: domain.Order{
				UserId:      123,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "none",
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^INSERT INTO orders").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
				assert.Equal(t, int64(0), id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			id, err := repo.CreateOrder(context.Background(), &tt.order)

			tt.check(t, id, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetActiveTpSlOrders(t *testing.T) {
	tests := []struct {
		name     string
		exchange string
		prepare  func(mock sqlmock.Sqlmock)
		check    func(t *testing.T, orders []*domain.Order, err error)
	}{
		{
			name:     "success",
			exchange: "exchangeA",
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "user_id", "tp_sl"}).
					AddRow(1, "exchangeA", "BTC/USD", "active", "buy", "market", "1", "50000", "exec123", "GTC", "user123", "none").
					AddRow(2, "exchangeA", "ETH/USD", "active", "sell", "limit", "2", "2500", "exec456", "GTC", "user456", "none")
				mock.ExpectQuery("^SELECT id, exchange, symbol").WillReturnRows(rows)
			},
			check: func(t *testing.T, orders []*domain.Order, err error) {
				assert.NoError(t, err)
				assert.Len(t, orders, 2)
				assert.Equal(t, int64(1), orders[0].Id)
				assert.Equal(t, "exchangeA", orders[0].Exchange)
				// проверить другие поля заказа
			},
		},
		{
			name:     "no_orders",
			exchange: "exchangeA",
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "user_id", "tp_sl"})
				mock.ExpectQuery("^SELECT id, exchange, symbol").WillReturnRows(rows)
			},
			check: func(t *testing.T, orders []*domain.Order, err error) {
				assert.NoError(t, err)
				assert.Len(t, orders, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			orders, err := repo.GetActiveTpSlOrders(context.Background(), tt.exchange)

			tt.check(t, orders, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_CancelOrder(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		symbol   string
		exchange string
		prepare  func(mock sqlmock.Sqlmock)
		check    func(t *testing.T, err error)
	}{
		{
			name:     "success",
			id:       1,
			symbol:   "BTC/USD",
			exchange: "exchangeA",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders SET").
					WithArgs("canceled", sqlmock.AnyArg(), 1, "BTC/USD", "exchangeA").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "database_error",
			id:       1,
			symbol:   "BTC/USD",
			exchange: "exchangeA",
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders SET").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			err = repo.CancelOrder(context.Background(), tt.id, tt.symbol, tt.exchange)

			tt.check(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_CreateOrderWithSettings(t *testing.T) {
	tests := []struct {
		name     string
		order    *domain.Order
		settings *domain.Settings
		prepare  func(mock sqlmock.Sqlmock)
		check    func(t *testing.T, id int64, err error)
	}{
		{
			name: "success",
			order: &domain.Order{
				UserId:      1,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "none",
			},
			settings: &domain.Settings{
				TpPercent: "10",
				SlPercent: "5",
				TpPrice:   "51000",
				SlPrice:   "49000",
				Ts:        "5",
				TpType:    "price",
				SlType:    "price",
				OrderId:   1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("^INSERT INTO orders").WillReturnRows(rows)

				rowsSettings := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("^INSERT INTO order_settings").WillReturnRows(rowsSettings)
			},
			check: func(t *testing.T, id int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id)
			},
		},
		{
			name: "database_error",
			order: &domain.Order{
				UserId:      1,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "none",
			},
			settings: &domain.Settings{
				TpPercent: "10",
				SlPercent: "5",
				TpPrice:   "51000",
				SlPrice:   "49000",
				Ts:        "5",
				TpType:    "price",
				SlType:    "price",
				OrderId:   1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^INSERT INTO orders").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			id, err := repo.CreateOrderWithSettings(context.Background(), tt.order, tt.settings)

			tt.check(t, id, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_UpdateTpSl(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		price    string
		settings *domain.Settings
		prepare  func(mock sqlmock.Sqlmock)
		check    func(t *testing.T, err error)
	}{
		{
			name:  "success",
			id:    1,
			price: "100.00",
			settings: &domain.Settings{
				TpPrice:   "50.00",
				TpPercent: "0.05",
				SlPrice:   "80.00",
				SlPercent: "0.03",
				OrderId:   1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec("^UPDATE order_settings").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:  "database_error",
			id:    1,
			price: "100.00",
			settings: &domain.Settings{
				TpPrice:   "50.00",
				TpPercent: "0.05",
				SlPrice:   "80.00",
				SlPercent: "0.03",
				OrderId:   1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			err = repo.UpdateTpSl(context.Background(), tt.id, tt.price, tt.settings)

			tt.check(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_ExecuteOrder(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		prepare func(mock sqlmock.Sqlmock)
		check   func(t *testing.T, err error)
	}{
		{
			name: "success",
			id:   1,
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "database_error",
			id:   1,
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			err = repo.ExecuteOrder(context.Background(), tt.id)

			tt.check(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_ActivateOrder(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		prepare func(mock sqlmock.Sqlmock)
		check   func(t *testing.T, err error)
	}{
		{
			name: "success",
			id:   1,
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "database_error",
			id:   1,
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^UPDATE orders").WillReturnError(errors.New("database error"))
			},
			check: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.prepare(mock)

			err = repo.ActivateOrder(context.Background(), tt.id)

			tt.check(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetOrder(t *testing.T) {
	tests := []struct {
		name     string
		orderID  int64
		symbol   string
		exchange string
		queryRow func(mock sqlmock.Sqlmock)
		want     *domain.Order
		wantErr  error
	}{
		{
			name:     "success",
			orderID:  1,
			symbol:   "BTC/USD",
			exchange: "exchangeA",
			queryRow: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "tp_sl"}).
					AddRow(1, 1, "exchangeA", "BTC/USD", "active", "buy", "market", "1", "50000", 123, "GTC", "none")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			want: &domain.Order{
				Id:          1,
				UserId:      1,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "none",
			},
			wantErr: nil,
		},
		{
			name:     "not_found",
			orderID:  2,
			symbol:   "ETH/USD",
			exchange: "exchangeB",
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name:     "database_error",
			orderID:  3,
			symbol:   "LTC/USD",
			exchange: "exchangeC",
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRow(mock)

			got, err := repo.GetOrder(context.Background(), tt.orderID, tt.symbol, tt.exchange)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetTpSlOrderByBaseOrder(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		symbol   string
		exchange string
		tpsl     string
		queryRow func(mock sqlmock.Sqlmock)
		want     *domain.Order
		wantErr  error
	}{
		{
			name:     "success",
			id:       1,
			symbol:   "BTC/USD",
			exchange: "exchangeA",
			tpsl:     "tpsl",
			queryRow: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "tp_sl"}).
					AddRow(1, 1, "exchangeA", "BTC/USD", "active", "buy", "market", "1", "50000", 123, "GTC", "tpsl")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			want: &domain.Order{
				Id:          1,
				UserId:      1,
				Exchange:    "exchangeA",
				Symbol:      "BTC/USD",
				Status:      "active",
				Side:        "buy",
				OrderType:   "market",
				Quantity:    "1",
				Price:       "50000",
				ExecOrderId: 123,
				TimeInForce: "GTC",
				TpSl:        "tpsl",
			},
			wantErr: nil,
		},
		{
			name:     "not_found",
			id:       2,
			symbol:   "ETH/USD",
			exchange: "exchangeB",
			tpsl:     "tpsl",
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name:     "database_error",
			id:       3,
			symbol:   "LTC/USD",
			exchange: "exchangeC",
			tpsl:     "tpsl",
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRow(mock)

			got, err := repo.GetTpSlOrderByBaseOrder(context.Background(), tt.id, tt.symbol, tt.exchange, tt.tpsl)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetOpposingTpSlOrder(t *testing.T) {
	tests := []struct {
		name     string
		order    *domain.Order
		queryRow func(mock sqlmock.Sqlmock)
		want     *domain.Order
		wantErr  error
	}{
		{
			name: "success_tp",
			order: &domain.Order{
				Exchange:    "exchangeA",
				ExecOrderId: 1,
				TpSl:        "tp",
			},
			queryRow: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "exec_order_id", "symbol"}).
					AddRow(2, "exchangeA", 1, "BTC/USD")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			want: &domain.Order{
				Id:          2,
				Exchange:    "exchangeA",
				ExecOrderId: 1,
				Symbol:      "BTC/USD",
			},
			wantErr: nil,
		},
		{
			name: "success_sl",
			order: &domain.Order{
				Exchange:    "exchangeB",
				ExecOrderId: 3,
				TpSl:        "sl",
			},
			queryRow: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "exec_order_id", "symbol"}).
					AddRow(4, "exchangeB", 3, "ETH/USD")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			want: &domain.Order{
				Id:          4,
				Exchange:    "exchangeB",
				ExecOrderId: 3,
				Symbol:      "ETH/USD",
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			order: &domain.Order{
				Exchange:    "exchangeC",
				ExecOrderId: 5,
				TpSl:        "tp",
			},
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "database_error",
			order: &domain.Order{
				Exchange:    "exchangeD",
				ExecOrderId: 6,
				TpSl:        "tp",
			},
			queryRow: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRow(mock)

			got, err := repo.GetOpposingTpSlOrder(context.Background(), tt.order)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetActiveSymbols(t *testing.T) {
	tests := []struct {
		name        string
		queryRows   func(mock sqlmock.Sqlmock)
		wantSymbols domain.SymbolList
		wantErr     error
	}{
		{
			name: "success",
			queryRows: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"symbol"}).
					AddRow("BTC/USD").
					AddRow("ETH/USD").
					AddRow("LTC/USD")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			wantSymbols: domain.SymbolList{"BTC/USD", "ETH/USD", "LTC/USD"},
			wantErr:     nil,
		},
		{
			name: "no_rows",
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantSymbols: nil,
			wantErr:     nil,
		},
		{
			name: "database_error",
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			wantSymbols: nil,
			wantErr:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRows(mock)

			gotSymbols, err := repo.GetActiveSymbols(context.Background())

			assert.Equal(t, tt.wantSymbols, gotSymbols)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetTpSlOrdersByBaseOrder(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		queryRows  func(mock sqlmock.Sqlmock)
		wantOrders []*domain.Order
		wantErr    error
	}{
		{
			name: "success",
			id:   1,
			queryRows: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "user_id", "tp_sl"}).
					AddRow(1, "exchangeA", "BTC/USD", "active", "buy", "market", "1", "50000", 123, "GTC", 1, "tpsl").
					AddRow(2, "exchangeB", "ETH/USD", "active", "sell", "limit", "2", "60000", 124, "GTC", 2, "sl")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			wantOrders: []*domain.Order{
				{
					Id:          1,
					Exchange:    "exchangeA",
					Symbol:      "BTC/USD",
					Status:      "active",
					Side:        "buy",
					OrderType:   "market",
					Quantity:    "1",
					Price:       "50000",
					ExecOrderId: 123,
					TimeInForce: "GTC",
					UserId:      1,
					TpSl:        "tpsl",
				},
				{
					Id:          2,
					Exchange:    "exchangeB",
					Symbol:      "ETH/USD",
					Status:      "active",
					Side:        "sell",
					OrderType:   "limit",
					Quantity:    "2",
					Price:       "60000",
					ExecOrderId: 124,
					TimeInForce: "GTC",
					UserId:      2,
					TpSl:        "sl",
				},
			},
			wantErr: nil,
		},
		{
			name: "no_rows",
			id:   3,
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantOrders: nil,
			wantErr:    nil,
		},
		{
			name: "database_error",
			id:   4,
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			wantOrders: nil,
			wantErr:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRows(mock)

			gotOrders, err := repo.GetTpSlOrdersByBaseOrder(context.Background(), tt.id)

			assert.Equal(t, tt.wantOrders, gotOrders)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetLimitOrders(t *testing.T) {
	tests := []struct {
		name       string
		exchange   string
		queryRows  func(mock sqlmock.Sqlmock)
		wantOrders []*domain.Order
		wantErr    error
	}{
		{
			name:     "success",
			exchange: "exchangeA",
			queryRows: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "exchange", "symbol", "status", "side", "order_type", "quantity", "price", "exec_order_id", "time_in_force", "user_id", "tp_sl"}).
					AddRow(1, "exchangeA", "BTC/USD", "active", "buy", "limit", "1", "50000", 123, "GTC", 1, "tpsl").
					AddRow(2, "exchangeA", "ETH/USD", "active", "sell", "limit", "2", "60000", 124, "GTC", 2, "tpsl")
				mock.ExpectQuery("^SELECT").WillReturnRows(rows)
			},
			wantOrders: []*domain.Order{
				{
					Id:          1,
					Exchange:    "exchangeA",
					Symbol:      "BTC/USD",
					Status:      "active",
					Side:        "buy",
					OrderType:   "limit",
					Quantity:    "1",
					Price:       "50000",
					ExecOrderId: 123,
					TimeInForce: "GTC",
					UserId:      1,
					TpSl:        "tpsl",
				},
				{
					Id:          2,
					Exchange:    "exchangeA",
					Symbol:      "ETH/USD",
					Status:      "active",
					Side:        "sell",
					OrderType:   "limit",
					Quantity:    "2",
					Price:       "60000",
					ExecOrderId: 124,
					TimeInForce: "GTC",
					UserId:      2,
					TpSl:        "tpsl",
				},
			},
			wantErr: nil,
		},
		{
			name:     "no_rows",
			exchange: "exchangeB",
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantOrders: nil,
			wantErr:    nil,
		},
		{
			name:     "database_error",
			exchange: "exchangeC",
			queryRows: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT").WillReturnError(errors.New("database error"))
			},
			wantOrders: nil,
			wantErr:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := repository.NewOrderRepository(db)
			tt.queryRows(mock)

			gotOrders, err := repo.GetLimitOrders(context.Background(), tt.exchange)

			assert.Equal(t, tt.wantOrders, gotOrders)
			assert.Equal(t, tt.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
