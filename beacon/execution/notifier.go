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

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/types"
	"github.com/itsdevbear/bolaris/types/config"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

type forkchoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

type EngineNotifier struct {
	engine.Caller
	payloadCache *cache.ProposerPayloadIDsCache
	beaconCfg    *config.Beacon
	etherbase    common.Address
	logger       log.Logger
	fcsp         forkchoiceStoreProvider
	engine       engine.Caller
}

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

func (s *EngineNotifier) NotifyForkchoiceUpdate(
	ctx context.Context, slot primitives.Slot, arg *NotifyForkchoiceUpdateArg,
	withAttrs, withRetry bool,
) (*enginev1.PayloadIDBytes, error) {
	if withRetry {
		return s.notifyForkchoiceUpdateWithSyncingRetry(ctx, slot, arg, withAttrs)
	}
	return s.notifyForkchoiceUpdate(ctx, slot, arg, withAttrs)
}

// It returns true if the EL has returned VALID for the block.
func (s *EngineNotifier) NotifyNewPayload(ctx context.Context /*preStateVersion*/, _ int,
	preStateHeader interfaces.ExecutionData, /*, blk interfaces.ReadOnlySignedBeaconBlock*/
) (bool, error) {
	lastValidHash, err := s.engine.NewPayload(ctx, preStateHeader,
		[]common.Hash{}, &common.Hash{} /*empty version hashes and root before Deneb*/)
	return lastValidHash != nil, err
}
