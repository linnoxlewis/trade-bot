package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"github.com/linnoxlewis/trade-bot/internal/service"
	"time"
)

type Databaser interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type OrderRepository struct {
	db          Databaser
	transaction *Transaction
}

func (o *OrderRepository) Atomic(ctx context.Context, fn func(ctx context.Context, orderRepo service.OrderRepo) error) (err error) {
	return o.transaction.Atomic(ctx, fn)
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
		transaction: &Transaction{
			db: db,
		},
	}
}

func (o *OrderRepository) CreateOrderWithSettings(ctx context.Context, order *domain.Order, settings *domain.Settings) (id int64, err error) {
	now := time.Now()
	db := o.transaction.GetDb(ctx)
	query := `INSERT INTO orders (user_id,
                    exchange,
                    symbol,
                    status,
                    side,
                    order_type,
                    quantity,
                    price,
                    exec_order_id, 
                    time_in_force, 
                    tp_sl, 
                    created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING ID`

	if err = db.QueryRowContext(ctx, query,
		order.UserId,
		order.Exchange,
		order.Symbol,
		order.Status,
		order.Side,
		order.OrderType,
		order.Quantity,
		order.Price,
		order.ExecOrderId,
		order.TimeInForce,
		order.TpSl, now).Scan(&id); err != nil {
		return 0, err
	}

	if settings != nil {
		query = `INSERT INTO order_settings (order_id, tp_percent, sl_percent, tp_price, sl_price, ts, tp_type, sl_type, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = db.ExecContext(ctx, query,
			id,
			settings.TpPercent,
			settings.SlPercent,
			settings.TpPrice,
			settings.SlPrice,
			settings.Ts,
			settings.TpType,
			settings.SlType, now)
		if err != nil {
			return 0, err
		}
	}

	return
}

func (o *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (id int64, err error) {
	query := `INSERT INTO orders (user_id,
                    exchange,
                    symbol,
                    status,
                    side,
                    order_type,
                    quantity,
                    price, 
                    exec_order_id, 
                    time_in_force, 
                    tp_sl,created_at) VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11,$12) RETURNING ID`
	if err = o.transaction.GetDb(ctx).QueryRowContext(ctx, query,
		order.UserId,
		order.Exchange,
		order.Symbol,
		order.Status,
		order.Side,
		order.OrderType,
		order.Quantity,
		order.Price,
		order.ExecOrderId,
		order.TimeInForce,
		order.TpSl,
		time.Now()).Scan(&id); err != nil {
		return 0, err
	}

	return
}

func (o *OrderRepository) GetActiveTpSlOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	query := `SELECT id,
       exchange,
       symbol,
       status,
       side,
       order_type,
       quantity,
       price, 
       exec_order_id,
       time_in_force,
       user_id,
       tp_sl FROM orders 
             WHERE exchange = $1 
               AND status = $2 
               AND tp_sl != $3`
	rows, err := o.db.QueryContext(ctx, query,
		exchange,
		consts.OrderStatusActive,
		consts.BaseOrderType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	orderList := make([]*domain.Order, 0)
	for rows.Next() {
		var mdl domain.Order
		err = rows.Scan(&mdl.Id,
			&mdl.Exchange,
			&mdl.Symbol,
			&mdl.Status,
			&mdl.Side,
			&mdl.OrderType,
			&mdl.Quantity,
			&mdl.Price,
			&mdl.ExecOrderId,
			&mdl.TimeInForce,
			&mdl.UserId,
			&mdl.TpSl)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, &mdl)
	}

	return orderList, rows.Err()
}

func (o *OrderRepository) CancelOrder(ctx context.Context, id int64, symbol, exchange string) error {
	query := `UPDATE orders SET 
                  status = $1, 
                  updated_at = $2
              WHERE id = $3
                AND symbol = $4 
                AND exchange = $5`
	_, err := o.transaction.GetDb(ctx).ExecContext(ctx, query,
		consts.OrderStatusCanceled,
		time.Now(),
		id,
		symbol,
		exchange)

	return err
}

func (o *OrderRepository) CancelOrderByEx(ctx context.Context, id int64, symbol, exchange string) error {
	query := `UPDATE orders SET status = $1 
              WHERE ( exec_order_id = $2 OR exec_order_id = 
                                       (SELECT id FROM orders WHERE exec_order_id = $3)) 
                AND status = $4 
                AND symbol = $5
                AND exchange = $6`
	_, err := o.transaction.GetDb(ctx).ExecContext(ctx,
		query,
		consts.OrderStatusCanceled,
		id,
		id,
		consts.OrderStatusActive,
		symbol,
		exchange)

	return err
}

func (o *OrderRepository) ExecuteOrder(ctx context.Context, id int64) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := o.transaction.GetDb(ctx).ExecContext(ctx,
		query,
		consts.OrderStatusFilled,
		id)

	return err
}

func (o *OrderRepository) UpdateTpSl(ctx context.Context, id int64, price string, settings *domain.Settings) error {
	db := o.transaction.GetDb(ctx)
	query := `UPDATE orders SET 
                  price = $1, 
                  updated_at = $2 
              WHERE id = $3`
	if _, err := db.ExecContext(ctx, query,
		price,
		time.Now(), id); err != nil {
		return err
	}

	if settings != nil {
		query = `UPDATE order_settings SET 
                          tp_price = $1, 
                          tp_percent = $2, 
                          sl_price = $3,
                          sl_percent = $4,
                      	  updated_at = $5
                  WHERE order_id = $6`
		if _, err := db.ExecContext(ctx, query,
			settings.TpPrice,
			settings.TpPercent,
			settings.SlPrice,
			settings.SlPercent,
			time.Now(),
			settings.OrderId); err != nil {
			return err
		}
	}

	return nil
}

func (o *OrderRepository) GetOrder(ctx context.Context, orderId int64, symbol, exchange string) (*domain.Order, error) {
	var result domain.Order

	query := `SELECT o.id,
       o.user_id, 
       o.exchange, 
       o.symbol, 
       o.status, 
       o.side, 
       o.order_type, 
       o.quantity, 
       o.price, 
       o.exec_order_id, 
       o.time_in_force, 
       o.tp_sl FROM orders o 
               WHERE exchange = $1 
                 AND id = $2 
                 AND symbol = $3 LIMIT 1`

	err := o.transaction.GetDb(ctx).
		QueryRowContext(ctx,
			query,
			exchange,
			orderId,
			symbol).Scan(&result.Id,
		&result.UserId,
		&result.Exchange,
		&result.Symbol,
		&result.Status,
		&result.Side,
		&result.OrderType,
		&result.Quantity,
		&result.Price,
		&result.ExecOrderId,
		&result.TimeInForce,
		&result.TpSl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (o *OrderRepository) GetTpSlOrderByBaseOrder(ctx context.Context, id int64, symbol, exchange, tpsl string) (*domain.Order, error) {
	var result domain.Order
	query := `SELECT o.id,
       o.user_id, 
       o.exchange, 
       o.symbol, 
       o.status, 
       o.side, 
       o.order_type, 
       o.quantity, 
       o.price, 
       o.exec_order_id, 
       o.time_in_force, 
       o.tp_sl FROM orders o
               WHERE exchange = $1
                 AND exec_order_id = $2 
                 AND symbol = $3 
                 AND tp_sl = $4 LIMIT 1`

	err := o.transaction.GetDb(ctx).
		QueryRowContext(ctx,
			query,
			exchange,
			id,
			symbol,
			tpsl).Scan(&result.Id,
		&result.UserId,
		&result.Exchange,
		&result.Symbol,
		&result.Status,
		&result.Side,
		&result.OrderType,
		&result.Quantity,
		&result.Price,
		&result.ExecOrderId,
		&result.TimeInForce,
		&result.TpSl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (o *OrderRepository) GetOpposingTpSlOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	var result domain.Order
	var tpSl string
	if order.TpSl == consts.TpOrderType {
		tpSl = consts.SlOrderType
	} else {
		tpSl = consts.TpOrderType
	}
	query := `SELECT id,
       exchange, 
       exec_order_id, 
       symbol FROM orders 
              WHERE exchange = $1 
                AND status = $2 
                AND tp_sl = $3 
                AND exec_order_id = $4`

	err := o.transaction.GetDb(ctx).
		QueryRowContext(ctx,
			query,
			order.Exchange,
			consts.OrderStatusActive,
			tpSl,
			order.ExecOrderId).Scan(&result.Id,
		&result.Exchange,
		&result.ExecOrderId,
		&result.Symbol)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (o *OrderRepository) GetActiveSymbols(ctx context.Context) (domain.SymbolList, error) {
	query := `SELECT DISTINCT symbol FROM orders WHERE status = $1`
	rows, err := o.db.QueryContext(ctx, query, consts.OrderStatusActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	symbolList := make([]domain.Symbol, 0)
	for rows.Next() {
		var symbol domain.Symbol
		if err := rows.Scan(&symbol); err != nil {
			return nil, err
		}
		symbolList = append(symbolList, symbol)
	}

	return symbolList, rows.Err()
}

func (o *OrderRepository) ActivateOrder(ctx context.Context, id int64) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := o.transaction.GetDb(ctx).ExecContext(ctx,
		query,
		consts.OrderStatusActive,
		id)

	return err
}

func (o *OrderRepository) GetTpSlOrdersByBaseOrder(ctx context.Context, id int64) ([]*domain.Order, error) {
	query := `SELECT id, 
       exchange,
       symbol,
       status, 
       side, 
       order_type,
       quantity,
       price,
       exec_order_id, 
       time_in_force, 
       user_id, 
       tp_sl FROM orders 
             WHERE exec_order_id = $1 AND tp_sl != $2`
	rows, err := o.transaction.GetDb(ctx).QueryContext(ctx, query,
		id,
		consts.BaseOrderType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	orderList := make([]*domain.Order, 0)
	for rows.Next() {
		var mdl domain.Order
		err = rows.Scan(&mdl.Id,
			&mdl.Exchange,
			&mdl.Symbol,
			&mdl.Status,
			&mdl.Side,
			&mdl.OrderType,
			&mdl.Quantity,
			&mdl.Price,
			&mdl.ExecOrderId,
			&mdl.TimeInForce,
			&mdl.UserId,
			&mdl.TpSl)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, &mdl)
	}

	return orderList, rows.Err()
}

func (o *OrderRepository) GetLimitOrders(ctx context.Context, exchange string) ([]*domain.Order, error) {
	query := `SELECT id, 
       exchange, 
       symbol, 
       status,
       side, 
       order_type, 
       quantity,
       price,
       exec_order_id, 
       time_in_force, 
       user_id, 
       tp_sl FROM orders WHERE tp_sl = $1 
                           AND status = $2 
                           AND exchange = $3`
	rows, err := o.db.QueryContext(ctx, query,
		consts.BaseOrderType,
		consts.OrderStatusActive,
		exchange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	orderList := make([]*domain.Order, 0)
	for rows.Next() {
		var mdl domain.Order
		err = rows.Scan(&mdl.Id,
			&mdl.Exchange,
			&mdl.Symbol,
			&mdl.Status,
			&mdl.Side,
			&mdl.OrderType,
			&mdl.Quantity,
			&mdl.Price,
			&mdl.ExecOrderId,
			&mdl.TimeInForce,
			&mdl.UserId,
			&mdl.TpSl)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, &mdl)
	}

	return orderList, rows.Err()
}
