//go:build simulated

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

package simulated

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/node-core/components"
)

func FixedComponents(t *testing.T) []any {
	t.Helper()
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
	c = append(c, components.ProvideNodeAPIHandlers,
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

// ProvideElectraGenesisChainSpec provides a chain spec with pectra as the genesis.
func ProvideElectraGenesisChainSpec() (chain.Spec, error) {
	specData := spec.TestnetChainSpecData()
	// Both Deneb1 and Electra happen in genesis.
	specData.GenesisTime = 0
	specData.Deneb1ForkTime = 0
	specData.ElectraForkTime = 0
	chainSpec, err := chain.NewSpec(specData)
	if err != nil {
		return nil, err
	}
	return chainSpec, nil
}

// ProvideSimulationChainSpec provides a default chain-spec equivalent to testnet.
// Bypasses the need for environment variables.
func ProvideSimulationChainSpec() (chain.Spec, error) {
	specData := spec.TestnetChainSpecData()
	specData.GenesisTime = 0
	// Arbitrary number
	specData.Deneb1ForkTime = 30
	chainSpec, err := chain.NewSpec(specData)
	if err != nil {
		return nil, err
	}
	return chainSpec, nil
}

// ProvidePectraForkTestChainSpec provides a chain spec with pectra at timestamp 10
func ProvidePectraForkTestChainSpec() (chain.Spec, error) {
	specData := spec.TestnetChainSpecData()
	specData.GenesisTime = 0
	specData.Deneb1ForkTime = 0
	specData.ElectraForkTime = 10
	chainSpec, err := chain.NewSpec(specData)
	if err != nil {
		return nil, err
	}
	return chainSpec, nil
}
