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

package types

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	eip4844 "github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// RawBeaconBlockBody is the interface for a beacon block body.
type RawBeaconBlockBody interface {
	WriteOnlyBeaconBlockBody
	ReadOnlyBeaconBlockBody
	Length() uint64
}

// WriteOnlyBeaconBlockBody is the interface for a write-only beacon block body.
type WriteOnlyBeaconBlockBody interface {
	SetDeposits([]*types.Deposit)
	SetEth1Data(*types.Eth1Data)
	SetExecutionData(*types.ExecutionPayload) error
	SetBlobKzgCommitments(eip4844.KZGCommitments[gethprimitives.ExecutionHash])
	SetRandaoReveal(crypto.BLSSignature)
	SetGraffiti(common.Bytes32)
	SetAttestations([]*types.AttestationData)
	SetSlashingInfo([]*types.SlashingInfo)
}

// ReadOnlyBeaconBlockBody is the interface for
// a read-only beacon block body.
type ReadOnlyBeaconBlockBody interface {
	constraints.SSZMarshallableRootable
	IsNil() bool

	// Execution returns the execution data of the block.
	GetDeposits() []*types.Deposit
	GetEth1Data() *types.Eth1Data
	GetGraffiti() common.Bytes32
	GetRandaoReveal() crypto.BLSSignature
	GetExecutionPayload() *types.ExecutionPayload
	GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash]
	GetTopLevelRoots() ([][32]byte, error)
}

// RawBeaconBlock is the interface for a beacon block.
type RawBeaconBlock[BeaconBlockBodyT RawBeaconBlockBody] interface {
	constraints.SSZMarshallableRootable
	SetStateRoot(common.Root)
	GetStateRoot() common.Root
	IsNil() bool
	Version() uint32
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetBody() BeaconBlockBodyT
	GetHeader() *types.BeaconBlockHeader
}

// executionPayloadBody is the interface for the execution data of a block.
type executionPayloadBody interface {
	constraints.SSZMarshallableRootable
	constraints.JSONMarshallable
	IsNil() bool
	Version() uint32
	GetPrevRandao() common.Bytes32
	GetBlockHash() gethprimitives.ExecutionHash
	GetParentHash() gethprimitives.ExecutionHash
	GetNumber() math.U64
	GetGasLimit() math.U64
	GetGasUsed() math.U64
	GetTimestamp() math.U64
	GetExtraData() []byte
	GetBaseFeePerGas() math.Wei
	GetFeeRecipient() gethprimitives.ExecutionAddress
	GetStateRoot() common.Bytes32
	GetReceiptsRoot() common.Bytes32
	GetLogsBloom() []byte
	GetBlobGasUsed() math.U64
	GetExcessBlobGas() math.U64
}
