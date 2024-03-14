package core

import (
	"fmt"

	"github.com/berachain/beacon-kit/beacon/core/state"
	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/config"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
type StateProcessor struct {
	cfg *config.Beacon
	st  state.BeaconState
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor(
	cfg *config.Beacon,
	st state.BeaconState,
) *StateProcessor {
	return &StateProcessor{
		cfg: cfg,
		st:  st,
	}
}

// ProcessBlock processes the block and ensures it matches the local state.
func (sp *StateProcessor) ProcessBlock(
	blk types.BeaconBlock,
) error {

	// Ensure Body is non nil.
	body := blk.GetBody()
	if body.IsNil() {
		return types.ErrNilBlkBody
	}

	// process the eth1 vote.
	payload := body.GetExecutionPayload()
	if payload.IsNil() {
		return types.ErrNilPayloadInBlk
	}

	// common.ProcessHeader

	// process the withdrawals.
	if err := sp.processWithdrawals(payload.GetWithdrawals()); err != nil {
		return err
	}

	// phase0.ProcessProposerSlashings
	// phase0.ProcessAttesterSlashings

	// process the randao reveal.
	if err := sp.processRandaoReveal(); err != nil {
		return err
	}

	// phase0.ProcessEth1Vote ? forkchoice?

	// process the deposits and ensure they match the local state.
	if err := sp.processDeposits(body.GetDeposits()); err != nil {
		return err
	}

	// ProcessVoluntaryExits

	return nil
}

// ProcessDeposits processes the deposits and ensures they match the
// local state.
func (sp *StateProcessor) processDeposits(
	deposits []*types.Deposit,
) error {
	if uint64(len(deposits)) > sp.cfg.Limits.MaxDepositsPerBlock {
		return fmt.Errorf(
			"too many deposits, expected: %d, got: %d",
			sp.cfg.Limits.MaxDepositsPerBlock, len(deposits),
		)
	}

	// Ensure the deposits match the local state.
	localDeposits, err := sp.st.ExpectedDeposits(uint64(len(deposits)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, dep := range deposits {
		if dep == nil {
			return types.ErrNilDeposit
		}
		if dep.Index != localDeposits[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localDeposits[i].Index, dep.Index)
		}
	}
	return nil
}

func (sp *StateProcessor) processWithdrawals(
	withdrawals []*enginetypes.Withdrawal,
) error {
	return nil
}

func (sp *StateProcessor) processRandaoReveal() error {
	return nil
}
