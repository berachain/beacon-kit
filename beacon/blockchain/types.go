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

package blockchain

import (
	"context"
	"reflect"

	"github.com/berachain/beacon-kit/beacon/execution"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

// LocalBuilder is the interface for the builder service.
type LocalBuilder interface {
	BuildLocalPayload(
		ctx context.Context,
		parentEth1Hash primitives.ExecutionHash,
		slot primitives.Slot,
		timestamp uint64,
		parentBlockRoot [32]byte,
	) (*enginetypes.PayloadID, error)
}

type ExecutionService interface {
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		fcuConfig *execution.FCUConfig,
	) (*enginetypes.PayloadID, error)

	// NotifyNewPayload notifies the execution client of a new payload.
	NotifyNewPayload(
		ctx context.Context,
		slot primitives.Slot,
		payload enginetypes.ExecutionPayload,
		versionedHashes []primitives.ExecutionHash,
		parentBlockRoot [32]byte,
	) (bool, error)

	ProcessLogsInETH1Block(
		ctx context.Context,
		blockHash primitives.ExecutionHash,
	) ([]*reflect.Value, error)
}

type StakingService interface{}

type SyncService interface {
	IsInitSync() bool
	Status() error
}
