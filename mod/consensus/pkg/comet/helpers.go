package cometbft

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	abci "github.com/cometbft/cometbft/abci/types"
)

func convertValidatorUpdates(
	updates transition.ValidatorUpdates,
) []abci.ValidatorUpdate {
	valUpdates := make([]abci.ValidatorUpdate, len(updates))

	for i, update := range updates {
		valUpdates[i] = abci.ValidatorUpdate{
			PubKeyBytes: update.Pubkey[:],
			PubKeyType:  crypto.CometBLSType,
			//#nosec:G701 // this is safe.
			Power: int64(update.EffectiveBalance.Unwrap()),
		}
	}

	return valUpdates
}
