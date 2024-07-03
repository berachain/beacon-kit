package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	. "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

func Test_Schema(t *testing.T) {
	beaconState := &deneb.BeaconState{}
	x := Walk(func(st *deneb.BeaconState) *types.Fork { return st.Fork })
}
