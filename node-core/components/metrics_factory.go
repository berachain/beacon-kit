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
	"github.com/berachain/beacon-kit/config/config"
	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/observability/metrics/discard"
	"github.com/berachain/beacon-kit/observability/metrics/prometheus"
)

type MetricsFactoryInput struct {
	depinject.In
	Config *config.Config
}

// ProvideMetricsFactory provides a metrics factory based on configuration.
// When Config.Telemetry.Enabled is true, creates real Prometheus metrics.
// When false, creates no-op metrics with zero runtime overhead.
// This setting affects ALL metrics in the beacon-kit system.
func ProvideMetricsFactory(in MetricsFactoryInput) metrics.Factory {
	if in.Config.Telemetry.Enabled {
		return prometheus.NewFactory("beacon_kit")
	}
	return discard.NewFactory()
}
