package state_transition

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// processSlashingsReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#slashings-balances-updates
//
//nolint:lll
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processSlashingsReset(
	st BeaconStateT,
) error {
	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	index := (uint64(sp.cs.SlotToEpoch(slot)) + 1) % sp.cs.EpochsPerSlashingsVector()
	return st.UpdateSlashingAtIndex(index, 0)
}

// processProposerSlashing as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#proposer-slashings
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processProposerSlashing(
	_ BeaconStateT,
	// ps ProposerSlashing,
) error {
	return nil
}

// processAttesterSlashing as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#attester-slashings
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processAttesterSlashing(
	_ BeaconStateT,
	// as AttesterSlashing,
) error {
	return nil
}

// processSlashings as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#slashings
//
// processSlashings processes the slashings and ensures they match the local
// state.
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processSlashings(
	st BeaconStateT,
) error {
	totalBalance, err := st.GetTotalActiveBalances(sp.cs.SlotsPerEpoch())
	if err != nil {
		return err
	}

	totalSlashings, err := st.GetTotalSlashing()
	if err != nil {
		return err
	}

	adjustedTotalSlashingBalance := min(
		uint64(totalSlashings)*sp.cs.ProportionalSlashingMultiplier(),
		uint64(totalBalance),
	)

	vals, err := st.GetValidators()
	if err != nil {
		return err
	}

	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	//nolint:mnd // this is in the spec
	slashableEpoch := (uint64(sp.cs.SlotToEpoch(slot)) + sp.cs.EpochsPerSlashingsVector()) / 2

	// Iterate through the validators and slash if needed.
	for _, val := range vals {
		if val.IsSlashed() &&
			(slashableEpoch == uint64(val.GetWithdrawableEpoch())) {
			if err = sp.processSlash(
				st, val,
				adjustedTotalSlashingBalance,
				uint64(totalBalance),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

// processSlash handles the logic for slashing a validator.
//
//nolint:unused // will be used later
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, ValidatorT, _, _, _,
]) processSlash(
	st BeaconStateT,
	val ValidatorT,
	adjustedTotalSlashingBalance uint64,
	totalBalance uint64,
) error {
	// Calculate the penalty.
	increment := sp.cs.EffectiveBalanceIncrement()
	balDivIncrement := uint64(val.GetEffectiveBalance()) / increment
	penaltyNumerator := balDivIncrement * adjustedTotalSlashingBalance
	penalty := penaltyNumerator / totalBalance * increment

	// Get the val index and decrease the balance of the validator.
	idx, err := st.ValidatorIndexByPubkey(val.GetPubkey())
	if err != nil {
		return err
	}

	return st.DecreaseBalance(idx, math.Gwei(penalty))
}
