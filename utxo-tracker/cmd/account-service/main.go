package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/infra/api/restv1"
	"github.com/hannesdejager/utxo-tracker/internal/infra/env"
	"github.com/hannesdejager/utxo-tracker/internal/infra/httpsvr"
	"github.com/hannesdejager/utxo-tracker/internal/infra/sys"
)

func main() {
	fmt.Printf("Starting, PID=%d\n", os.Getpid())
	svr := httpsvr.StartAsync(
		env.HTTPConfig(),
		restv1.NewHandler("/rest/v1"))
	sys.AwaitTermination()
	log.Println("Shutting down...")
	httpsvr.StopGracefully(svr, 30*time.Second)
	log.Println("Bye!")
}
