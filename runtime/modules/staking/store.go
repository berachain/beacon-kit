package staking

import (
	sdkcollections "cosmossdk.io/collections"
	stakingtypes "cosmossdk.io/x/staking/types"
)

// Validators key: valAddr | value: Validator
func (k *Keeper) ValidatorsByValAddress() sdkcollections.Map[[]byte, stakingtypes.Validator] {
	return k.stakingKeeper.Validators
}
