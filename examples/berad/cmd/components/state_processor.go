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

package components

import (
	"cosmossdk.io/depinject"
	spec "github.com/berachain/beacon-kit/examples/berad/pkg/chain-spec"
	core "github.com/berachain/beacon-kit/examples/berad/pkg/state-transition"
	"github.com/berachain/beacon-kit/examples/berad/pkg/state-transition/state"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	beaconcomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
)

// StateProcessorInput is the input for the state processor for the depinject
// framework.
type StateProcessorInput[
	ExecutionPayloadT beaconcomponents.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT beaconcomponents.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	WithdrawalT beaconcomponents.Withdrawal[WithdrawalT],
	WithdrawalsT beaconcomponents.Withdrawals[WithdrawalT],
] struct {
	depinject.In
	ChainSpec       spec.BeraChainSpec
	ExecutionEngine *engine.Engine[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
		beaconcomponents.PayloadID,
		WithdrawalsT,
	]
	Signer crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor[
	BeaconBlockT beaconcomponents.BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT beaconcomponents.BeaconBlockBody[
		BeaconBlockBodyT, *beaconcomponents.AttestationData, DepositT,
		*beaconcomponents.Eth1Data, ExecutionPayloadT, *beaconcomponents.SlashingInfo,
	],
	BeaconBlockHeaderT beaconcomponents.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*beaconcomponents.Eth1Data, ExecutionPayloadHeaderT,
		*beaconcomponents.Fork, KVStoreT, *beaconcomponents.Validator,
		beaconcomponents.Validators, WithdrawalT, WithdrawalsT,
	],
	BeaconStateMarshallableT any,
	DepositT beaconcomponents.Deposit[
		DepositT, *beaconcomponents.ForkData, beaconcomponents.WithdrawalCredentials,
	],
	ExecutionPayloadT beaconcomponents.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT beaconcomponents.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	KVStoreT state.KVStore[
		KVStoreT, BeaconBlockHeaderT, ExecutionPayloadHeaderT,
		*beaconcomponents.Fork, *beaconcomponents.Validator, beaconcomponents.Validators,
	],
	WithdrawalsT beaconcomponents.Withdrawals[WithdrawalT],
	WithdrawalT beaconcomponents.Withdrawal[WithdrawalT],
](
	in StateProcessorInput[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT, WithdrawalsT,
	],
) *core.StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, *beaconcomponents.Context, DepositT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, *beaconcomponents.Fork, *beaconcomponents.ForkData, KVStoreT, *beaconcomponents.Validator,
	beaconcomponents.Validators, WithdrawalT, WithdrawalsT, beaconcomponents.WithdrawalCredentials,
] {
	return core.NewStateProcessor[
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		*beaconcomponents.Context,
		DepositT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		*beaconcomponents.Fork,
		*beaconcomponents.ForkData,
		KVStoreT,
		*beaconcomponents.Validator,
		beaconcomponents.Validators,
		WithdrawalT,
		WithdrawalsT,
		beaconcomponents.WithdrawalCredentials,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
