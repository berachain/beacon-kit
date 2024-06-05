package validator

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives"
)

// verifyIncomingBlockStateRoot verifies the state root of an incoming block
// and logs the process.
//
//nolint:gocognit // todo fix.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) ReceiveBeaconBlock(
	ctx context.Context,
	blk BeaconBlockT,
) error {
	// Grab a copy of the state to verify the incoming block.
	st := s.bsb.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	//
	// TODO: This is a super hacky. It should be handled better elsewhere,
	// ideally via some broader sync service.
	s.forceStartupSyncOnce.Do(func() { s.forceStartupHead(ctx, st) })

	// If the block is nil or a nil pointer, exit early.
	if blk.IsNil() {
		s.logger.Error(
			"aborting block verification on nil block ⛔️ ",
		)

		if s.localPayloadBuilder.Enabled() &&
			s.cfg.EnableOptimisticPayloadBuilds {
			go func() {
				if pErr := s.rebuildPayloadForRejectedBlock(
					ctx, st,
				); pErr != nil {
					s.logger.Error(
						"failed to rebuild payload for nil block",
						"error", pErr,
					)
				}
			}()
		}

		return ErrNilBlk
	}

	s.logger.Info(
		"received incoming beacon block 📫 ",
		"state_root", blk.GetStateRoot(),
	)

	// We purposefully make a copy of the BeaconState in orer
	// to avoid modifying the underlying state, for the event in which
	// we have to rebuild a payload for this slot again, if we do not agree
	// with the incoming block.
	stCopy := st.Copy()

	// Verify the state root of the incoming block.
	if err := s.verifyStateRoot(
		ctx, stCopy, blk,
	); err != nil {
		// TODO: this is expensive because we are not caching the
		// previous result of HashTreeRoot().
		localStateRoot, htrErr := st.HashTreeRoot()
		if htrErr != nil {
			return htrErr
		}

		s.logger.Error(
			"rejecting incoming block ❌ ",
			"block_state_root",
			blk.GetStateRoot(),
			"local_state_root",
			primitives.Root(localStateRoot),
			"error",
			err,
		)

		if s.localPayloadBuilder.Enabled() &&
			s.cfg.EnableOptimisticPayloadBuilds {
			go func() {
				if pErr := s.rebuildPayloadForRejectedBlock(
					ctx, st,
				); pErr != nil {
					s.logger.Error(
						"failed to rebuild payload for rejected block",
						"for_slot", blk.GetSlot(),
						"error", pErr,
					)
				}
			}()
		}

		return err
	}

	s.logger.Info(
		"state root verification succeeded - accepting incoming block 🏎️ ",
		"state_root", blk.GetStateRoot(),
	)

	if s.localPayloadBuilder.Enabled() && s.cfg.EnableOptimisticPayloadBuilds {
		go func() {
			if err := s.optimisticPayloadBuild(ctx, stCopy, blk); err != nil {
				s.logger.Error(
					"failed to build optimistic payload",
					"for_slot", blk.GetSlot()+1,
					"error", err,
				)
			}
		}()
	}

	return nil
}

// ReceiveBlobs receives blobs from the network and processes them.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositT,
]) ReceiveBlobs(
	ctx context.Context,
	blk BeaconBlockT,
	blobs BlobSidecarsT,
) error {
	if blk.IsNil() {
		s.logger.Error(
			"aborting blob verification on nil block ⛔️ ",
		)
		return ErrNilBlk
	}

	s.logger.Info(
		"received incoming blob sidecars 🚔 ",
		"state_root", blk.GetStateRoot(),
	)

	if err := s.verifyBlobProofs(blk.GetSlot(), blobs); err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars ❌ ",
			"error", err,
		)
		return err
	}

	s.logger.Info(
		"blob sidecars verification succeeded - accepting incoming blobs 💦 ",
		"num_blobs", blobs.Len(),
	)
	return nil
}
