package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"testing"
)

func TestAddSymbol(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewSymbolRepository(db)
	symbol := "BTCUSD"

	mock.ExpectExec("INSERT INTO default_symbols").WithArgs(symbol, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddSymbol(context.Background(), symbol)
	if err != nil {
		t.Errorf("error was not expected while adding symbol: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDefaultSymbols(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewSymbolRepository(db)

	rows := sqlmock.NewRows([]string{"symbol"}).
		AddRow("BTCUSD").
		AddRow("ETHUSD")

	mock.ExpectQuery("SELECT symbol FROM default_symbols").WillReturnRows(rows)

	symbols, err := repo.GetDefaultSymbols(context.Background())
	if err != nil {
		t.Fatalf("error was not expected while getting default symbols: %s", err)
	}

	expectedSymbols := domain.SymbolList{"BTCUSD", "ETHUSD"}
	if len(symbols) != len(expectedSymbols) {
		t.Errorf("expected %d symbols, got %d", len(expectedSymbols), len(symbols))
	}

	for i, symbol := range symbols {
		if symbol != expectedSymbols[i] {
			t.Errorf("expected symbol %s, got %s", expectedSymbols[i], symbol)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
