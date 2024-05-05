package main

import (
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner/pruner"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize your database here. This could be a filedb.DB or pebble.DB instance.
	var db interfaces.DB

	// TODO: Time could be configured
	interval := time.Minute * 1

	// Create a new pruner instance.

	p := pruner.NewPruner(interval, db)

	go p.Start()
	// Start the pruner in a separate goroutine.
	//go func() {
	//	for {
	//		select {
	//		case <-time.After(interval):
	//			fmt.Println("starting pruner")
	//			err := p.Start()
	//			if err != nil {
	//				return
	//			}
	//		}
	//	}
	//}()

	// Wait for an interrupt signal.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// When you're done, stop the pruner.
	defer p.Stop()
}
