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
func newVersionMetrics(logger log.Logger[any], sink TelemetrySink) *versionMetrics {
	return &versionMetrics{
		sink: sink,
	}
}

// reportVersion increments the versionReported counter.
func (vm *versionMetrics) reportVersion(version string) {
	vm.logger.Info("this node is running", "version", version)
	vm.sink.IncrementCounter(
		"beacon_kit.runtime.version.reported", "version", version)
}
