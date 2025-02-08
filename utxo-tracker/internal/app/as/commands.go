package as

import (
	"context"
	"fmt"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
)

type NewAccountCmd struct {
	domain.Account
}

func (NewAccountCmd) ID() domain.CommandID {
	return "NewAccount"
}

type NewAccountHandler struct {
	Repo WriteRepo
}

func (h NewAccountHandler) Handle(ctx context.Context, c NewAccountCmd) error {
	// TODO Try to find the addresses
	return h.Repo.CreateUserAndAccount(ctx, c.Account)
}

type DeleteAccountCmd struct {
	AccountName domain.AccountName
	UserName    domain.UserName
}

func (DeleteAccountCmd) ID() domain.CommandID {
	return "DeleteAccount"
}

type DeleteAccountHandler struct {
	Repo WriteRepo
}

func (h DeleteAccountHandler) Handle(ctx context.Context, c DeleteAccountCmd) error {
	found, err := h.Repo.DeleteAccount(ctx, c.UserName, c.AccountName)
	if err != nil {
		return fmt.Errorf(
			"could not delete account %s for user %s: %w",
			c.AccountName, c.UserName, err,
		)
	}
	if !found {
		return fmt.Errorf(
			"could not delete account %s for user %s: %w",
			c.AccountName, c.UserName, err,
		)
	}
	return nil
}
