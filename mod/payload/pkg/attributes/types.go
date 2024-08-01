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

package attributes

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// BeaconState is an interface for accessing the beacon state.
type BeaconState[WithdrawalT any] interface {
	// GetRandaoMixAtIndex returns the randao mix at the given index.
	GetRandaoMixAtIndex(index uint64) (common.Root, error)
}

// StateProcessor is the interface for the state processor.
type StateProcessor[BeaconStateT any, WithdrawalT any] interface {
	// ProcessState processes the state.
	ExpectedWithdrawals(
		BeaconStateT,
	) ([]WithdrawalT, error)
}

// PayloadAttributes is the interface for the payload attributes.
type PayloadAttributes[SelfT any, WithdrawalT any] interface {
	engineprimitives.PayloadAttributer
	// New creates a new payload attributes instance.
	New(
		uint32,
		uint64,
		common.Bytes32,
		gethprimitives.ExecutionAddress,
		[]WithdrawalT,
		common.Root,
	) (SelfT, error)
}
