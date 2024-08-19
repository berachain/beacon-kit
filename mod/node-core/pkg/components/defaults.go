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
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
)

type LoggerT = *phuslu.Logger

func DefaultComponents() []any {
	components := []any{
		ProvideABCIMiddleware[LoggerT],
		ProvideAttributesFactory[LoggerT],
		ProvideAvailabilityPruner[LoggerT],
		ProvideAvailibilityStore[LoggerT],
		ProvideBeaconDepositContract,
		ProvideBlockPruner[LoggerT],
		ProvideBlockStore[LoggerT],
		ProvideBlockStoreService[LoggerT],
		ProvideBlsSigner,
		ProvideBlobProcessor[LoggerT],
		ProvideBlobProofVerifier,
		ProvideBlobVerifier,
		ProvideChainService[LoggerT],
		ProvideChainSpec,
		ProvideConfig,
		ProvideConsensusEngine,
		ProvideDAService[LoggerT],
		ProvideDBManager[LoggerT],
		ProvideDepositPruner[LoggerT],
		ProvideDepositService[LoggerT],
		ProvideDepositStore,
		ProvideEngineClient[LoggerT],
		ProvideExecutionEngine[LoggerT],
		ProvideJWTSecret,
		ProvideLocalBuilder[LoggerT],
		ProvideReportingService[LoggerT],
		ProvideServiceRegistry[LoggerT],
		ProvideSidecarFactory,
		ProvideStateProcessor,
		ProvideKVStore,
		ProvideSSZBackend,
		ProvideStorageBackend,
		ProvideTelemetrySink,
		ProvideTrustedSetup,
		ProvideValidatorService[LoggerT],
		// TODO Hacks
		ProvideKVStoreService,
		ProvideKVStoreKey,
	}
	components = append(components,
		ProvideNodeAPIServer[LoggerT],
		ProvideNodeAPIEngine,
		ProvideNodeAPIBackend,
	)
	components = append(components, DefaultNodeAPIHandlers()...)
	components = append(components, DefaultBrokerProviders()...)
	return components
}
