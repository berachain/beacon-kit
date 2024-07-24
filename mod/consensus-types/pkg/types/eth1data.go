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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Eth1Data struct {
	// DepositRoot is the root of the deposit tree.
	DepositRoot common.Root `json:"depositRoot"  ssz-size:"32"`
	// DepositCount is the number of deposits in the deposit tree.
	DepositCount uint64 `json:"depositCount"`
	// BlockHash is the hash of the block corresponding to the Eth1Data.
	BlockHash gethprimitives.ExecutionHash `json:"blockHash"    ssz-size:"32"`
}

// New creates a new Eth1Data.
func (e *Eth1Data) New(
	depositRoot common.Root,
	depositCount math.U64,
	blockHash gethprimitives.ExecutionHash,
) *Eth1Data {
	e = &Eth1Data{
		DepositRoot:  depositRoot,
		DepositCount: uint64(depositCount),
		BlockHash:    blockHash,
	}
	return e
}

// GetDepositCount returns the deposit count.
func (e *Eth1Data) GetDepositCount() math.U64 {
	return math.U64(e.DepositCount)
}
