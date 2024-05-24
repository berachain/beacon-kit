package runtime

import (
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// convertValidatorUpdate abstracts the conversion of
// a transition.ValidatorUpdate  to an appmodulev2.ValidatorUpdate
//
//	format.
func convertValidatorUpdate(
	u **transition.ValidatorUpdate,
) (appmodulev2.ValidatorUpdate, error) {
	update := *u
	return appmodulev2.ValidatorUpdate{
		PubKey:     update.Pubkey[:],
		PubKeyType: crypto.CometBLSType,
		//#nosec:G701 // this is safe.
		Power: int64(update.EffectiveBalance.Unwrap()),
	}, nil
}
