package engineprimitives_test

import (
	"encoding/binary"
	"encoding/json"
	"github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives/mocks"
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

type MockExecutionPayloadT struct {
	Value string
}

func (m MockExecutionPayloadT) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (m MockExecutionPayloadT) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Value)
}

func TestExecutionPayloadEnvelope(t *testing.T) {

	BlobsBundle := &mocks.BlobsBundle{}

	// Convert the int to a byte slice
	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(100))

	// Use math.NewU256L
	blockValue, err := math.NewU256L(valueBytes)
	if err != nil {
		t.Fatalf("Failed to convert int to U256L: %v", err)
	}

	envelope := &engineprimitives.ExecutionPayloadEnvelope[MockExecutionPayloadT, *mocks.BlobsBundle]{
		ExecutionPayload: MockExecutionPayloadT{Value: "test"},
		BlockValue:       blockValue,
		BlobsBundle:      BlobsBundle,
		Override:         true,
	}

	payload := envelope.GetExecutionPayload()
	require.Equal(t, envelope.ExecutionPayload, payload)

	value := envelope.GetValue()
	require.Equal(t, envelope.BlockValue, value)

	bundle := envelope.GetBlobsBundle()
	require.Equal(t, envelope.BlobsBundle, bundle)

	override := envelope.ShouldOverrideBuilder()
	require.Equal(t, envelope.Override, override)
}
