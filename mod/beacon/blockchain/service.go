package blockchain

import (
	"context"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is the event handler for the blockchain service.
type Service[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
] struct {
	blockBroker           EventFeed[*asynctypes.Event[BeaconBlockT]]
	genesisBroker         EventFeed[*asynctypes.Event[GenesisT]]
	validatorUpdateBroker EventFeed[*asynctypes.Event[transition.ValidatorUpdates]]
	logger                log.Logger[any]
	processor             Processor[
		BeaconBlockT,
		BeaconBlockBodyT,
		DepositT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		GenesisT,
	]
}

// NewService creates a new event handler for the blockchain service.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
](
	processor Processor[
		BeaconBlockT,
		BeaconBlockBodyT,
		DepositT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		GenesisT,
	],
	blockBroker EventFeed[*asynctypes.Event[BeaconBlockT]],
	genesisBroker EventFeed[*asynctypes.Event[GenesisT]],
	//nolint:lll // compiler vs linter
	validatorUpdateBroker EventFeed[*asynctypes.Event[transition.ValidatorUpdates]],
	logger log.Logger[any],
) *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	DepositT,
	ExecutionPayloadT,
	ExecutionPayloadHeaderT,
	GenesisT,
] {
	return &Service[
		BeaconBlockT,
		BeaconBlockBodyT,
		DepositT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		GenesisT,
	]{
		blockBroker:           blockBroker,
		genesisBroker:         genesisBroker,
		validatorUpdateBroker: validatorUpdateBroker,
		logger:                logger,
		processor:             processor,
	}
}

// Start starts the event handler.
func (e *Service[
	BeaconBlockT, _, _, _, _, GenesisT,
]) Start(ctx context.Context) error {
	subBlkCh, err := e.blockBroker.Subscribe()
	if err != nil {
		return err
	}
	subGenCh, err := e.genesisBroker.Subscribe()
	if err != nil {
		return err
	}
	go e.start(ctx, subBlkCh, subGenCh)
	return nil
}

func (e *Service[
	BeaconBlockT, _, _, _, _, GenesisT,
]) start(
	ctx context.Context,
	subBlkCh chan *asynctypes.Event[BeaconBlockT],
	subGenCh chan *asynctypes.Event[GenesisT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-subBlkCh:
			switch msg.Type() {
			case events.BeaconBlockReceived:
				e.handleBeaconBlockReceived(msg)
			case events.BeaconBlockFinalizedRequest:
				e.handleBeaconBlockFinalization(msg)
			}
		case msg := <-subGenCh:
			if msg.Type() == events.GenesisDataProcessRequest {
				e.handleProcessGenesisDataRequest(msg)
			}
		}
	}
}

// Name returns the name of the event handler.
func (e *Service[
	_, _, _, _, _, _,
]) Name() string {
	return "blockchain"
}

func (e *Service[
	_, _, _, _, _, GenesisT,
]) handleProcessGenesisDataRequest(msg *asynctypes.Event[GenesisT]) {
	if msg.Error() != nil {
		e.logger.Error("Error processing genesis data", "error", msg.Error())
		return
	}

	// Process the genesis data.
	valUpdates, err := e.processor.ProcessGenesisData(msg.Context(), msg.Data())
	if err != nil {
		e.logger.Error("Failed to process genesis data", "error", err)
	}

	// Publish the validator set updated event.
	if err = e.validatorUpdateBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.ValidatorSetUpdated,
			valUpdates,
			err,
		),
	); err != nil {
		e.logger.Error(
			"Failed to publish validator set updated event",
			"error",
			err,
		)
	}
}

// handleBeaconBlockReceived handles the beacon block received event.
func (e *Service[
	BeaconBlockT, _, _, _, _, _,
]) handleBeaconBlockReceived(
	msg *asynctypes.Event[BeaconBlockT],
) {
	// If the block is nil, exit early.
	if msg.Error() != nil {
		e.logger.Error("Error processing beacon block", "error", msg.Error())
		return
	}

	// Publish the verified block event.
	if err := e.blockBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.BeaconBlockVerified,
			msg.Data(),
			e.processor.VerifyIncomingBlock(msg.Context(), msg.Data()),
		),
	); err != nil {
		e.logger.Error("Failed to publish verified block", "error", err)
	}
}

// handleBeaconBlockFinalization handles the beacon block finalized event.
func (e *Service[
	BeaconBlockT, _, _, _, _, _,
]) handleBeaconBlockFinalization(
	msg *asynctypes.Event[BeaconBlockT],
) {
	// If there's an error in the event, log it and return
	if msg.Error() != nil {
		e.logger.Error("Error verifying beacon block", "error", msg.Error())
		return
	}

	// Process the verified block
	valUpdates, err := e.processor.ProcessBeaconBlock(msg.Context(), msg.Data())
	if err != nil {
		e.logger.Error("Failed to process verified beacon block", "error", err)
	}

	// If required, we want to forkchoice at the end of post
	// block processing.
	// TODO: this is hood as fuck.
	// We won't send a fcu if the block is bad, should be addressed
	// via ticker later.
	if err = e.blockBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(), events.BeaconBlockFinalized, msg.Data(),
		),
	); err != nil {
		e.logger.Error("Failed to publish finalized block", "error", err)
	}

	// Publish the validator set updated event
	if err = e.validatorUpdateBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.ValidatorSetUpdated,
			valUpdates,
			err,
		)); err != nil {
		e.logger.Error(
			"Failed to publish validator set updated event",
			"error",
			err,
		)
	}
}
