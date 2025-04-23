// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core

import (
	"context"

	"github.com/berachain/beacon-kit/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

type ReadOnlyBeaconState interface {
	GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
	GetSlot() (math.Slot, error)
	GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
}

// ReadOnlyContext defines an interface for managing state transition context.
type ReadOnlyContext interface {
	ConsensusCtx() context.Context
	ConsensusTime() math.U64
	ProposerAddress() []byte
	VerifyPayload() bool
	VerifyRandao() bool
	VerifyResult() bool
	MeterGas() bool
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// NotifyNewPayload notifies the execution client of the new payload.
	NotifyNewPayload(
		ctx context.Context,
		req ctypes.NewPayloadRequest,
		retryOnSyncingStatus bool,
	) error
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	SetGauge(key string, value int64, args ...string)
	// IncrementCounter increments the counter identified by
	// the provided key.
	IncrementCounter(key string, args ...string)
}

type ChainSpec interface {
	chain.HysteresisSpec
	chain.BalancesSpec
	chain.DepositSpec
	chain.ForkSpec
	chain.DomainTypeSpec
	chain.WithdrawalsSpec
	SlotsPerEpoch() uint64
	SlotToEpoch(slot math.Slot) math.Epoch
	SlotsPerHistoricalRoot() uint64
	EpochsPerHistoricalVector() uint64
	GenesisForkVersion() common.Version
	ActiveForkVersionForTimestamp(timestamp math.U64) common.Version
	ValidatorSetCap() uint64
	HistoricalRootsLimit() uint64
}
