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
type StateProcessorInput[
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	WithdrawalT any,
	WithdrawalsT Withdrawals[WithdrawalT],
] struct {
	depinject.In
	ChainSpec       common.ChainSpec
	ExecutionEngine ExecutionEngineT
	Signer          crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[
		BeaconBlockT, AttestationDataT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, AttestationDataT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	ContextT Context[ContextT],
	DepositT Deposit[DepositT, ForkDataT, WithdrawalCredentialsT],
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT Fork[ForkT],
	ForkDataT ForkData[ForkDataT],
	KVStoreT any,
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	SlashingInfoT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT Validators[ValidatorT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
	WithdrawalCredentialsT ~[32]byte,
](
	in StateProcessorInput[
		ExecutionEngineT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalT, WithdrawalsT,
	],
) *core.StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
	ContextT, DepositT, Eth1DataT, ExecutionEngineT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
	ValidatorsT, WithdrawalT, WithdrawalsT, WithdrawalCredentialsT,
] {
	return core.NewStateProcessor[
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		ContextT,
		DepositT,
		Eth1DataT,
		ExecutionEngineT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		ForkT,
		ForkDataT,
		KVStoreT,
		ValidatorT,
		ValidatorsT,
		WithdrawalT,
		WithdrawalsT,
		WithdrawalCredentialsT,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
