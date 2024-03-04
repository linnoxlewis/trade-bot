package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"time"
)

type ApiKeyRepo struct {
	db *sql.DB
}

func NewApiKeyRepo(db *sql.DB) *ApiKeyRepo {
	return &ApiKeyRepo{
		db: db,
	}
}

func (a *ApiKeyRepo) AddApiKeys(ctx context.Context, keys *domain.ApiKeys) (err error) {
	query := `INSERT INTO user_api_keys (user_id, 
                           exchange, 
                           pub_key, 
                           priv_key,
                           passphrase, 
                           created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = a.db.ExecContext(ctx, query,
		keys.UserId,
		keys.Exchange,
		keys.PubKey,
		keys.PrivKey,
		keys.Passphrase,
		time.Now())

	return
}

func (a *ApiKeyRepo) DeleteApiKey(ctx context.Context, userId int64, exchange string) (err error) {
	query := `DELETE FROM user_api_keys WHERE user_id = $1 AND exchange = $2`
	_, err = a.db.ExecContext(ctx, query, userId, exchange)

	return
}

func (a *ApiKeyRepo) ClearApiKey(ctx context.Context, userId int64) (err error) {
	query := `DELETE FROM user_api_keys WHERE user_id = $1`
	_, err = a.db.ExecContext(ctx, query, userId)

	return
}

func (a *ApiKeyRepo) GetApiKeysByUserIdAndExchange(ctx context.Context, userId int64, exchange string) (*domain.ApiKeys, error) {
	var apiKeys domain.ApiKeys
	query := `SELECT user_id,
       exchange, 
       pub_key, 
       priv_key, 
       passphrase FROM user_api_keys 
                  WHERE user_id = $1 AND exchange = $2`
	row := a.db.QueryRowContext(ctx, query, userId, exchange)
	err := row.Scan(&apiKeys.UserId,
		&apiKeys.Exchange,
		&apiKeys.PubKey,
		&apiKeys.PrivKey,
		&apiKeys.Passphrase)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &apiKeys, nil
}

func (a *ApiKeyRepo) GetApiKeysByUserId(ctx context.Context, userId int64) (apiKeys []*domain.ApiKeys, err error) {
	query := `SELECT user_id, 
       exchange, 
       pub_key, 
       priv_key, 
       passphrase FROM user_api_keys
                  WHERE user_id = $1`
	rows, err := a.db.QueryContext(ctx, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	apiKeysList := make([]*domain.ApiKeys, 0)
	for rows.Next() {
		mdl := &domain.ApiKeys{}
		err = rows.Scan(&mdl.UserId,
			&mdl.Exchange,
			&mdl.PubKey,
			&mdl.PrivKey,
			&mdl.Passphrase)
		if err != nil {
			return nil, err
		}
		apiKeysList = append(apiKeysList, mdl)
	}

	return apiKeysList, rows.Err()
}
