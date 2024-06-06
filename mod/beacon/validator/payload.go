package validator

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
)

// forceStartupHead sends a force head FCU to the execution client.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) forceStartupHead(
	ctx context.Context,
	st BeaconStateT,
) {
	slot, err := st.GetSlot()
	if err != nil {
		s.logger.Error(
			"failed to get slot for force startup head",
			"error", err,
		)
		return
	}

	// TODO: Verify if the slot number is correct here, I believe in current
	// form
	// it should be +1'd. Not a big deal until hardforks are in play though.
	if err = s.localPayloadBuilder.SendForceHeadFCU(ctx, st, slot+1); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error) {
	// Get the payload for the block.
	envelope, err := s.localPayloadBuilder.
		RetrievePayload(
			ctx,
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
		)
	if err != nil {
		s.metrics.failedToRetrievePayload(
			blk.GetSlot(),
			err,
		)

		// The latest execution payload header will be from the previous block
		// during the block building phase.
		var lph *types.ExecutionPayloadHeader
		lph, err = st.GetLatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}

		// If we failed to retrieve the payload, request a synchrnous payload.
		//
		// NOTE: The state here is properly configured by the
		// prepareStateForBuilding
		//
		// call that needs to be called before requesting the Payload.
		// TODO: We should decouple the PayloadBuilder from BeaconState to make
		// this less confusing.
		return s.localPayloadBuilder.RequestPayloadSync(
			ctx,
			st,
			blk.GetSlot(),
			// TODO: this is hood.
			max(
				//#nosec:G701
				uint64(time.Now().Unix()+1),
				uint64((lph.GetTimestamp()+1)),
			),
			blk.GetParentBlockRoot(),
			lph.GetBlockHash(),
			lph.GetParentHash(),
		)
	}
	return envelope, nil
}
