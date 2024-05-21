package privval

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/cometbft/cometbft/types"
)

// extension of cometbft PrivValidator to allow signing arbitrary
// things.
type Privval interface {
	types.PrivValidator
	crypto.BLSSigner
}
