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
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// defaultReportingInterval is the default interval at which the version is
// reported.
const defaultReportingInterval = 5 * time.Minute

// ReportingService is a service that periodically logs the running chain
// version.
type ReportingService[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	PayloadAttributesT client.PayloadAttributes,
] struct {
	// logger is used to log information about the running chain version.
	logger log.Logger
	// version represents the current version of the running chain.
	version string
	// reportingInterval is the interval at which the version is reported.
	reportingInterval time.Duration
	// sink is the telemetry sink used to report metrics.
	sink TelemetrySink
	// client to query the execution layer
	client *client.EngineClient[ExecutionPayloadT, PayloadAttributesT]
}

// NewReportingService creates a new VersionReporterService.
func NewReportingService[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	PayloadAttributesT client.PayloadAttributes,
](
	logger log.Logger,
	telemetrySink TelemetrySink,
	version string,
	engineClient *client.EngineClient[ExecutionPayloadT, PayloadAttributesT],
) *ReportingService[
	ExecutionPayloadT, PayloadAttributesT,
] {
	return &ReportingService[
		ExecutionPayloadT, PayloadAttributesT,
	]{
		logger:            logger,
		version:           version,
		reportingInterval: defaultReportingInterval,
		sink:              telemetrySink,
		client:            engineClient,
	}
}

// Name returns the name of the service.
func (*ReportingService[_, _]) Name() string {
	return "reporting"
}

// Start begins the periodic logging of the chain version.
func (rs *ReportingService[_, _]) Start(ctx context.Context) error {
	ticker := time.NewTicker(rs.reportingInterval)
	rs.handleReport(ctx)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rs.handleReport(ctx)
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (rs *ReportingService[_, _]) handleReport(ctx context.Context) {
	systemInfo := runtime.GOOS + "/" + runtime.GOARCH

	rs.logger.Info(fmt.Sprintf(`


	+==========================================================================+
	+ ⭐️ Star BeaconKit on GitHub @ https://github.com/berachain/beacon-kit    +
	+ 🧩 Your node is running version: %-40s+
	+ 💾 Your system: %-57s+
	+ 🦺 Please report issues @ https://github.com/berachain/beacon-kit/issues +
	+==========================================================================+


`,
		rs.version,
		systemInfo,
	))

	// TODO: Delete this counter as it should be included in the new beacon_kit.runtime.version metric.
	rs.sink.IncrementCounter(
		"beacon_kit.runtime.version.reported",
		"version", rs.version, "system", systemInfo,
	)

	// Get the client version from the execution layer.
	info, err := rs.client.GetClientVersionV1(ctx)
	if err != nil {
		rs.logger.Error("Failed to get client version", "err", err)
		return
	}
	rs.logger.Info("GetClientVersionV1", "info", info)

	// the spec says we should have at least one client version
	if len(info) == 0 {
		rs.logger.Warn("No client version returned")
		return
	}

	// Report the version to the telemetry sink and include labels for beacon node version and eth name and version
	var args = [8]string{
		"version", rs.version,
		"system", systemInfo,
		"eth_version", info[0].Version,
		"eth_name", info[0].Name,
	}
	rs.sink.SetGauge("beacon_kit.runtime.version", 1, args[:]...)
}
