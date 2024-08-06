package state_transition

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/sourcegraph/conc/iter"
)

// processSyncCommitteeUpdates processes the sync committee updates.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, ValidatorT, _, _, _,
]) processSyncCommitteeUpdates(
	st BeaconStateT,
) (transition.ValidatorUpdates, error) {
	vals, err := st.GetValidatorsByEffectiveBalance()
	if err != nil {
		return nil, err
	}

	return iter.MapErr(
		vals,
		func(val *ValidatorT) (*transition.ValidatorUpdate, error) {
			v := (*val)
			return &transition.ValidatorUpdate{
				Pubkey:           v.GetPubkey(),
				EffectiveBalance: v.GetEffectiveBalance(),
			}, nil
		},
	)
}
