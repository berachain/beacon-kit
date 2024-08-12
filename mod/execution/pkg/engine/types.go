// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package engine

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type EngineClient[
	ExecutionPayloadT any,
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
] interface {
	Start(ctx context.Context) error
	GetPayload(
		ctx context.Context,
		payloadID engineprimitives.PayloadID,
		forkVersion uint32,
	) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error)
	NewPayload(
		ctx context.Context,
		payload ExecutionPayloadT,
		versionedHashes []common.ExecutionHash,
		parentBeaconBlockRoot *common.Root,
	) (*common.ExecutionHash, error)
	ForkchoiceUpdated(
		ctx context.Context,
		state *engineprimitives.ForkchoiceStateV1,
		attrs PayloadAttributesT,
		forkVersion uint32,
	) (PayloadIDT, *common.ExecutionHash, error)
}

// ExecutionPayload represents the payload of an execution block.
type ExecutionPayload[ExecutionPayloadT, WithdrawalsT any] interface {
	constraints.EngineType[ExecutionPayloadT]
	GetPrevRandao() common.Bytes32
	GetBlockHash() common.ExecutionHash
	GetParentHash() common.ExecutionHash
	GetNumber() math.U64
	GetGasLimit() math.U64
	GetGasUsed() math.U64
	GetTimestamp() math.U64
	GetExtraData() []byte
	GetBaseFeePerGas() *math.U256
	GetFeeRecipient() common.ExecutionAddress
	GetStateRoot() common.Bytes32
	GetReceiptsRoot() common.Bytes32
	GetLogsBloom() bytes.B256
	GetBlobGasUsed() math.U64
	GetExcessBlobGas() math.U64
	GetWithdrawals() WithdrawalsT
	GetTransactions() engineprimitives.Transactions
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// IncrementCounter increments a counter metric identified by the provided
	// keys.
	IncrementCounter(key string, args ...string)
}

// Withdrawal is the interface for a withdrawal.
type Withdrawal[WithdrawalT any] interface {
	// GetAmount returns the amount of the withdrawal.
	GetAmount() math.Gwei
	// GetIndex returns the public key of the validator.
	GetIndex() math.U64
	// GetValidatorIndex returns the index of the validator.
	GetValidatorIndex() math.ValidatorIndex
	// GetAddress returns the address of the withdrawal.
	GetAddress() common.ExecutionAddress
}
