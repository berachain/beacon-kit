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
	"sync"
	"time"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/itsdevbear/bolaris/engine/client"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/sourcegraph/conc"
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

	isSyncedLock      *sync.RWMutex
	isSyncedCond      *sync.Cond
	waitForSyncStopCh chan struct{}
	isELSynced        bool
	isCLSynced        bool
}

// Start initiates the synchronization service.
func (s *Service) Start(ctx context.Context) {
	s.isSyncedLock = &sync.RWMutex{}
	s.isSyncedCond = &sync.Cond{L: s.isSyncedLock}
	s.waitForSyncStopCh = make(chan struct{})

	// Start the synchronization loop in a new goroutine.
	go func() {
		// Call syncLoop to continuously check and update the sync status.
		s.syncLoop(ctx)
		// Once the context is done, close
		// the notifySyncSignal channel to signal completion.
		<-ctx.Done()
		close(s.waitForSyncStopCh)
		s.isSyncedCond.Broadcast()
	}()
}

// Status returns the current synchronization status.
func (s *Service) Status() error {
	s.isSyncedLock.RLock()
	defer s.isSyncedLock.RUnlock()
	if err := s.status(); err != nil {
		return err
	}
	return s.BaseService.Status()
}

// status returns the current synchronization status.
func (s *Service) status() error {
	switch {
	case !s.isCLSynced:
		return ErrConsensusClientIsSyncing
	case !s.isELSynced:
		return ErrExecutionClientIsSyncing
	default:
		return nil
	}
}

// SetClientContext sets the client context for the service.
func (s *Service) SetClientContext(clientCtx cosmosclient.Context) {
	s.clientCtx = &clientCtx
}

// WaitForHealthy waits for the service to be synced.
func (s *Service) WaitForHealthy(ctx context.Context) {
	s.isSyncedCond.L.Lock()
	defer s.isSyncedCond.L.Unlock()

	// If we are not sync'd we immediately request sync progress.
	if s.status() != nil {
		go s.requestSyncProgress(ctx)
		for s.status() != nil {
			select {
			case <-s.waitForSyncStopCh:
				return
			case <-ctx.Done():
				return
			default:
				// Then we wait until we are sync'd.
				s.isSyncedCond.Wait()
			}
		}
	}
}

// syncLoop continuously runs and reports if our client is out of sync.
func (s *Service) syncLoop(ctx context.Context) {
	ticker := time.NewTicker(syncLoopInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.requestSyncProgress(ctx)
		}
	}
}

// requestSyncProgress gets the sync progress from the consensus and execution
// clients.
func (s *Service) requestSyncProgress(ctx context.Context) {
	wg := conc.NewWaitGroup()
	wg.Go(func() { s.CheckCLSync(ctx) })
	wg.Go(func() { s.CheckELSync(ctx) })
	wg.Wait()
	if s.isCLSynced && s.isELSynced {
		s.isSyncedCond.Broadcast()
	}
}
