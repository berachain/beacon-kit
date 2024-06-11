package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/stretchr/testify/require"
)

func TestBlobsBundleV1(t *testing.T) {
	bundle := &engineprimitives.BlobsBundleV1[[48]byte, [48]byte, [131072]byte]{
		Commitments: [][48]byte{{1, 2, 3}, {4, 5, 6}},
		Proofs:      [][48]byte{{7, 8, 9}, {10, 11, 12}},
		Blobs:       []*[131072]byte{{13, 14, 15}, {16, 17, 18}},
	}

	commitments := bundle.GetCommitments()
	require.Equal(t, bundle.Commitments, commitments)

	proofs := bundle.GetProofs()
	require.Equal(t, bundle.Proofs, proofs)

	blobs := bundle.GetBlobs()
	require.Equal(t, bundle.Blobs, blobs)
}
