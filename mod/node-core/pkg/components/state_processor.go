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
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// StateProcessorInput is the input for the state processor for the depinject
// framework.
type StateProcessorInput[
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	WithdrawalsT Withdrawals,
] struct {
	depinject.In
	ChainSpec       common.ChainSpec
	ExecutionEngine ExecutionEngineT
	Signer          crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor[

	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	WithdrawalsT Withdrawals,
](
	in StateProcessorInput[
		ExecutionEngineT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalsT,
	],
) *StateProcessor {
	return core.NewStateProcessor[
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*Context,
		*Deposit,
		*Eth1Data,
		*ExecutionEngine,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Fork,
		*ForkData,
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		engineprimitives.Withdrawals,
		WithdrawalCredentials,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
