package logging

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/hannesdejager/utxo-tracker/pkg/gons/middleware"
)

func APIRequestLogger(
	log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter,
			r *http.Request) {
			peeker := &middleware.ResponsePeeker{
				ResponseWriter: w}

			defer func() {
				var l = log.InfoContext
				if peeker.Status >= 500 {
					l = log.ErrorContext
				}
				l(r.Context(),
					"API Request made", "path", r.URL.Path, "status",
					peeker.Status)
			}()

			next.ServeHTTP(peeker, r.WithContext(r.Context()))
		})
	}
}

func Recoverer(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					log.ErrorContext(r.Context(),
						"Panic ocurred in HTTP handler",
						"error", fmt.Errorf("%v", rvr),
						"stack_trace", string(debug.Stack()))

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
