// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package sync

import (
	"context"
	"time"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/itsdevbear/bolaris/engine/client"
	"github.com/itsdevbear/bolaris/runtime/service"
	"golang.org/x/sync/errgroup"
)

// syncLoopInterval is the interval at which the sync loop checks for sync
// status.
//

const syncLoopInterval = 6 * time.Second

// Service is responsible for tracking the synchornization status
// of both the beacon and execution chains.
type Service struct {
	service.BaseService
	engineClient *client.EngineClient
	clientCtx    *cosmosclient.Context
}

// SetClientContext sets the client context for the service.
func (s *Service) SetClientContext(clientCtx cosmosclient.Context) {
	s.clientCtx = &clientCtx
}

// Start initiates the synchronization service.
func (s *Service) Start(ctx context.Context) {
	// Start the synchronization loop in a new goroutine.
	go func() {
		// Call syncLoop to continuously check and update the sync status.
		s.syncLoop(ctx)
		// Once the context is done, close
		// the notifySyncSignal channel to signal completion.
		<-ctx.Done()
	}()
}

// syncLoop continuously runs and reports if our client is out of sync.
func (s *Service) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(syncLoopInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.SetStatus(ErrNotRunning)
			return
		case <-ticker.C:
			s.SetStatus(s.RequestSyncProgress(ctx))
		}
	}
}

// Request the sync progress from the consensus and execution clients.
func (s *Service) RequestSyncProgress(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Query Consensus Client for sync progress.
	g.Go(func() error {
		return s.CheckCLSync(ctx)
	})

	// Query Execution Client for sync progress.
	g.Go(func() error {
		return s.CheckELSync(ctx)
	})

	return g.Wait()
}
