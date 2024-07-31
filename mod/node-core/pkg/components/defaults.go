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

func DefaultComponentsWithStandardTypes() []any {
	return []any{
		ProvideABCIMiddleware,
		ProvideAttributesFactory,
		ProvideAvailabilityPruner,
		ProvideAvailibilityStore,
		ProvideBeaconDepositContract,
		ProvideBlockFeed,
		ProvideBlockPruner,
		ProvideBlockStore,
		ProvideBlockStoreService,
		ProvideBlsSigner,
		ProvideBlobFeed,
		ProvideBlobProcessor,
		ProvideBlobProofVerifier,
		ProvideBlobVerifier,
		ProvideBrokerRegistry,
		ProvideChainService,
		ProvideChainSpec,
		ProvideConfig,
		ProvideConsensusEngine,
		ProvideDAService,
		ProvideDBManager,
		ProvideDepositPruner,
		ProvideDepositService,
		ProvideDepositStore,
		ProvideDispatcher,
		ProvideEngineClient,
		ProvideExecutionEngine,
		ProvideGenesisBroker,
		ProvideJWTSecret,
		ProvideLocalBuilder,
		ProvideNodeAPIBackend,
		ProvideNodeAPIBeaconHandler,
		ProvideNodeAPIBuilderHandler,
		ProvideNodeAPIConfigHandler,
		ProvideNodeAPIDebugHandler,
		ProvideNodeAPIEngine,
		ProvideNodeAPIEventsHandler,
		ProvideNodeAPIHandlers,
		ProvideNodeAPINodeHandler,
		ProvideNodeAPIProofHandler,
		ProvideNodeAPIServer,
		ProvideServiceRegistry,
		ProvideSidecarFactory,
		ProvideSlotBroker,
		ProvideStateProcessor,
		ProvideStatusBroker,
		ProvideStorageBackend,
		ProvideTelemetrySink,
		ProvideTrustedSetup,
		ProvideValidatorService,
		ProvideValidatorUpdateBroker,
	}
}
