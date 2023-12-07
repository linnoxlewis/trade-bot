package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/internal/domain"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (user *domain.User, err error) {
	query := `SELECT id,tg_id FROM users WHERE id = $1`
	if err = u.db.QueryRow(query, id).Scan(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return
}

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.User) (err error) {
	query := `INSERT INTO users (id,tg_id, tg_username, created_at) VALUES ($1, $2, $3,$4)`
	_, err = u.db.ExecContext(ctx,
		query,
		user.Id.String(), user.TgId, user.Username, time.Now())

	return
}
