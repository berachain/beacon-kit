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

package mock_consensus_test

import (
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components"
	"testing"
)

func DefaultComponents(t *testing.T) []any {
	c := []any{
		//components.ProvideAttributesFactory[*phuslu.Logger],
		//components.ProvideAvailabilityStore[*phuslu.Logger],
		//components.ProvideDepositContract,
		//components.ProvideBlockStore[*phuslu.Logger],
		components.ProvideBlsSigner,
		//components.ProvideBlobProcessor[*phuslu.Logger],
		//components.ProvideBlobProofVerifier,
		//components.ProvideChainService[*phuslu.Logger],
		//components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		//components.ProvideServerConfig,
		// Using in-memory Deposit Store
		components.ProvideDepositStoreInMemory[*phuslu.Logger],
		components.ProvideEngineClient[*phuslu.Logger],
		components.ProvideExecutionEngine[*phuslu.Logger],
		//components.ProvideJWTSecret,
		//components.ProvideLocalBuilder[*phuslu.Logger],
		//components.ProvideReportingService[*phuslu.Logger],
		//components.ProvideCometBFTService[*phuslu.Logger],
		//components.ProvideServiceRegistry[*phuslu.Logger],
		//components.ProvideSidecarFactory,
		components.ProvideStateProcessor[*phuslu.Logger],
		//components.ProvideKVStore,
		//components.ProvideStorageBackend,
		components.ProvideTelemetrySink,
		//components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		//components.ProvideValidatorService[*phuslu.Logger],
		//components.ProvideShutDownService[*phuslu.Logger],
		//clicomponents.ProvideClientContext,
	}
	c = append(c,
		//components.ProvideNodeAPIServer[*phuslu.Logger, echo.Context],
		//components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend,
	)
	//
	c = append(c) //	components.ProvideNodeAPIHandlers[echo.Context],
	//	components.ProvideNodeAPIBeaconHandler[echo.Context],
	//	components.ProvideNodeAPIBuilderHandler[echo.Context],
	//	components.ProvideNodeAPIConfigHandler[echo.Context],
	//	components.ProvideNodeAPIDebugHandler[echo.Context],
	//	components.ProvideNodeAPIEventsHandler[echo.Context],
	//	components.ProvideNodeAPINodeHandler[echo.Context],
	//	components.ProvideNodeAPIProofHandler[echo.Context],

	return c
}
