package deposit

import (
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// depositMetrics is a struct that contains metrics for the deposit service.
type depositMetrics struct {
	// sink is the telemetry sink.
	sink TelemetrySink
}

// newDepositMetrics creates a new instance of the depositMetrics struct.
func newDepositMetrics(sink TelemetrySink) *depositMetrics {
	return &depositMetrics{
		sink: sink,
	}
}

// markFailedToGetBlockLogs increments the counter for failed to get block logs.
func (m *depositMetrics) markFailedToGetBlockLogs(blockNum math.U64) {
	m.sink.IncrementCounter(
		"beacon_kit.execution.deposit.failed_to_get_block_logs",
		"block_num",
		strconv.FormatUint(uint64(blockNum), 10),
	)
}
