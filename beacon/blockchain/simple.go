package blockchain

import (
	"context"
	"time"

	"cosmossdk.io/core/header"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

func (s *Service) BuildNewBlock(ctx context.Context, beaconBlock header.Info, headHash []byte) (interfaces.ExecutionData, error) {
	payloadIDNew, err := s.notifyForkchoiceUpdate(ctx, uint64(beaconBlock.Height), &notifyForkchoiceUpdateArg{
		headHash: headHash,
	}, true)

	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)

	payload, _, _, err := s.engine.GetPayload(ctx, [8]byte(payloadIDNew[:]), primitives.Slot(beaconBlock.Height))
	return payload, err
}
