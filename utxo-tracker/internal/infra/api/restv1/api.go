package restv1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hannesdejager/utxo-tracker/internal/app/as"
	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/hannesdejager/utxo-tracker/internal/infra/jaeger"
	"github.com/hannesdejager/utxo-tracker/internal/infra/logging"
	"github.com/hannesdejager/utxo-tracker/internal/infra/prometheus"
	"github.com/hannesdejager/utxo-tracker/pkg/gons/restdocs"
)

//go:generate ../../../../scripts/gen-rest-v1-api.sh

func NewHandler(log *slog.Logger, baseURL string, rr as.ReadRepo, wr as.WriteRepo) http.Handler {
	r := chi.NewRouter()
	r.Use(prometheus.APIMiddleware)
	r.Use(jaeger.TracingMiddleware)
	r.Use(logging.APIRequestLogger(log))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(logging.Recoverer(log))
	r.Get(baseURL+"/spec", SpecHandler())
	r.Mount(baseURL+"/docs", restdocs.New(
		restdocs.WithSpecURL(baseURL+"/spec"),
		restdocs.WithBaseURL(baseURL+"/docs"),
		restdocs.WithHtmlTitle("Account Service API"),
	).Handler())
	r.Get(baseURL, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURL+"/docs", http.StatusMovedPermanently)
	})
	return HandlerFromMuxWithBaseURL(
		&impl{rr: rr, wr: wr, log: log},
		r,
		baseURL,
	)
}

// impl is our implementation of the ServerInterface that was generated.
type impl struct {
	rr  as.ReadRepo
	wr  as.WriteRepo
	log *slog.Logger
}

// GetAccounts returns a 501 status code indicating that the functionality is not implemented.
func (s *impl) GetAccounts(
	w http.ResponseWriter,
	r *http.Request,
	params GetAccountsParams,
) {
	accounts := make([]Account, 0)
	response := GetAccountsResponse{
		Accounts: &accounts,
	}
	err := s.rr.GetAccountsWithAddresses(
		r.Context(),
		domain.UserName(params.XUserID),
		func(a domain.AccountAddresses) error {
			acct := Account{
				Id:        string(a.Account),
				Name:      string(a.Account),
				Addresses: make([]string, len(a.Items)),
			}
			for i, from := range a.Items {
				acct.Addresses[i] = string(from)
			}
			accounts = append(accounts, acct)
			return nil
		},
	)

	if err != nil {
		s.log.ErrorContext(r.Context(), "error fetching accounts", "error", err)
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}

	// Set content type and return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(response); err != nil {
		s.log.ErrorContext(r.Context(), "could not encode response as JSON", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateAccount creates an account.
func (s *impl) CreateAccount(w http.ResponseWriter, r *http.Request, params CreateAccountParams) {
	var payload NewAccountRequest
	d := json.NewDecoder(r.Body)
	err := d.Decode(&payload)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = s.wr.CreateUserAndAccount(r.Context(), domain.Account{
		User: domain.UserName(params.XUserID),
		ID:   domain.AccountName(payload.Name),
		XPub: domain.XPubEtc(payload.Xpub),
		Type: domain.AccountType(payload.Type),
	})
	if err != nil {
		if s.wr.IsDuplicateError(err) {
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}
		s.log.ErrorContext(r.Context(), "Could not save account", "error", fmt.Errorf("%v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Location", fmt.Sprintf("/accounts/%s", payload.Name))
}

// GetAccountById returns a 501 status code indicating that the functionality is not implemented.
func (s *impl) GetAccountById(w http.ResponseWriter, r *http.Request, accountId string, params GetAccountByIdParams) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// Delete an account by ID
// (DELETE /accounts/{accountId})
func (s *impl) DeleteAccount(w http.ResponseWriter, r *http.Request, accountId string, params DeleteAccountParams) {
	wasDeleted, err := s.wr.DeleteAccount(r.Context(), domain.UserName(params.XUserID), domain.AccountName(accountId))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if !wasDeleted {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
