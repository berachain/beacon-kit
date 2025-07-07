// SPDX-License-Identifier: BUSL-1.1
package proof

import (
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetValidatorBalance returns the balance of a validator along with a Merkle
// proof that can be verified against the beacon block root.
func (h *Handler) GetValidatorBalance(c handlers.Context) (any, error) {
	params, err := utils.BindAndValidate[types.ValidatorBalanceRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}

	validatorIndex, err := math.U64FromString(params.ValidatorIndex)
	if err != nil {
		return nil, err
	}

	slot, beaconState, blockHeader, err := h.resolveTimestampID(params.TimestampID)
	if err != nil {
		return nil, err
	}

	h.Logger().Info(
		"Generating validator balance proofs", "slot", slot, "validator_index", validatorIndex,
	)

	bsm, err := beaconState.GetMarshallable()
	if err != nil {
		return nil, err
	}

	balProof, beaconBlockRoot, err := merkle.ProveValidatorBalanceInBlock(
		validatorIndex, blockHeader, bsm,
	)
	if err != nil {
		return nil, err
	}

	balance, err := beaconState.GetBalance(validatorIndex)
	if err != nil {
		return nil, err
	}

	return types.ValidatorBalanceResponse{
		BeaconBlockHeader: blockHeader,
		BeaconBlockRoot:   beaconBlockRoot,
		Balance:           balance,
		BalanceProof:      balProof,
	}, nil
}
