package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.User) (err error) {
	query := `INSERT INTO users (id, username, created_at) VALUES ($1, $2, $3)`
	_, err = u.db.ExecContext(ctx,
		query,
		user.ID, user.Username, time.Now())

	return
}

func (u *UserRepository) ExistUser(ctx context.Context, id int64) (exist bool) {
	query := `SELECT EXISTS( SELECT 1 FROM users where id = $1)`
	if err := u.db.QueryRowContext(ctx, query, id).Scan(&exist); err != nil {
		return false
	}

	return
}

func (u *UserRepository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var exist bool
	query := `SELECT EXISTS( SELECT 1 FROM users where id = $1 AND is_admin = $2)`
	if err := u.db.QueryRowContext(ctx, query, id, true).Scan(&exist); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exist, nil
}

func (u *UserRepository) GetAdminIds(ctx context.Context) ([]int, error) {
	query := `SELECT id FROM users WHERE is_admin = $1`
	rows, err := u.db.QueryContext(ctx, query, true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}
