package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"time"
)

type SymbolRepository struct {
	db *sql.DB
}

func NewSymbolRepository(db *sql.DB) *SymbolRepository {
	return &SymbolRepository{
		db: db,
	}
}

func (s *SymbolRepository) AddSymbol(ctx context.Context, symbol string) error {
	query := `INSERT INTO default_symbols (symbol, created_at) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, symbol, time.Now())

	return err
}

func (s *SymbolRepository) GetDefaultSymbols(ctx context.Context) (domain.SymbolList, error) {
	query := `SELECT symbol FROM default_symbols`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	symbolList := make(domain.SymbolList, 0)
	for rows.Next() {
		var symbol domain.Symbol
		if err := rows.Scan(&symbol); err != nil {
			return nil, err
		}
		symbolList = append(symbolList, symbol)
	}

	return symbolList, rows.Err()
}
