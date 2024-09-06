package middleware

import (
	"context"
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

/* -------------------------------------------------------------------------- */
/*                                 InitGenesis                                */
/* -------------------------------------------------------------------------- */

// vm.ValidatorManager should be filled with the content in returned ValidatorUpdates
func (vm *VMMiddleware) InitGenesis(ctx context.Context, bz []byte) (transition.ValidatorUpdates, error) {
	waitCtx, cancel := context.WithTimeout(ctx, AwaitTimeout)
	defer cancel()

	data := new(GenesisT)
	if err := json.Unmarshal(bz, data); err != nil {
		vm.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	if err := vm.dispatcher.Publish(
		async.NewEvent(ctx, async.GenesisDataReceived, *data),
	); err != nil {
		return nil, err
	}
	return vm.waitForGenesisProcessed(waitCtx)
}

// waitForGenesisProcessed waits until the genesis data has been processed and
// returns the validator updates, or err if the context is cancelled.
func (vm *VMMiddleware) waitForGenesisProcessed(ctx context.Context) (transition.ValidatorUpdates, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInitGenesisTimeout(ctx.Err())
	case gdpEvent := <-vm.subGenDataProcessed:
		return gdpEvent.Data(), gdpEvent.Error()
	}
}

/* -------------------------------------------------------------------------- */
/*                               BuildBlock                              */
/* -------------------------------------------------------------------------- */

// BuildBlock is the internal handler for preparing proposals.
func (vm *VMMiddleware) BuildBlock(ctx context.Context, slotData SlotDataT) ([]byte, []byte, error) {
	awaitCtx, cancel := context.WithTimeout(ctx, AwaitTimeout)
	defer cancel()

	// flush the channels to ensure that we are not handling old data.
	if numMsgs := async.ClearChan(vm.subBuiltBeaconBlock); numMsgs > 0 {
		vm.logger.Error(
			"WARNING: messages remaining in built beacon block channel",
			"num_msgs", numMsgs)
	}
	if numMsgs := async.ClearChan(vm.subBuiltSidecars); numMsgs > 0 {
		vm.logger.Error(
			"WARNING: messages remaining in built sidecars channel",
			"num_msgs", numMsgs)
	}

	if err := vm.dispatcher.Publish(
		async.NewEvent(
			ctx, async.NewSlot, slotData,
		),
	); err != nil {
		return nil, nil, err
	}

	// wait for built beacon block
	builtBeaconBlock, err := vm.waitForBuiltBeaconBlock(awaitCtx)
	if err != nil {
		return nil, nil, err
	}

	// wait for built sidecars
	builtSidecars, err := vm.waitForBuiltSidecars(awaitCtx)
	if err != nil {
		return nil, nil, err
	}

	return vm.handleBuiltBeaconBlockAndSidecars(builtBeaconBlock, builtSidecars)
}

// waitForBuiltBeaconBlock waits for the built beacon block to be received.
func (vm *VMMiddleware) waitForBuiltBeaconBlock(ctx context.Context) (BeaconBlockT, error) {
	select {
	case <-ctx.Done():
		return *new(BeaconBlockT), ErrBuildBeaconBlockTimeout(ctx.Err())
	case bbEvent := <-vm.subBuiltBeaconBlock:
		return bbEvent.Data(), bbEvent.Error()
	}
}

// waitForBuiltSidecars waits for the built sidecars to be received.
func (vm *VMMiddleware) waitForBuiltSidecars(ctx context.Context) (BlobSidecarsT, error) {
	select {
	case <-ctx.Done():
		return *new(BlobSidecarsT), ErrBuildSidecarsTimeout(ctx.Err())
	case scEvent := <-vm.subBuiltSidecars:
		return scEvent.Data(), scEvent.Error()
	}
}

// handleBuiltBeaconBlockAndSidecars gossips the built beacon block and blob
// sidecars to the network.
func (vm *VMMiddleware) handleBuiltBeaconBlockAndSidecars(bb BeaconBlockT, sc BlobSidecarsT) ([]byte, []byte, error) {
	bbBz, bbErr := bb.MarshalSSZ()
	if bbErr != nil {
		return nil, nil, bbErr
	}
	scBz, scErr := sc.MarshalSSZ()
	if scErr != nil {
		return nil, nil, scErr
	}
	return bbBz, scBz, nil
}

/* -------------------------------------------------------------------------- */
/*                               VerifyBlock                              */
/* -------------------------------------------------------------------------- */

// VerifyBlock processes the proposal for the ABCI middleware.
// It handles both the beacon block and blob sidecars concurrently.
// Returns error if block does not verify, nil otherwise
func (vm *VMMiddleware) VerifyBlock(ctx context.Context, outerBlk OuterBlock) error {
	awaitCtx, cancel := context.WithTimeout(ctx, AwaitTimeout)
	defer cancel()

	// flush the channels to ensure that we are not handling old data.
	if numMsgs := async.ClearChan(vm.subBBVerified); numMsgs > 0 {
		vm.logger.Error(
			"WARNING: messages remaining in beacon block verification channel",
			"num_msgs", numMsgs)
	}
	if numMsgs := async.ClearChan(vm.subSCVerified); numMsgs > 0 {
		vm.logger.Error(
			"WARNING: messages remaining in sidecar verification channel",
			"num_msgs", numMsgs)
	}

	// Request the beacon block.
	forkVersion := vm.ActiveForkVersionForSlot(math.U64(outerBlk.Height()))
	blk, err := UnmarshalBeaconBlockFromOuterBlock(outerBlk, forkVersion)
	if err != nil {
		return errors.WrapNonFatal(err)
	}

	// notify that the beacon block has been received.
	if err = vm.dispatcher.Publish(
		async.NewEvent(ctx, async.BeaconBlockReceived, blk),
	); err != nil {
		return errors.WrapNonFatal(err)
	}

	// Request the blob sidecars.
	sidecars, err := UnmarshalBlobSidecarsFromOuterBlock(outerBlk)
	if err != nil {
		return errors.WrapNonFatal(err)
	}

	// notify that the sidecars have been received.
	if err = vm.dispatcher.Publish(
		async.NewEvent(ctx, async.SidecarsReceived, sidecars),
	); err != nil {
		return errors.WrapNonFatal(err)
	}

	// err if the built beacon block or sidecars failed verification.
	_, err = vm.waitForBeaconBlockVerification(awaitCtx)
	if err != nil {
		return err
	}
	_, err = vm.waitForSidecarVerification(awaitCtx)
	return err
}

// waitForBeaconBlockVerification waits for the built beacon block to be
// verified.
func (vm *VMMiddleware) waitForBeaconBlockVerification(ctx context.Context) (BeaconBlockT, error) {
	select {
	case <-ctx.Done():
		return *new(BeaconBlockT), ErrVerifyBeaconBlockTimeout(ctx.Err())
	case vEvent := <-vm.subBBVerified:
		return vEvent.Data(), vEvent.Error()
	}
}

// waitForSidecarVerification waits for the built sidecars to be verified.
func (vm *VMMiddleware) waitForSidecarVerification(ctx context.Context) (BlobSidecarsT, error) {
	select {
	case <-ctx.Done():
		return *new(BlobSidecarsT), ErrVerifySidecarsTimeout(ctx.Err())
	case vEvent := <-vm.subSCVerified:
		return vEvent.Data(), vEvent.Error()
	}
}

/* -------------------------------------------------------------------------- */
/*                                AcceptBlock                               */
/* -------------------------------------------------------------------------- */

// AcceptBlock returns the validator set updates from the beacon state.
func (vm *VMMiddleware) AcceptBlock(ctx context.Context, outerBlk OuterBlock) (transition.ValidatorUpdates, error) {
	awaitCtx, cancel := context.WithTimeout(ctx, AwaitTimeout)
	defer cancel()

	// flush the channel to ensure that we are not handling old data.
	if numMsgs := async.ClearChan(vm.subFinalValidatorUpdates); numMsgs > 0 {
		vm.logger.Error(
			"WARNING: messages remaining in final validator updates channel",
			"num_msgs", numMsgs)
	}

	forkVersion := vm.ActiveForkVersionForSlot(math.Slot(outerBlk.Height()))
	blk, err := UnmarshalBeaconBlockFromOuterBlock(outerBlk, forkVersion)
	if err != nil {
		// If we don't have a block, we can't do anything.
		return nil, nil
	}
	blobs, err := UnmarshalBlobSidecarsFromOuterBlock(outerBlk)
	if err != nil {
		// If we don't have a block, we can't do anything.
		return nil, nil
	}

	// notify that the final beacon block has been received.
	if err = vm.dispatcher.Publish(
		async.NewEvent(ctx, async.FinalBeaconBlockReceived, blk),
	); err != nil {
		return nil, err
	}

	// notify that the final blob sidecars have been received.
	if err = vm.dispatcher.Publish(
		async.NewEvent(ctx, async.FinalSidecarsReceived, blobs),
	); err != nil {
		return nil, err
	}

	// wait for the final validator updates.
	return vm.waitForFinalValidatorUpdates(awaitCtx)
}

// waitForFinalValidatorUpdates waits for the final validator updates to be
// received.
func (vm *VMMiddleware) waitForFinalValidatorUpdates(
	ctx context.Context,
) (transition.ValidatorUpdates, error) {
	select {
	case <-ctx.Done():
		return nil, ErrFinalValidatorUpdatesTimeout(ctx.Err())
	case event := <-vm.subFinalValidatorUpdates:
		return event.Data(), event.Error()
	}
}
