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

package engineprimitives

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// Marshallable is an interface that combines the ssz.Marshaler and
// ssz.Unmarshaler interfaces.
type SSZMarshallable interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
	// SizeSSZ returns the size in bytes that the object would take when
	// marshaled.
	SizeSSZ() int
}

// ExecutionPayloadBody is the interface for the execution data of a block.
// It contains all the fields that are part of both an execution payload header
// and a full execution payload.
type ExecutionPayloadBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	json.Marshaler
	json.Unmarshaler
	IsNil() bool
	Version() uint32
	GetPrevRandao() primitives.Bytes32
	GetBlockHash() common.ExecutionHash
	GetParentHash() common.ExecutionHash
	GetNumber() math.U64
	GetGasLimit() math.U64
	GetGasUsed() math.U64
	GetTimestamp() math.U64
	GetExtraData() []byte
	GetBaseFeePerGas() math.Wei
	GetFeeRecipient() common.ExecutionAddress
	GetStateRoot() primitives.Bytes32
	GetReceiptsRoot() primitives.Bytes32
	GetLogsBloom() []byte
	GetBlobGasUsed() math.U64
	GetExcessBlobGas() math.U64
}

// ExecutionPayload represents the execution data of a block.
type ExecutionPayload[WithdrawalT any] interface {
	ExecutionPayloadBody
	GetTransactions() [][]byte
	GetWithdrawals() []WithdrawalT
}

// ExecutionPayloadHeader represents the execution header of a block.
type ExecutionPayloadHeader interface {
	ExecutionPayloadBody
	GetTransactionsRoot() primitives.Root
	GetWithdrawalsRoot() primitives.Root
}

// PayloadAttributer represents payload attributes of a block.
type PayloadAttributer interface {
	// IsNil returns true if the PayloadAttributer is nil.
	IsNil() bool
	// Version returns the version of the PayloadAttributer.
	Version() uint32
	// Validate checks if the PayloadAttributer is valid and returns an error if
	// it is not.
	Validate() error
	// GetSuggestedFeeRecipient returns the suggested fee recipient for the
	// block.
	GetSuggestedFeeRecipient() common.ExecutionAddress
}
