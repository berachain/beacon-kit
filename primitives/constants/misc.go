// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package constants

import (
	"github.com/berachain/beacon-kit/primitives/math"
)

// This file contains various constants as defined:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#misc
const (
	// GenesisSlot represents the initial slot in the system.
	GenesisSlot math.Slot = 0
	// GenesisEpoch represents the initial epoch in the system.
	GenesisEpoch math.Epoch = 0
	// FarFutureEpoch represents a far future epoch value.
	FarFutureEpoch = ^uint64(0)
)

// Berachain constants.
const (
	// FirstDepositIndex represents the index of the first deposit in the system, set at genesis.
	FirstDepositIndex uint64 = 0
)

// State list lengths.
const (
	// ValidatorsRegistryLimit is the maximum number of validators that can be registered.
	// https://github.com/ethereum/consensus-specs/blob/dev/presets/mainnet/phase0.yaml#L49
	// 2**40 (= 1,099,511,627,776) validator spots.
	ValidatorsRegistryLimit = 1_099_511_627_776

	// FullExitRequestAmount is the request amount for a full exit request, i.e. when a validator wants to withdraw
	// its entire balance.
	FullExitRequestAmount = 0

	// PendingPartialWithdrawalsLimit is the maximum number of pending partial withdrawals.
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#state-list-lengths
	// 2**27 (= 134,217,728) pending partial withdrawals
	// If the limit is hit, any new partial withdrawal requests will be dropped. This is not likely to happen but
	// theoretically possible.
	PendingPartialWithdrawalsLimit = 134_217_728
)
