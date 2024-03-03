package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linnoxlewis/trade-bot/internal/domain"
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

func TestAddApiKeys(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	apiKeys := &domain.ApiKeys{
		UserId:     1,
		Exchange:   "testExchange",
		PubKey:     "pubTestKey",
		PrivKey:    "privTestKey",
		Passphrase: "testPassphrase",
	}

	mock.ExpectExec(`INSERT INTO user_api_keys`).WithArgs(apiKeys.UserId, apiKeys.Exchange, apiKeys.PubKey, apiKeys.PrivKey, apiKeys.Passphrase, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewApiKeyRepo(db)
	err = repo.AddApiKeys(context.Background(), apiKeys)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteApiKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := int64(1)
	exchange := "testExchange"

	mock.ExpectExec(`DELETE FROM user_api_keys WHERE user_id = \$1 AND exchange = \$2`).WithArgs(userId, exchange).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := repository.NewApiKeyRepo(db)
	err = repo.DeleteApiKey(context.Background(), userId, exchange)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClearApiKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := int64(1)

	mock.ExpectExec(`DELETE FROM user_api_keys WHERE user_id = \$1`).WithArgs(userId).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := repository.NewApiKeyRepo(db)
	err = repo.ClearApiKey(context.Background(), userId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetApiKeysByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userId := int64(1)
	rows := sqlmock.NewRows([]string{"user_id", "exchange", "pub_key", "priv_key", "passphrase"}).
		AddRow(userId, "testExchange", "pubTestKey", "privTestKey", "testPassphrase")

	mock.ExpectQuery(`SELECT user_id, exchange, pub_key, priv_key, passphrase FROM user_api_keys WHERE user_id = \$1`).WithArgs(userId).WillReturnRows(rows)

	repo := repository.NewApiKeyRepo(db)
	apiKeys, err := repo.GetApiKeysByUserId(context.Background(), userId)

	assert.NoError(t, err)
	assert.NotNil(t, apiKeys)
	assert.Len(t, apiKeys, 1)
	assert.Equal(t, "testExchange", apiKeys[0].Exchange)
	assert.NoError(t, mock.ExpectationsWereMet())
}
