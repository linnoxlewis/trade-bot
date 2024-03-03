package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linnoxlewis/trade-bot/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetApiKeysByUserIdAndExchange(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "exchange", "pub_key", "priv_key", "passphrase"}).
		AddRow(1, "testExchange", "pubTestKey", "privTestKey", "testPassphrase")

	mock.ExpectQuery(`SELECT user_id, exchange, pub_key, priv_key, passphrase FROM user_api_keys WHERE user_id = \$1 AND exchange = \$2`).
		WithArgs(1, "testExchange").
		WillReturnRows(rows)

	repo := repository.NewApiKeyRepo(db)
	apiKeys, err := repo.GetApiKeysByUserIdAndExchange(context.Background(), 1, "testExchange")

	assert.NoError(t, err)
	assert.NotNil(t, apiKeys)
	assert.Equal(t, int64(1), apiKeys.UserId)
	assert.Equal(t, "testExchange", apiKeys.Exchange)
	assert.Equal(t, "pubTestKey", apiKeys.PubKey)
	assert.Equal(t, "privTestKey", apiKeys.PrivKey)
	assert.Equal(t, "testPassphrase", apiKeys.Passphrase)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
