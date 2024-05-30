// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package version

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/log"
)

// defaultReportingInterval is the default interval at which the version is
// reported.
const defaultReportingInterval = 5 * time.Minute

// ReportingService is a service that periodically logs the running chain
// version.
type ReportingService struct {
	// logger is used to log information about the running chain version.
	logger log.Logger[any]
	// version represents the current version of the running chain.
	version string
	// ticker is used to trigger periodic logging at specified intervals.
	ticker *time.Ticker
	// metrics contains the metrics for the version service.
	metrics *versionMetrics
}

// NewReportingService creates a new VersionReporterService.
func NewReportingService(
	logger log.Logger[any],
	telemetrySink TelemetrySink,
	version string,
) *ReportingService {
	return &ReportingService{
		logger:  logger,
		version: version,
		ticker: time.NewTicker(
			defaultReportingInterval,
		),
		metrics: newVersionMetrics(logger, telemetrySink),
	}
}

// Name returns the name of the service.
func (*ReportingService) Name() string {
	return "ReportingService"
}

// Start begins the periodic logging of the chain version.
func (v *ReportingService) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				v.ticker.Stop()
				return
			case <-v.ticker.C:
				v.metrics.reportVersion(v.version)
			}
		}
	}()
	return nil
}

// Status returns nil if the service is healthy.
func (*ReportingService) Status() error {
	return nil
}

// WaitForHealthy waits for all registered services to be healthy.
func (*ReportingService) WaitForHealthy(context.Context) {}
