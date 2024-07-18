package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// GetGenesis returns the genesis state of the beacon chain.
func (h Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetGenesis(ctx context.Context) (common.Root, error) {
	// needs genesis_time and gensis_fork_version
	st, err := h.StateFromContext(ctx, "stateID")
	if err != nil {
		return common.Root{}, err
	}
	return st.GetGenesisValidatorsRoot()
}
