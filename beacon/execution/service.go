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

package execution

import (
	"context"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/itsdevbear/bolaris/primitives"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// Service is responsible for delivering beacon chain notifications to
// the execution client and processing logs received from the execution client.
type Service struct {
	service.BaseService
	// engine gives the notifier access to the engine api of the execution
	// client.
	engine *engineclient.EngineClient
	// logFactory is the factory for creating objects from Ethereum logs.
	logFactory LogFactory
}

// Start spawns any goroutines required by the service.
func (s *Service) Start(ctx context.Context) {
	go s.engine.Start(ctx)
}

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error {
	return s.engine.Status()
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
// TODO: handle the bools better i.e attrs, retry, async.
func (s *Service) NotifyForkchoiceUpdate(
	ctx context.Context, fcuConfig *FCUConfig,
) (*enginetypes.PayloadID, error) {
	var (
		err       error
		payloadID *enginetypes.PayloadID
	)
	// Push the forkchoice request to the forkchoice dispatcher, we want to
	// block until
	if e := s.GCD().GetQueue(forkchoiceDispatchQueue).Sync(func() {
		payloadID, err = s.notifyForkchoiceUpdate(ctx, fcuConfig)
	}); e != nil {
		return nil, e
	}

	return payloadID, err
}

// GetPayload returns the payload and blobs bundle for the given slot.
func (s *Service) GetPayload(
	ctx context.Context, payloadID enginetypes.PayloadID, slot primitives.Slot,
) (enginetypes.ExecutionPayload, *enginetypes.BlobsBundleV1, bool, error) {
	return s.engine.GetPayload(
		ctx, payloadID,
		s.BeaconCfg().ActiveForkVersion(slot),
	)
}

// NotifyNewPayload notifies the execution client of a new payload.
// It returns true if the EL has returned VALID for the block.
func (s *Service) NotifyNewPayload(
	ctx context.Context,
	slot primitives.Slot,
	payload enginetypes.ExecutionPayload,
	versionedHashes []common.Hash,
	parentBlockRoot [32]byte,
) (bool, error) {
	return s.notifyNewPayload(
		ctx,
		slot,
		payload,
		versionedHashes,
		parentBlockRoot,
	)
}

// ProcessLogsInETH1Block gets logs in the Eth1 block
// received from the execution client and uses LogFactory to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service) ProcessLogsInETH1Block(
	ctx context.Context,
	blockHash common.Hash,
) ([]*reflect.Value, error) {
	// Gather all the logs corresponding to
	// the addresses of interest from this block.
	logsInBlock, err := s.engine.GetLogs(
		ctx,
		blockHash,
		s.logFactory.GetRegisteredAddresses(),
	)
	if err != nil {
		return nil, err
	}
	return s.logFactory.ProcessLogs(logsInBlock, blockHash)
}
