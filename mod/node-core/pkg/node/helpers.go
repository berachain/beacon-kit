package node

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/berachain/beacon-kit/mod/log"

	"golang.org/x/sync/errgroup"
)

// listenForQuitSignals listens for SIGINT and SIGTERM. When a signal is
// received,
// the cleanup function is called, indicating the caller can gracefully exit or
// return.
//
// Note, the blocking behavior of this depends on the block argument.
// The caller must ensure the corresponding context derived from the cancelFn is
// used correctly.
func listenForQuitSignals(
	g *errgroup.Group,
	block bool,
	cancelFn context.CancelFunc,
	logger log.Logger[any],
) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	f := func() {
		sig := <-sigCh
		cancelFn()

		logger.Info("caught signal", "signal", sig.String())
	}

	if block {
		g.Go(func() error {
			f()
			return nil
		})
	} else {
		go f()
	}
}
