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

import "github.com/berachain/beacon-kit/mod/log"

// versionMetrics holds metrics related to the version reporting.
type versionMetrics struct {
	// logger is the logger used to log information about the version.
	logger log.Logger[any]
	// sink is the telemetry sink used to report metrics.
	sink TelemetrySink
}

// newVersionMetrics creates a new instance of versionMetrics.
func newVersionMetrics(
	logger log.Logger[any],
	sink TelemetrySink,
) *versionMetrics {
	return &versionMetrics{
		logger: logger,
		sink:   sink,
	}
}

// reportVersion increments the versionReported counter.
func (vm *versionMetrics) reportVersion(version string) {
	vm.logger.Info("this node is running", "version", version)
	vm.sink.IncrementCounter(
		"beacon_kit.runtime.version.reported", "version", version)
}
