package aspostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hannesdejager/utxo-tracker/internal/app/as"
	"github.com/hannesdejager/utxo-tracker/internal/app/config"
	"github.com/hannesdejager/utxo-tracker/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

type WriteRepo struct {
	db *sql.DB
}

func NewWriteRepo(c config.AccountServiceDB) (as.WriteRepo, *sql.DB, error) {
	var err error
	db, err := sql.Open("pgx", c.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool limits
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0) // No timeout

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &WriteRepo{db: db}, db, nil
}

func (r *WriteRepo) CreateUserAndAccount(ctx context.Context, a domain.Account) error {
	_, err := r.db.ExecContext(ctx, `
		WITH ins AS (
			INSERT INTO users (name) VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		)
		INSERT INTO accounts (fk_users, name, xpub, acc_type)
		SELECT id, $2, $3, $4 FROM ins;
	`, a.User, a.ID, a.XPub, a.Type)
	return err
}

func (r *WriteRepo) DeleteAccount(ctx context.Context, u domain.UserName, a domain.AccountName) (bool, error) {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM accounts
		WHERE fk_users = (SELECT id FROM users WHERE name = $1)
		AND name = $2;
	`, u, a)
	count, _ := res.RowsAffected()
	return count > 0, err
}

func (r *WriteRepo) IsDuplicateError(e error) bool {
	var pqErr *pq.Error
	if errors.As(e, &pqErr) {
		return pqErr.Code == "23505" // Unique violation
	}
	return false
}
