// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package main

import (
	"github.com/berachain/beacon-kit/node-core/components"
)

func DefaultComponents() []any {
	c := []any{
		components.ProvideAttributesFactory,
		components.ProvideAvailabilityStore,
		components.ProvideDepositContract,
		components.ProvideBlockStore,
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor,
		components.ProvideBlobProofVerifier,
		components.ProvideChainService,
		components.ProvideNode,
		components.ProvideConfig,
		components.ProvideServerConfig,
		components.ProvideDepositStore,
		components.ProvideEngineClient,
		components.ProvideExecutionEngine,
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder,
		components.ProvideReportingService,
		components.ProvideCometBFTService,
		components.ProvideServiceRegistry,
		components.ProvideSidecarFactory,
		components.ProvideStateProcessor,
		components.ProvideKVStore,
		components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService,
		components.ProvideShutDownService,
	}
	c = append(c,
		components.ProvideNodeAPIServer,
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend,
	)

	c = append(c,
		components.ProvideNodeAPIHandlers,
		components.ProvideNodeAPIBeaconHandler,
		components.ProvideNodeAPIBuilderHandler,
		components.ProvideNodeAPIConfigHandler,
		components.ProvideNodeAPIDebugHandler,
		components.ProvideNodeAPIEventsHandler,
		components.ProvideNodeAPINodeHandler,
		components.ProvideNodeAPIProofHandler,
	)

	return c
}
