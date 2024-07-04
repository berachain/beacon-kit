package ssz_test

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/tree"
)

func Test_Schema(t *testing.T) {
	beaconState := &deneb.BeaconState{
		Validators: []*types.Validator{{Slashed: true}},
	}
	container := beaconState.Default(spec.DevnetChainSpec())

	gi := container.GIndex(
		math.U64(1),
		tree.ObjectPath("validators/0/slashed"),
	)
	// TODO assert the generalized index
	// TODO assert invalidation on base struct carries to container construction
	fmt.Println(gi)
}
