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
	"github.com/berachain/beacon-kit/mod/cli/pkg/flags"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
)

// DefaultClientComponents returns the default components for
// the client.
func DefaultClientComponents() []any {
	return []any{
		ProvideClientContext,
		ProvideKeyring,
	}
}

// TODO: Decouple cli flags package from here, and / or
// remove this function.
func DefaultNodeComponents() []any {
	return []any{
		components.ProvideABCIMiddleware,
		components.ProvideAvailabilityPruner,
		components.GetAVSProvider[*components.BeaconBlockBody](
			sdkflags.FlagHome,
		),
		components.GetBlsSignerProvider(sdkflags.FlagHome),
		components.ProvideBlockFeed,
		components.ProvideBlobProcessor[*components.BeaconBlockBody],
		components.GetBlobProofVerifierProvider(
			flags.KZGImplementation,
		),
		components.ProvideChainService,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideDBManager,
		components.ProvideDepositPruner,
		components.ProvideDepositService,
		components.GetDepositStoreProvider[*components.Deposit](
			sdkflags.FlagHome,
		),
		components.ProvideBeaconDepositContract[
			*components.Deposit, *components.ExecutionPayload,
			*components.Withdrawal, types.WithdrawalCredentials,
		],
		components.ProvideEngineClient[*components.ExecutionPayload],
		components.ProvideExecutionEngine[*components.ExecutionPayload],
		components.ProvideFinalizeBlockMiddleware,
		components.GetJWTSecretProvider(flags.JWTSecretPath),
		components.ProvideLocalBuilder,
		components.ProvideServiceRegistry,
		components.ProvideStateProcessor,
		components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		components.GetTrustedSetupProvider(flags.KZGTrustedSetupPath),
		components.ProvideValidatorMiddleware,
		components.ProvideValidatorService,
	}
}
