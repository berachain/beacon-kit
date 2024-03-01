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
	"github.com/itsdevbear/bolaris/async/dispatch"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// BaseService is a base service that provides common functionality for all
// services.
type BaseService struct {
	BeaconStorageBackend
	name   string
	cfg    *config.Config
	gcd    *dispatch.GrandCentralDispatch
	logger log.Logger

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
		BeaconStorageBackend: bsp,
		gcd:                  gcd,
		logger:               logger,
		cfg:                  cfg,
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

// BeaconCfg returns the configuration settings of the beacon node from
// the BaseService. It provides access to various configuration parameters
// used by the beacon node.
func (s *BaseService) BeaconCfg() *config.Beacon {
	return &s.cfg.Beacon
}

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
func (s *BaseService) ActiveForkVersionForSlot(slot primitives.Slot) int {
	return s.BeaconCfg().ActiveForkVersion(primitives.Epoch(slot))
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
