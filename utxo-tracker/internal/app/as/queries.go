package as

import (
	"context"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
)

type GetAccountsQuery struct {
	User domain.UserName
}

// ID is an implementation for domain.Query
func (q GetAccountsQuery) ID() domain.QueryID {
	return "GetAccounts"
}

// GetAccountsQueryHandler handles domain.Query's of type GetAccountsQuery
type GetAccountsQueryHandler struct {
	Repo ReadRepo
}

// Handle is an implementation of domain.QueryHandler
func (h *GetAccountsQueryHandler) Handle(ctx context.Context, q GetAccountsQuery, callback domain.QueryCallback[domain.AccountAddresses]) error {
	return h.Repo.GetAccountsWithAddresses(ctx, q.User, func(aa domain.AccountAddresses) error {
		return callback(aa)
	})
}

type GetSingleAccountQuery struct {
	User domain.UserName
}

// ID is an implementation for domain.Query
func (q GetSingleAccountQuery) ID() domain.QueryID {
	return "GetSingleAccount"
}

// GetSingleAccountQueryHandler handles domain.Query's of type GetSingleAccountQuery.
type GetSingleAccountQueryHandler struct {
	Repo ReadRepo
}

// Handle is an implementation of domain.QueryHandler
func (h *GetSingleAccountQueryHandler) Handle(ctx context.Context, q GetSingleAccountQuery, callback domain.QueryCallback[domain.AccountAddresses]) error {
	return h.Repo.GetAccountsWithAddresses(ctx, q.User, func(aa domain.AccountAddresses) error {
		return callback(aa)
	})
}
