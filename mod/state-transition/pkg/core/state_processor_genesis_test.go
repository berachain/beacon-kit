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

package core_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
)

type (
	TestBeaconStateMarshallableT = types.BeaconState[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.BeaconBlockHeader,
		types.Eth1Data,
		types.ExecutionPayloadHeader,
		types.Fork,
		types.Validator,
	]

	TestKVStoreT = beacondb.KVStore[
		*types.BeaconBlockHeader,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		types.Validators,
	]

	TestBeaconStateT = statedb.StateDB[
		*types.BeaconBlockHeader,
		*TestBeaconStateMarshallableT,
		*types.Eth1Data,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		types.WithdrawalCredentials,
	]
)

func TestInitialize(t *testing.T) {
	cs := spec.TestnetChainSpec()
	// in.ExecutionEngine,
	mocksSigner := &mocks.BLSSigner{}

	_ = core.NewStateProcessor[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		*TestBeaconStateT,

		*transition.Context,
		*types.Deposit,
		*types.Eth1Data,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*TestKVStoreT,
		*types.Validator,
		types.Validators,
		*engineprimitives.Withdrawal,
		engineprimitives.Withdrawals,
		types.WithdrawalCredentials,
	](
		cs,
		nil,
		mocksSigner,
	)

	// vals, err := sp.InitializePreminedBeaconStateFromEth1()
	// require.NoError(t, err)

	// This is not the right assert, just want to get to compile
	// require.Empty(t, vals)

}
