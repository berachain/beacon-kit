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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/blobreactor"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/state-transition/core/state"
)

// Simple metrics providers (no additional dependencies)

func ProvideBlobFactoryMetrics(factory metrics.Factory) *dablob.FactoryMetrics {
	return dablob.NewFactoryMetrics(factory)
}

func ProvideBlobProcessorMetrics(factory metrics.Factory) *dablob.ProcessorMetrics {
	return dablob.NewProcessorMetrics(factory)
}

func ProvideBlobVerifierMetrics(factory metrics.Factory) *dablob.VerifierMetrics {
	return dablob.NewVerifierMetrics(factory)
}

func ProvideBlobFetcherMetrics(factory metrics.Factory) *blockchain.BlobFetcherMetrics {
	return blockchain.NewBlobFetcherMetrics(factory)
}

func ProvideBlobReactorMetrics(factory metrics.Factory) *blobreactor.Metrics {
	return blobreactor.NewMetrics(factory)
}

func ProvideBlockchainMetrics(factory metrics.Factory) *blockchain.Metrics {
	return blockchain.NewMetrics(factory)
}

func ProvideValidatorMetrics(factory metrics.Factory) *validator.Metrics {
	return validator.NewMetrics(factory)
}

func ProvideStateDBMetrics(factory metrics.Factory) *state.Metrics {
	return state.NewMetrics(factory)
}

func ProvideStateProcessorMetrics(factory metrics.Factory) *core.Metrics {
	return core.NewMetrics(factory)
}

type ExecutionClientMetricsInput struct {
	depinject.In
	Factory metrics.Factory
	Logger  *phuslu.Logger
}

func ProvideExecutionClientMetrics(in ExecutionClientMetricsInput) *client.Metrics {
	return client.NewMetrics(in.Factory, in.Logger.With("service", "execution-client"))
}

type ExecutionEngineMetricsInput struct {
	depinject.In
	Factory metrics.Factory
	Logger  *phuslu.Logger
}

func ProvideExecutionEngineMetrics(in ExecutionEngineMetricsInput) *engine.Metrics {
	return engine.NewMetrics(in.Factory, in.Logger.With("service", "execution-engine"))
}

// AllMetricsProviders returns all metrics provider functions for depinject.
// This helper groups all metrics providers together to reduce boilerplate
// in component lists (defaults.go, components.go, etc.).
func AllMetricsProviders() []any {
	return []any{
		ProvidePrometheusRegisterer, // Must be first - creates registerer used by factory
		ProvideMetricsFactory,       // Must be second - creates factory used by all others
		ProvideStateDBMetrics,
		ProvideValidatorMetrics,
		ProvideExecutionClientMetrics,
		ProvideExecutionEngineMetrics,
		ProvideBlobFactoryMetrics,
		ProvideBlobProcessorMetrics,
		ProvideBlobVerifierMetrics,
		ProvideBlockchainMetrics,
		ProvideBlobFetcherMetrics,
		ProvideBlobReactorMetrics,
		ProvideStateProcessorMetrics,
		ProvideCometBFTMetrics,
	}
}
