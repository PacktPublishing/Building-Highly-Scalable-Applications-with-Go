package aspostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hannesdejager/utxo-tracker/internal/app/as"
	"github.com/hannesdejager/utxo-tracker/internal/app/config"
	"github.com/hannesdejager/utxo-tracker/internal/domain"
)

var ErrNotFound = errors.New("record not found")

type readRepo struct {
	db *sql.DB
}

func NewReadRepo(c config.AccountServiceDB) (as.ReadRepo, *sql.DB, error) {
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

	return &readRepo{db: db}, db, nil
}

func (r *readRepo) GetAccountsWithAddresses(ctx context.Context, u domain.UserName, callback func(v domain.AccountAddresses) error) error {
	query := `
		SELECT a.id, a.name, COALESCE(aa.address, '')
		FROM accounts a
		JOIN users u ON a.fk_users = u.id
		LEFT JOIN account_addresses aa ON a.id = aa.fk_accounts
		WHERE u.name = $1
		ORDER BY a.name, aa.address;
	`

	rows, err := r.db.QueryContext(ctx, query, u)
	if err != nil {
		return fmt.Errorf("querying accounts with addresses: %w", err)
	}
	defer rows.Close()

	var currentAccount domain.AccountAddresses
	var lastAccountID int

	for rows.Next() {
		var (
			accountID int
			accName   domain.AccountName
			address   domain.Address
		)

		if err := rows.Scan(&accountID, &accName, &address); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}

		if accountID != lastAccountID {
			// Send the previous account to the callback if it's not the first iteration
			if lastAccountID != 0 {
				if err := callback(currentAccount); err != nil {
					if errors.Is(err, as.ErrIterationExit) {
						return nil // Graceful exit
					}
					return err
				}
			}

			// Start a new account
			currentAccount = domain.AccountAddresses{
				Account: accName,
				Items:   []domain.Address{},
			}
			lastAccountID = accountID
		}

		if address != "" {
			currentAccount.Items = append(currentAccount.Items, address)
		}
	}

	// Send the last account if there was any
	if lastAccountID != 0 {
		if err := callback(currentAccount); err != nil {
			if errors.Is(err, as.ErrIterationExit) {
				return nil
			}
			return err
		}
	}

	return rows.Err()
}

func (*readRepo) IsNotFoundError(e error) bool {
	return errors.Is(e, sql.ErrNoRows)
}
