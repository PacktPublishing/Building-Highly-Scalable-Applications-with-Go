package httpsvr

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/app/config"
)

func Routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello, World!"))
		})
	return router
}

// Start sets up and starts the HTTP server
func StartAsync(c config.HTTPServer, h http.Handler) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", c.Port),
		Handler:           h,
		ReadTimeout:       c.ReadTimeout,
		WriteTimeout:      c.WriteTimeout,
		IdleTimeout:       c.IdleTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf(
				"Error listening on port %d: %v\n",
				c.Port,
				err,
			)
		}
	}()

	return server
}

func StopGracefully(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf(
			"Server forced to shutdown: %v",
			err,
		)
	}
}
