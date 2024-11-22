package restv1

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hannesdejager/utxo-tracker/internal/infra/jaeger"
	"github.com/hannesdejager/utxo-tracker/internal/infra/logging"
	"github.com/hannesdejager/utxo-tracker/internal/infra/prometheus"
	"github.com/hannesdejager/utxo-tracker/pkg/gons/restdocs"
)

//go:generate ../../../../scripts/gen-rest-v1-api.sh

func NewHandler(log *slog.Logger, baseURL string) http.Handler {
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
		&impl{},
		r,
		baseURL,
	)
}

// AccountService is our implementation of the ServerInterface.
type impl struct{}

// GetAccounts returns a 501 status code indicating that the functionality is not implemented.
func (s *impl) GetAccounts(
	w http.ResponseWriter,
	r *http.Request,
	params GetAccountsParams,
) {
	http.Error(
		w,
		"Not Implemented",
		http.StatusNotImplemented,
	)
}

// CreateAccount returns a 501 status code indicating that the functionality is not implemented.
func (s *impl) CreateAccount(w http.ResponseWriter, r *http.Request, params CreateAccountParams) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// GetAccountById returns a 501 status code indicating that the functionality is not implemented.
func (s *impl) GetAccountById(w http.ResponseWriter, r *http.Request, accountId string, params GetAccountByIdParams) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
