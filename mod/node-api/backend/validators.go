package backend

import (
	"context"
	"strconv"

	serverType "github.com/berachain/beacon-kit/mod/node-api/server/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) GetStateValidator(
	ctx context.Context,
	stateID string,
	validatorID string,
) (*serverType.ValidatorData[ValidatorT], error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return nil, err
	}
	index, indexErr := b.getValidatorIndex(st, validatorID)
	if indexErr != nil {
		return nil, indexErr
	}
	validator, validatorErr := st.ValidatorByIndex(index)
	if validatorErr != nil {
		return nil, validatorErr
	}
	balance, balanceErr := st.GetBalance(index)
	if balanceErr != nil {
		return nil, balanceErr
	}
	return &serverType.ValidatorData[ValidatorT]{
		Index:     index.Unwrap(),
		Balance:   balance.Unwrap(),
		Status:    "active",
		Validator: validator,
	}, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _,
]) GetStateValidators(
	ctx context.Context,
	stateID string,
	id []string,
	_ []string,
) ([]*serverType.ValidatorData[ValidatorT], error) {
	validators := make([]*serverType.ValidatorData[ValidatorT], 0)
	for _, indexOrKey := range id {
		validatorData, err := b.GetStateValidator(ctx, stateID, indexOrKey)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validatorData)
	}
	return validators, nil
}

func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetStateValidatorBalances(
	ctx context.Context,
	stateID string,
	id []string,
) ([]*serverType.ValidatorBalanceData, error) {
	st, err := b.StateFromContext(ctx, stateID)
	if err != nil {
		return nil, err
	}
	balances := make([]*serverType.ValidatorBalanceData, 0)
	for _, indexOrKey := range id {
		index, indexErr := b.getValidatorIndex(st, indexOrKey)
		if indexErr != nil {
			return nil, indexErr
		}
		balance, err := st.GetBalance(index)
		if err != nil {
			return nil, err
		}
		balances = append(balances, &serverType.ValidatorBalanceData{
			Index:   index.Unwrap(),
			Balance: balance.Unwrap(),
		})
	}
	return balances, nil
}
func (b Backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) getValidatorIndex(st BeaconStateT, keyOrIndex string) (math.U64, error) {
	if index, err := strconv.ParseUint(keyOrIndex, 10, 64); err == nil {
		return math.U64(index), nil
	}
	var key crypto.BLSPubkey
	if err := key.UnmarshalText([]byte(keyOrIndex)); err != nil {
		return math.U64(0), err
	}
	return st.ValidatorIndexByPubkey(key)
}
