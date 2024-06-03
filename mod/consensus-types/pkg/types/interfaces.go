// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	eip4844 "github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// BeaconBlockBody is the interface for a beacon block body.
type BeaconBlockBody interface {
	WriteOnlyBeaconBlockBody
	ReadOnlyBeaconBlockBody
	Length() uint64
}

// WriteOnlyBeaconBlockBody is the interface for a write-only beacon block body.
type WriteOnlyBeaconBlockBody interface {
	SetDeposits([]*Deposit)
	SetEth1Data(*Eth1Data)
	SetExecutionData(*ExecutionPayload) error
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
	SetRandaoReveal(crypto.BLSSignature)
}

// ReadOnlyBeaconBlockBody is the interface for
// a read-only beacon block body.
type ReadOnlyBeaconBlockBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool

	// Execution returns the execution data of the block.
	GetDeposits() []*Deposit
	GetEth1Data() *Eth1Data
	GetGraffiti() bytes.B32
	GetRandaoReveal() crypto.BLSSignature
	GetExecutionPayload() *ExecutionPayload
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
	GetTopLevelRoots() ([][32]byte, error)
}

// BeaconBlock is the interface for a beacon block.
type RawBeaconBlock[BeaconBlockBodyT BeaconBlockBody] interface {
	SetStateRoot(common.Root)
	GetStateRoot() common.Root
	ReadOnlyBeaconBlock[BeaconBlockBodyT]
}

type BeaconBlockG[BodyT any] struct {
	ReadOnlyBeaconBlock[BodyT]
}

// ReadOnlyBeaconBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconBlock[BodyT any] interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Version() uint32
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	GetBody() BodyT
	GetHeader() *BeaconBlockHeader
}
