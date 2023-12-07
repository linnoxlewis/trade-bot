package repository

import (
	"context"
	"database/sql"
	"github.com/linnoxlewis/trade-bot/internal/service"
)

const txKey = "TX"

type Transaction struct {
	db *sql.DB
}

func (t *Transaction) Atomic(ctx context.Context,
	fn func(ctx context.Context, order service.OrderRepo) error) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		defer context.WithValue(ctx, txKey, nil)
		if p := recover(); p != nil {
			_ = tx.Rollback()

			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	orderStore := &OrderRepository{
		db: tx,
	}
	ctx = context.WithValue(ctx, txKey, tx)
	err = fn(ctx, orderStore)

	return
}

func (t *Transaction) getTxFromContext(ctx context.Context) Databaser {
	result := ctx.Value(txKey)
	if result != nil {
		res, ok := result.(Databaser)
		if !ok {
			return nil
		}

		return res
	}

	return nil
}

func (t *Transaction) GetDb(ctx context.Context) Databaser {
	db := t.getTxFromContext(ctx)
	if db != nil {
		return db
	}

	return t.db
}
