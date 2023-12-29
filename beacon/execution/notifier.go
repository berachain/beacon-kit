// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"errors"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/types/config"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// EngineNotifier is responsible for delivering beacon chain notifications to
// the execution client.
type EngineNotifier struct {
	// engine gives the notifier access to the engine api of the execution client.
	engine engine.Caller
	// payloadCache is used to track currently building payload IDs for a given slot.
	payloadCache *cache.ProposerPayloadIDsCache
	// beaconCfg is the beacon chain configuration.
	beaconCfg *config.Beacon
	// etherbase is the address to which block rewards are sent.
	etherbase common.Address
	// logger is the logger for the service.
	logger log.Logger
	// fcsp is the fork choice store provider.
	fcsp forkchoiceStoreProvider
}

// New creates a new EngineNotifier with the provided options.
func New(opts ...Option) *EngineNotifier {
	ec := &EngineNotifier{
		payloadCache: cache.NewProposerPayloadIDsCache(),
	}
	for _, opt := range opts {
		if err := opt(ec); err != nil {
			ec.logger.Error("Failed to apply option", "error", err)
		}
	}

	return ec
}

// NotifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (s *EngineNotifier) NotifyForkchoiceUpdate(
	ctx context.Context, slot primitives.Slot, arg *NotifyForkchoiceUpdateArg,
	withAttrs, withRetry bool,
) (*enginev1.PayloadIDBytes, error) {
	if withRetry {
		return s.notifyForkchoiceUpdateWithSyncingRetry(ctx, slot, arg, withAttrs)
	}
	return s.notifyForkchoiceUpdate(ctx, slot, arg, withAttrs)
}

// NotifyNewPayload notifies the execution client of a new payload.
// It returns true if the EL has returned VALID for the block.
func (s *EngineNotifier) NotifyNewPayload(ctx context.Context /*preStateVersion*/, _ int,
	preStateHeader interfaces.ExecutionData, /*, blk interfaces.ReadOnlySignedBeaconBlock*/
) (bool, error) {
	// var lastValidHash []byte
	// if blk.Version() >= version.Deneb {
	// 	var versionedHashes []common.Hash
	// 	versionedHashes, err = kzgCommitmentsToVersionedHashes(blk.Block().Body())
	// 	if err != nil {
	// 		return false, errors.Wrap(err, "could not get versioned hashes to feed the engine")
	// 	}
	// 	pr := common.Hash(blk.Block().ParentRoot())
	// 	lastValidHash, err = s.cfg.ExecutionEngineCaller.NewPayload
	//			(ctx, payload, versionedHashes, &pr)
	// } else {
	// 	lastValidHash, err = s.cfg.ExecutionEngineCaller.NewPayload(ctx, payload,
	// []common.Hash{}, &common.Hash{} /*empty version hashes and root before Deneb*/)
	// }

	lastValidHash, err := s.engine.NewPayload(ctx, preStateHeader,
		[]common.Hash{}, &common.Hash{} /*empty version hashes and root before Deneb*/)
	return lastValidHash != nil, err
}

// GetBuiltPayload returns the payload and blobs bundle for the given slot.
func (s *EngineNotifier) GetBuiltPayload(
	ctx context.Context, slot primitives.Slot,
) (interfaces.ExecutionData, *enginev1.BlobsBundle, bool, error) {
	_, payloadID, found := s.payloadCache.GetProposerPayloadIDs(
		slot, [32]byte{}, // TODO: support building on multiple heads as a safety fallback feature.
	)
	if !found {
		return nil, nil, false, errors.New("payload not found")
	}
	return s.engine.GetPayload(ctx, payloadID, slot)
}
