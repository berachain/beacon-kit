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

package beacon

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

type ValidatorResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                any  `json:"data"`
}

type GenesisData struct {
	GenesisTime           string         `json:"genesis_time"`
	GenesisValidatorsRoot common.Bytes32 `json:"genesis_validators_root"`
	GenesisForkVersion    string         `json:"genesis_fork_version"`
}

type RootData struct {
	Root common.Root `json:"root"`
}

type ValidatorData[ValidatorT any] struct {
	ValidatorBalanceData
	Status    string     `json:"status"`
	Validator ValidatorT `json:"validator"`
}

type ValidatorBalanceData struct {
	Index   uint64 `json:"index,string"`
	Balance uint64 `json:"balance,string"`
}

//nolint:staticcheck // todo: figure this out.
type CommitteeData struct {
	Index      uint64   `json:"index,string"`
	Slot       uint64   `json:"slot,string"`
	Validators []uint64 `json:"validators,string"`
}

type BlockRewardsData struct {
	ProposerIndex     uint64 `json:"proposer_index,string"`
	Total             uint64 `json:"total,string"`
	Attestations      uint64 `json:"attestations,string"`
	SyncAggregate     uint64 `json:"sync_aggregate,string"`
	ProposerSlashings uint64 `json:"proposer_slashings,string"`
	AttesterSlashings uint64 `json:"attester_slashings,string"`
}
