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

package version

import (
	"fmt"
	"runtime"

	"github.com/berachain/beacon-kit/mod/interfaces/pkg/telemetry"
	"github.com/berachain/beacon-kit/mod/log"
)

// versionMetrics holds metrics related to the version reporting.
type versionMetrics struct {
	// system is the current system the node is running on.
	system string
	// logger is the logger used to log information about the version.
	logger log.Logger[any]
	// sink is the telemetry sink used to report metrics.
	sink telemetry.Sink
}

// newVersionMetrics creates a new instance of versionMetrics.
func newVersionMetrics(
	logger log.Logger[any],
	sink telemetry.Sink,
) *versionMetrics {
	return &versionMetrics{
		system: runtime.GOOS + "/" + runtime.GOARCH,
		logger: logger,
		sink:   sink,
	}
}

// reportVersion increments the versionReported counter.
func (vm *versionMetrics) reportVersion(version string) {
	vm.logger.Info(fmt.Sprintf(`


	+==========================================================================+
	+ ⭐️ Star BeaconKit on GitHub @ https://github.com/berachain/beacon-kit    +
	+ 🧩 Your node is running version: %-40s+
	+ 💾 Your system: %-57s+
	+ 🦺 Please report issues @ https://github.com/berachain/beacon-kit/issues +
	+==========================================================================+


`,
		version,
		vm.system,
	))
	vm.sink.IncrementCounter(
		"beacon_kit.runtime.version.reported",
		"version", version, "system", vm.system,
	)
}
