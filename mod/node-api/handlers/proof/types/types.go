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
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
)

// BeaconBlockHeader is the interface for a beacon block header.
type BeaconBlockHeader interface {
	constraints.SSZRootable
	// GetTree is kept for FastSSZ compatibility.
	GetTree() (*fastssz.Node, error)
	// GetProposerIndex returns the proposer index.
	GetProposerIndex() math.ValidatorIndex
}

// BeaconState is the interface for a beacon state.
type BeaconState[
	BeaconStateMarshallableT, ExecutionPayloadHeaderT, ValidatorT any,
] interface {
	// GetLatestExecutionPayloadHeader returns the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (ExecutionPayloadHeaderT, error)
	// GetMarshallable returns the marshallable version of the beacon state.
	GetMarshallable() (BeaconStateMarshallableT, error)
	// ValidatorByIndex retrieves the validator at the given index.
	ValidatorByIndex(index math.ValidatorIndex) (ValidatorT, error)
}

// BeaconStateMarshallable is the interface for a beacon state that can be
// marshalled or hash tree rooted.
type BeaconStateMarshallable interface {
	// GetTree is kept for FastSSZ compatibility.
	GetTree() (*fastssz.Node, error)
}

// ExecutionPayloadHeader is the interface for an execution payload header.
type ExecutionPayloadHeader interface {
	// GetNumber returns the block number of the ExecutionPayloadHeader.
	GetNumber() math.U64
	// GetFeeRecipient returns the fee recipient address of the
	// ExecutionPayloadHeader.
	GetFeeRecipient() gethprimitives.ExecutionAddress
}

// Validator is the interface for a validator.
type Validator interface {
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
}
