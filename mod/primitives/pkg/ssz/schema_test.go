package ssz_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	. "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/stretchr/testify/require"
)

func Test_Schema(t *testing.T) {
	beaconState := &deneb.BeaconState{}
	container, err := NewContainer(beaconState)
	require.NoError(t, err)
	require.NotNil(t, container)
}
