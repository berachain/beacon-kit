package attributes

import (
	"github.com/berachain/beacon-kit/mod/log"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type AttributesFactory[
	BeaconStateT BeaconState[WithdrawalT], WithdrawalT any,
] struct {
	// chainSpec is the chain spec.
	chainSpec primitives.ChainSpec
	// logger is the logger.
	logger log.Logger[any]
	// suggestedFeeRecipient
	suggestedFeeRecipient common.ExecutionAddress
}

// NewAttributesFactory creates a new instance of AttributesFactory.
func NewAttributesFactory[
	BeaconStateT BeaconState[WithdrawalT], WithdrawalT any,
](
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	suggestedFeeRecipient common.ExecutionAddress,
) *AttributesFactory[BeaconStateT, WithdrawalT] {
	return &AttributesFactory[BeaconStateT, WithdrawalT]{
		chainSpec:             chainSpec,
		logger:                logger,
		suggestedFeeRecipient: suggestedFeeRecipient,
	}
}

// CreateAttributes creates a new instance of PayloadAttributes.
func (f *AttributesFactory[BeaconStateT, WithdrawalT]) BuildPayloadAttributes(
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (engineprimitives.PayloadAttributer, error) {
	var (
		prevRandao [32]byte
		epoch      = f.chainSpec.SlotToEpoch(slot)
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		f.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}

	// Get the previous randao mix.
	if prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % f.chainSpec.EpochsPerHistoricalVector(),
	); err != nil {
		return nil, err
	}

	return engineprimitives.NewPayloadAttributes(
		f.chainSpec.ActiveForkVersionForEpoch(epoch),
		timestamp,
		prevRandao,
		f.suggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}
