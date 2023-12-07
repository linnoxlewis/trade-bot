package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"time"
)

type ApiKeyRepo struct {
	db *sql.DB
}

func NewKeysStorage(db *sql.DB) *ApiKeyRepo {
	return &ApiKeyRepo{
		db: db,
	}
}

func (a *ApiKeyRepo) AddApiKeys(ctx context.Context, keys *domain.ApiKeys) (err error) {
	query := `INSERT INTO user_api_keys (user_id, exchange, pub_key, priv_key, passphrase, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = a.db.ExecContext(ctx,
		query,
		keys.UserId, keys.Exchange, keys.PubKey, keys.PrivKey, keys.Passphrase, time.Now())

	return
}

func (a *ApiKeyRepo) DeleteApiKey(ctx context.Context, userId uuid.UUID, exchange string) (err error) {
	query := `SELECT FROM user_api_keys WHERE user_id = $1 AND exchange = $2`
	_, err = a.db.ExecContext(ctx, query, userId, exchange)

	return
}

func (a *ApiKeyRepo) ClearApiKey(ctx context.Context, userId uuid.UUID) (err error) {
	query := `DELETE FROM user_api_keys WHERE user_id = $1`
	_, err = a.db.ExecContext(ctx, query, userId)

	return
}

func (a *ApiKeyRepo) GetApiKeysByUserIdAndExchange(ctx context.Context, userId uuid.UUID, exchange string) (apiKeys *domain.ApiKeys, err error) {
	query := `SELECT user_id, exchange, pub_key, priv_key, passphrase FROM users_api_keys 
				WHERE user_id = $1 AND exchange = $2`
	row := a.db.QueryRowContext(ctx, query, userId, exchange)
	err = row.Scan(&apiKeys)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return apiKeys, nil
}

func (a *ApiKeyRepo) GetByApiKeyByTgUserIdAndExchange(ctx context.Context, tgId int, exchange string) (apiKeys *domain.ApiKeys, err error) {
	query := `SELECT a.user_id, a.exchange, a. pub_key, a.priv_key, a.passphrase 
			  FROM finder_questions_answers a 
              INNER JOIN users u ON a.user_id = u.id
			  WHERE a.tg_id = $1 AND exchange = $2 LIMIT 1`
	row := a.db.QueryRowContext(ctx, query, tgId, exchange)
	err = row.Scan(&apiKeys)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return apiKeys, nil
}

func (a *ApiKeyRepo) GetApiKeysByUserId(ctx context.Context, userId uuid.UUID) (apiKeys []*domain.ApiKeys, err error) {
	query := `SELECT user_id, exchange, pub_key, priv_key, passphrase FROM finder_questions_answers 
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
		var mdl *domain.ApiKeys
		err = rows.Scan(&mdl.UserId, &mdl.PubKey, &mdl.PrivKey, &mdl.Passphrase, &mdl.Exchange)
		if err != nil {
			return nil, err
		}
		apiKeysList = append(apiKeysList, mdl)
	}

	return apiKeysList, rows.Err()
}

func (a *ApiKeyRepo) GetByApiKeyByTgUserId(ctx context.Context, tgId int64) (apiKeys []*domain.ApiKeys, err error) {
	query := `SELECT a.user_id, a.exchange, a. pub_key, a.priv_key, a.passphrase 
			  FROM finder_questions_answers a 
              INNER JOIN users u ON a.user_id = u.id
			  WHERE a.tg_id = $1 LIMIT 1`

	rows, err := a.db.QueryContext(ctx, query, tgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	apiKeysList := make([]*domain.ApiKeys, 0)
	for rows.Next() {
		var mdl *domain.ApiKeys
		err = rows.Scan(&mdl.UserId, &mdl.PubKey, &mdl.PrivKey, &mdl.Passphrase, &mdl.Exchange)
		if err != nil {
			return nil, err
		}
		apiKeysList = append(apiKeysList, mdl)
	}

	return apiKeysList, rows.Err()
}
