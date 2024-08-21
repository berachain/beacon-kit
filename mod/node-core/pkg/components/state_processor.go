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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// StateProcessorInput is the input for the state processor for the depinject
// framework.
type StateProcessorInput struct {
	depinject.In
	ChainSpec       common.ChainSpec
	ExecutionEngine *ExecutionEngine
	Signer          crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, *Deposit,
		*Eth1Data, *ExecutionPayload, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT, *Validator,
		Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	KVStoreT BeaconStore[
		KVStoreT, BeaconBlockHeaderT, *Eth1Data, *ExecutionPayloadHeader,
		*Fork, *Validator, Validators, *Withdrawal,
	],
](
	in StateProcessorInput,
) *core.StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, *Context, *Deposit, *Eth1Data, *ExecutionPayload,
	*ExecutionPayloadHeader, *Fork, *ForkData, KVStoreT, *Validator,
	Validators, *Withdrawal, Withdrawals, WithdrawalCredentials,
] {
	return core.NewStateProcessor[
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		*Context,
		*Deposit,
		*Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Fork,
		*ForkData,
		KVStoreT,
		*Validator,
		Validators,
		*Withdrawal,
		Withdrawals,
		WithdrawalCredentials,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
