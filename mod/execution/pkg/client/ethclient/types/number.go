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
	"errors"
	"fmt"
	stdmath "math"
	"strings"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BlockNumber represents a block number.
const (
	// SafeBlockNumberStr is the string representation of SafeBlockNumber.
	SafeBlockNumberStr = "safe"
	// FinalizedBlockNumberStr is the string representation of FinalizedBlockNumber.
	FinalizedBlockNumberStr = "finalized"
	// LatestBlockNumberStr is the string representation of LatestBlockNumber.
	LatestBlockNumberStr = "latest"
	// PendingBlockNumberStr is the string representation of PendingBlockNumber.
	PendingBlockNumberStr = "pending"
	// EarliestBlockNumberStr is the string representation of EarliestBlockNumber.
	EarliestBlockNumberStr = "earliest"

	// SafeBlockNumber represents a safe block to use for certain operations.
	SafeBlockNumber BlockNumber = -4
	// FinalizedBlockNumber represents the finalized block.
	FinalizedBlockNumber BlockNumber = -3
	// LatestBlockNumber represents the latest block in the chain.
	LatestBlockNumber BlockNumber = -2
	// PendingBlockNumber represents the pending block.
	PendingBlockNumber BlockNumber = -1
	// EarliestBlockNumber represents the earliest block in the chain.
	EarliestBlockNumber BlockNumber = 0
)

// BlockNumber represents a block number.
type BlockNumber int64

// UnmarshalJSON parses the given JSON fragment into a BlockNumber
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))

	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case EarliestBlockNumberStr:
		*bn = EarliestBlockNumber
		return nil
	case LatestBlockNumberStr:
		*bn = LatestBlockNumber
		return nil
	case PendingBlockNumberStr:
		*bn = PendingBlockNumber
		return nil
	case FinalizedBlockNumberStr:
		*bn = FinalizedBlockNumber
		return nil
	case SafeBlockNumberStr:
		*bn = SafeBlockNumber
		return nil
	}

	result := new(math.U64)
	if err := result.UnmarshalJSON([]byte(input)); err != nil {
		return err
	}
	if *result > math.U64(uint64(stdmath.MaxInt64)) {
		return errors.New("block number larger than int64")
	}
	//#nosec:G701 // handled by the guard above.
	*bn = BlockNumber(*result.UnwrapPtr())
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (bn BlockNumber) MarshalText() ([]byte, error) {
	return []byte(bn.String()), nil
}

// String returns the string representation of the BlockNumber.
func (bn BlockNumber) String() string {
	switch bn {
	case EarliestBlockNumber:
		return EarliestBlockNumberStr
	case LatestBlockNumber:
		return LatestBlockNumberStr
	case PendingBlockNumber:
		return PendingBlockNumberStr
	case FinalizedBlockNumber:
		return FinalizedBlockNumberStr
	case SafeBlockNumber:
		return SafeBlockNumberStr
	default:
		if bn < 0 {
			return fmt.Sprintf("<invalid %d>", bn)
		}
		return math.U64(bn).Hex()
	}
}
