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
	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	builderconfig "github.com/berachain/beacon-kit/mod/runtime/services/builder/config"
)

// BaseService is a base service that provides common functionality for all
// services.
type BaseService struct {
	bsb       BeaconStorageBackend
	name      string
	cfg       *config.Config
	chainSpec params.ChainSpec
	logger    log.Logger

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
	chainSpec params.ChainSpec,
	logger log.Logger,
) *BaseService {
	return &BaseService{
		bsb:       bsp,
		logger:    logger,
		cfg:       cfg,
		chainSpec: chainSpec,
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

// AvailabilityStore returns the availability store from the BaseService.
func (s *BaseService) AvailabilityStore(
	ctx context.Context,
) core.AvailabilityStore {
	return s.bsb.AvailabilityStore(ctx)
}

// BeaconState returns the beacon state from the BaseService.
func (s *BaseService) BeaconState(ctx context.Context) state.BeaconState {
	return s.bsb.BeaconState(ctx)
}

// ChainSpec returns the configuration settings of the beacon node from
// the BaseService.
func (s *BaseService) ChainSpec() params.ChainSpec {
	return s.chainSpec
}

// BuilderCfg returns the configuration settings of the builder from
// the BaseService.
func (s *BaseService) BuilderCfg() *builderconfig.Config {
	return &s.cfg.Builder
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
