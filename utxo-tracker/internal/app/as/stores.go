// Package as contains the application logic for the account service
package as

import (
	"context"
	"errors"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
)

// WriteRepo defines the requirements for the database where we will write
// account data.
type WriteRepo interface {
	// CreateUserAndAccount ensures the user exists and creates an account.
	CreateUserAndAccount(context.Context, domain.Account) error

	// DeleteAccount deletes an account for a user. If no error it will return
	// true if the account was indeed deleted or false if no such account was found.
	DeleteAccount(context.Context, domain.UserName, domain.AccountName) (bool, error)

	// Tells if the given error is because of a clash in
	// the database i.e. the account with that ID already exists
	IsDuplicateError(error) bool
}

// ErrIterationExit can be returned by a QueryCallback to prematurely end iteration
var ErrIterationExit = errors.New("query: early exit of iteration")

// WriteRepo defines the requirements for the database where we will read
// account data from.
type ReadRepo interface {
	// GetAccountsWithAddresses returns all accounts and also their addresses. It calls the
	// callback function repeatedly for each account. If the callback returns an error
	// then GetAccountsWithAddresses exits without going through all records
	// and it returns the error of the callback.
	GetAccountsWithAddresses(context.Context, domain.UserName, func(domain.AccountAddresses) error) error

	//GetAccountDetails(context.Context, domain.UserName, domain.AccountName) (domain.AccountAddresses, error)

	// Tells if the given error is because of a record that was not found.
	IsNotFoundError(e error) bool
}
