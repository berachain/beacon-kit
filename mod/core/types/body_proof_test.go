package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/berachain/beacon-kit/mod/tree/merkleize"
	"github.com/stretchr/testify/require"
)

func TestProofBs(t *testing.T) {
	commitments := []kzg.Commitment{
		{0x01},
		{0x02},
		{0x03},
		{0x04},
	}

	commitmentsLeaves := types.LeavesFromCommitments(commitments)

	htr1 := merkleize.Vector(commitmentsLeaves, uint64(len(commitmentsLeaves)))
	// require.NoError(t, err)

	htr2, err := types.GetBlobKzgCommitmentsRoot(commitments)
	require.NoError(t, err)

	require.Equal(t, htr1, htr2)
}
