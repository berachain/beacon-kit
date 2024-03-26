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

package service

import (
	"context"
	"sync"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/async/dispatch"
	"github.com/berachain/beacon-kit/beacon/core/state"
	"github.com/berachain/beacon-kit/beacon/forkchoice"
	"github.com/berachain/beacon-kit/beacon/forkchoice/ssf"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/config/params"
	"github.com/berachain/beacon-kit/primitives"
)

// BaseService is a base service that provides common functionality for all
// services.
type BaseService struct {
	bsb    BeaconStorageBackend
	name   string
	cfg    *config.Config
	gcd    *dispatch.GrandCentralDispatch
	logger log.Logger
	fcr    forkchoice.ForkChoicer

	// statusErrMu protects statusErr.
	statusErrMu *sync.RWMutex
	// statusErr is the error returned
	// by the last status check.
	statusErr error
}

// NewBaseService creates a new BaseService and applies the provided options.
func NewBaseService(
	cfg *config.Config,
	bsp BeaconStorageBackend,
	gcd *dispatch.GrandCentralDispatch,
	logger log.Logger,
) *BaseService {
	return &BaseService{
		bsb:    bsp,
		gcd:    gcd,
		logger: logger,
		cfg:    cfg,
		fcr:    ssf.New(bsp.ForkchoiceStore(context.Background())),
	}
}

// Name returns the name of the BaseService.
func (s *BaseService) Name() string {
	return s.name
}

// Logger returns the logger instance of the BaseService.
// It is used for logging messages in a structured manner.
func (s *BaseService) Logger() log.Logger {
	return s.logger
}

// GCD returns the GrandCentralDispatch instance of the BaseService.
// It is used for managing asynchronous tasks and dispatching events.
func (s *BaseService) GCD() *dispatch.GrandCentralDispatch {
	return s.gcd
}

// AvailabilityStore returns the availability store from the BaseService.
func (s *BaseService) AvailabilityStore(
	ctx context.Context,
) state.AvailabilityStore {
	return s.bsb.AvailabilityStore(ctx)
}

// BeaconState returns the beacon state from the BaseService.
func (s *BaseService) BeaconState(ctx context.Context) state.BeaconState {
	return s.bsb.BeaconState(ctx)
}

// ForkchoiceStore returns the forkchoice store from the BaseService.
func (s *BaseService) ForkchoiceStore(
	ctx context.Context,
) forkchoice.ForkChoicer {
	// TODO: Decouple from the Specific SingleSlotFinalityStore Impl.
	// TODO: SetContext isn't consistent with the rest of the methods.
	s.fcr.SetContext(ctx)
	return s.fcr
}

// BeaconCfg returns the configuration settings of the beacon node from
// the BaseService.
func (s *BaseService) BeaconCfg() *params.BeaconChainConfig {
	return &s.cfg.Beacon
}

// BuilderCfg returns the configuration settings of the builder from
// the BaseService.
func (s *BaseService) BuilderCfg() *config.Builder {
	return &s.cfg.Builder
}

// FeatureFlags returns the feature flags from the BaseService.
func (s *BaseService) FeatureFlags() *config.FeatureFlags {
	return &s.cfg.FeatureFlags
}

// Start is an intentional no-op for the BaseService.
func (s *BaseService) Start(context.Context) {}

// Status is an intentional no-op for the BaseService.
func (s *BaseService) Status() error {
	s.statusErrMu.RLock()
	defer s.statusErrMu.RUnlock()
	return s.statusErr
}

// WaitForHealthy is an intentional no-op for the BaseService.
func (s *BaseService) WaitForHealthy(context.Context) {}

// SetStatus sets the status error of the BaseService.
func (s *BaseService) SetStatus(err error) {
	s.statusErrMu.Lock()
	defer s.statusErrMu.Unlock()
	s.statusErr = err
}

// ActiveForkVersionForSlot returns the active fork version for the given slot.
func (s *BaseService) ActiveForkVersionForSlot(slot primitives.Slot) uint32 {
	return s.BeaconCfg().ActiveForkVersion(slot)
}

// // DispatchEvent sends a value to the feed associated with the provided key.
// func (s *BaseService) DispatchEvent(value dispatch.Event) int {
// 	return s.gcd.Dispatch(value)
// }

// // SubscribeToEvent subscribes a channel to the feed associated with the
// provided key.
// func (s *BaseService) SubscribeToEvent(key string, eventType reflect.Type) {
// 	channel := make(chan dispatch.Event)
// 	s.channels = append(s.channels, channel)
// 	s.gcd.Subscribe(key, channel)
// }
