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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

//go:generate go run github.com/karalabe/ssz/cmd/sszgen -type BeaconState -out gen_state_ssz.go .
type BeaconState struct {
	// Versioning
	//

	GenesisValidatorsRoot common.Root `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  math.Slot   `json:"slot"`
	Fork                  *Fork       `json:"fork"`

	// History
	LatestBlockHeader *BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        []common.Root      `json:"blockRoots"        ssz-size:"8192"`
	StateRoots        []common.Root      `json:"stateRoots"        ssz-size:"8192"`

	// Eth1
	Eth1Data                     *Eth1Data               `json:"eth1Data"`
	Eth1DepositIndex             uint64                  `json:"eth1DepositIndex"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader `json:"latestExecutionPayloadHeader"`

	// Registry
	Validators []*Validator `ssz-max:"1099511627776"`
	Balances   []uint64     `ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes []common.Bytes32 `json:"randaoMixes" ssz-size:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     [8192]uint64 `json:"slashings"     ssz-size:"8192"`
	TotalSlashing math.Gwei    `json:"totalSlashing"`
}
