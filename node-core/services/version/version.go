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

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/client/ethclient"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/constraints"
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
	// we print to console always at the beginning
	rs.printToConsole(engineprimitives.ClientVersionV1{
		Version: "unknown",
		Name:    "unknown"},
	)

	connectedTicker := time.NewTicker(time.Second)
	go func() {
		// wait until the client is connected
		connected := false
		for !connected {
			select {
			case <-connectedTicker.C:
				connected = rs.client.IsConnected()
			case <-ctx.Done():
				connectedTicker.Stop()
				return
			}
		}
		connectedTicker.Stop()

		rs.logger.Info("Connected to execution client")

		// log telemetry immediately after we are connected
		ethVersion, err := rs.GetEthVersion(ctx)
		if err != nil {
			rs.logger.Warn("Failed to get eth version", "err", err)
		}
		rs.logTelemetry(ethVersion)

		// then we start reporting at the reportingInterval interval
		ticker := time.NewTicker(rs.reportingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// since the eth client can be updated separately for beacon
				// node
				// we need to fetch the version every time
				ethVersion, err = rs.GetEthVersion(ctx)
				if err != nil {
					rs.logger.Warn("Failed to get eth version", "err", err)
				}

				// print to console and log telemetry
				rs.printToConsole(ethVersion)
				rs.logTelemetry(ethVersion)

				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (rs *ReportingService[_, _]) printToConsole(
	ethClient engineprimitives.ClientVersionV1) {
	rs.logger.Info(fmt.Sprintf(`


	+==========================================================================+
	+ ⭐️ Star BeaconKit on GitHub @ https://github.com/berachain/beacon-kit    +
	+ 🧩 Your node is running version: %-40s+
	+ ♦ Eth client: %-59s+
	+ 💾 Your system: %-57s+
	+ 🦺 Please report issues @ https://github.com/berachain/beacon-kit/issues +
	+==========================================================================+


`,
		rs.version,
		fmt.Sprintf("%s (version: %s)", ethClient.Name, ethClient.Version),
		runtime.GOOS+"/"+runtime.GOARCH,
	))
}

func (rs *ReportingService[_, _]) GetEthVersion(
	ctx context.Context) (engineprimitives.ClientVersionV1, error) {
	ethVersion := engineprimitives.ClientVersionV1{
		Version: "unknown",
		Name:    "unknown",
	}

	if rs.client.HasCapability(ethclient.GetClientVersionV1) {
		// Get the client version from the execution layer.
		info, err := rs.client.GetClientVersionV1(ctx)
		if err != nil {
			return ethVersion, fmt.Errorf(
				"failed to get client version: %w",
				err,
			)
		}

		// the spec says we should have at least one client version
		if len(info) == 0 {
			return ethVersion, errors.New("no client version returned")
		}

		ethVersion.Version = info[0].Version
		ethVersion.Name = info[0].Name
	} else {
		rs.logger.Warn("Client does not have capability to get client version")
	}

	return ethVersion, nil
}

func (rs *ReportingService[_, _]) logTelemetry(
	ethVersion engineprimitives.ClientVersionV1) {
	systemInfo := runtime.GOOS + "/" + runtime.GOARCH

	// TODO: Delete this counter as it should be included in the new
	// beacon_kit.runtime.version metric.
	rs.sink.IncrementCounter(
		"beacon_kit.runtime.version.reported",
		"version", rs.version, "system", systemInfo,
	)

	rs.logger.Info("Reporting version", "version", rs.version,
		"system", systemInfo,
		"eth_version", ethVersion.Version,
		"eth_name", ethVersion.Name)

	// Report the version to the telemetry sink and include labels
	// for beacon node version and eth name and version
	var args = [8]string{
		"version", rs.version,
		"system", systemInfo,
		"eth_version", ethVersion.Version,
		"eth_name", ethVersion.Name,
	}
	rs.sink.SetGauge("beacon_kit.runtime.version", 1, args[:]...)
}
