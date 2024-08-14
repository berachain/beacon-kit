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

const (
	SafeBlockNumbeStr      = "safe"
	FinalizedBlockNumStr   = "finalized"
	LatestBlockNumStr      = "latest"
	PendingBlockNumStr     = "pending"
	EarliestBlockNumberStr = "earliest"
)

const (
	SafeBlockNumber      = BlockNumber(-4)
	FinalizedBlockNumber = BlockNumber(-3)
	LatestBlockNumber    = BlockNumber(-2)
	PendingBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber  = BlockNumber(0)
)

type BlockNumber int64

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "safe", "finalized", "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))

	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case EarliestBlockNumberStr:
		*bn = EarliestBlockNumber
		return nil
	case LatestBlockNumStr:
		*bn = LatestBlockNumber
		return nil
	case PendingBlockNumStr:
		*bn = PendingBlockNumber
		return nil
	case FinalizedBlockNumStr:
		*bn = FinalizedBlockNumber
		return nil
	case SafeBlockNumbeStr:
		*bn = SafeBlockNumber
		return nil
	}

	result := new(math.U64)
	result.UnmarshalJSON([]byte(input))
	if *result > math.U64(uint64(stdmath.MaxInt64)) {
		return errors.New("block number larger than int64")
	}
	//#nosec:G701 // handled by the guard above.
	*bn = BlockNumber(*result.UnwrapPtr())
	return nil
}

// MarshalText implements encoding.TextMarshaler. It marshals:
// - "safe", "finalized", "latest", "earliest" or "pending" as strings
// - other numbers as hex
func (bn BlockNumber) MarshalText() ([]byte, error) {
	return []byte(bn.String()), nil
}

func (bn BlockNumber) String() string {
	switch bn {
	case EarliestBlockNumber:
		return EarliestBlockNumberStr
	case LatestBlockNumber:
		return LatestBlockNumStr
	case PendingBlockNumber:
		return PendingBlockNumStr
	case FinalizedBlockNumber:
		return FinalizedBlockNumStr
	case SafeBlockNumber:
		return SafeBlockNumbeStr
	default:
		if bn < 0 {
			return fmt.Sprintf("<invalid %d>", bn)
		}
		return math.U64(bn).Hex()
	}
}
