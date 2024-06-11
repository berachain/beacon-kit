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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// StateProcessor is the type alias for the state processor inteface.
type StateProcessor = blockchain.StateProcessor[
	*types.BeaconBlock,
	BeaconState,
	*datypes.BlobSidecars,
	*transition.Context,
	*types.Deposit,
]

// StateProcessorInput is the input for the state processor for the depinject
// framework.
type StateProcessorInput struct {
	depinject.In
	ChainSpec       primitives.ChainSpec
	ExecutionEngine *execution.Engine[*types.ExecutionPayload]
	Signer          crypto.BLSSigner
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor(
	in StateProcessorInput,
) StateProcessor {
	return core.NewStateProcessor[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
		*types.Eth1Data,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*types.Validator,
		*engineprimitives.Withdrawal,
		types.WithdrawalCredentials,
	](
		in.ChainSpec,
		in.ExecutionEngine,
		in.Signer,
	)
}
