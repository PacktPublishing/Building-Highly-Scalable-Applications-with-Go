package sys

import (
	"os"
	"os/signal"
	"syscall"
)

// AwaitTermination blocks until a shutdown signal is
// received from outside of the application.
func AwaitTermination() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
