package types

import (
	"fmt"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// WrapPayload sets the payload data from an `engine.ExecutionPayloadEnvelope`.
func WrapPayload(envelope interfaces.ExecutionData) (*WrappedPayloadEnvelope, error) {
	bz, err := envelope.MarshalSSZ()
	if err != nil {
		return nil, fmt.Errorf("failed to wrap payload: %w", err)
	}

	return &WrappedPayloadEnvelope{
		Data: bz,
	}, nil
}

// AsPayload extracts the payload as an `engine.ExecutionPayloadEnvelope`.
func (wpe *WrappedPayloadEnvelope) UnwrapPayload() interfaces.ExecutionData {
	payload := new(enginev1.ExecutionPayloadCapellaWithValue)
	payload.Payload = new(enginev1.ExecutionPayloadCapella)
	if err := payload.Payload.UnmarshalSSZ(wpe.Data); err != nil {
		return nil
	}

	// todo handle hardforks without needing codechange.
	data, err := blocks.WrappedExecutionPayloadCapella(
		payload.Payload, blocks.PayloadValueToGwei(payload.Value),
	)
	if err != nil {
		return nil
	}
	return data
}
