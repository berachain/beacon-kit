// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blob_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"testing"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/kzg/gokzg"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/version"
	testutils "github.com/berachain/beacon-kit/testing/utils"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	gethengine "github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/stretchr/testify/require"
)

// fakeELBlobs is a canned engine_getBlobsV2 backend recording what was asked.
type fakeELBlobs struct {
	resp      []*gethengine.BlobAndProofV2
	err       error
	gotHashes []common.ExecutionHash
}

func (f *fakeELBlobs) GetBlobsV2(
	_ context.Context, hashes []common.ExecutionHash,
) ([]*gethengine.BlobAndProofV2, error) {
	f.gotHashes = hashes
	return f.resp, f.err
}

// The verifier doubles as the prover and is expensive to construct, so it is shared across tests.
var (
	kzgOnce     sync.Once
	kzgVerifier *gokzg.Verifier
)

func loadKZG(t *testing.T) *gokzg.Verifier {
	t.Helper()
	kzgOnce.Do(func() {
		bz, err := os.ReadFile("../../testing/files/kzg-trusted-setup.json")
		require.NoError(t, err)
		ts := new(gokzg4844.JSONTrustedSetup)
		require.NoError(t, json.Unmarshal(bz, ts))
		kzgVerifier, err = gokzg.NewVerifier(ts)
		require.NoError(t, err)
	})
	require.NotNil(t, kzgVerifier)
	return kzgVerifier
}

// reconstructorEnv builds a signed block committing to one real blob, and a reconstructor backed by a fake EL
// and the real KZG prover and sidecar factory.
func reconstructorEnv(t *testing.T, el *fakeELBlobs) (*dablob.Reconstructor, *ctypes.SignedBeaconBlock, *eip4844.Blob) {
	t.Helper()
	verifier := loadKZG(t)

	blob := &eip4844.Blob{}
	blob[1] = 0x42 // arbitrary content; any canonical field elements will do

	commitment, err := verifier.Context.BlobToKZGCommitment((*gokzg4844.Blob)(blob), 0)
	require.NoError(t, err)

	blk := testutils.GenerateValidBeaconBlock(t, version.Deneb())
	blk.GetBody().SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{eip4844.KZGCommitment(commitment)})
	signedBlk := &ctypes.SignedBeaconBlock{BeaconBlock: blk, Signature: crypto.BLSSignature{0xab}}

	reconstructor := dablob.NewReconstructor(el, dablob.NewSidecarFactory(metrics.NewNoOpTelemetrySink()), verifier, log.NewNopLogger())
	return reconstructor, signedBlk, blob
}

// The reconstructed sidecars must be canonical: bound to the block's header and signature, carrying valid
// inclusion proofs, and a locally computed blob proof that verifies against the block's commitment.
func TestReconstructor_RebuildsCanonicalSidecars(t *testing.T) {
	t.Parallel()
	el := &fakeELBlobs{}
	reconstructor, signedBlk, blob := reconstructorEnv(t, el)
	el.resp = []*gethengine.BlobAndProofV2{{Blob: blob[:]}}

	sidecars, err := reconstructor.ReconstructSidecars(t.Context(), signedBlk)
	require.NoError(t, err)
	require.Len(t, sidecars, 1)

	commitments := signedBlk.GetBeaconBlock().GetBody().GetBlobKzgCommitments()
	require.Equal(t, commitments.ToVersionedHashes(), el.gotHashes, "must query the EL by the block's versioned hashes")

	sc := sidecars[0]
	require.Equal(t, commitments[0], sc.GetKzgCommitment())
	require.Equal(t, signedBlk.GetBeaconBlock().GetHeader().HashTreeRoot(), sc.GetBeaconBlockHeader().HashTreeRoot())
	require.Equal(t, signedBlk.GetSignature(), sc.GetSignature(), "the block signature must be embedded")
	require.True(t, sc.HasValidInclusionProof(), "inclusion proof must bind the commitment to the header")
	require.NoError(t,
		loadKZG(t).VerifyBlobProof(blob, sc.GetKzgProof(), sc.GetKzgCommitment()),
		"the recomputed blob proof must verify against the commitment")
}

// Every malformed or incomplete EL response is rejected: engine_getBlobsV2 is all-or-nothing, so a nil/empty
// response, a count mismatch, a nil element, or a truncated blob must never yield sidecars, and a transport
// error must propagate rather than read as a miss.
func TestReconstructor_BadResponsesRejected(t *testing.T) {
	t.Parallel()
	el := &fakeELBlobs{}
	reconstructor, signedBlk, blob := reconstructorEnv(t, el)

	cases := []struct {
		name    string
		resp    []*gethengine.BlobAndProofV2
		err     error
		wantErr string
	}{
		{name: "empty response is a miss", resp: nil, wantErr: dablob.ErrBlobsNotInELPool.Error()},
		{name: "count mismatch", resp: []*gethengine.BlobAndProofV2{{Blob: blob[:]}, {Blob: blob[:]}},
			wantErr: dablob.ErrBlobsNotInELPool.Error()},
		{name: "nil element", resp: []*gethengine.BlobAndProofV2{nil}, wantErr: dablob.ErrBlobsNotInELPool.Error()},
		{name: "truncated blob", resp: []*gethengine.BlobAndProofV2{{Blob: blob[:100]}}, wantErr: "invalid length"},
		{name: "EL error propagates", err: errors.New("engine down"), wantErr: "engine down"},
	}
	for _, tc := range cases {
		el.resp, el.err = tc.resp, tc.err
		_, err := reconstructor.ReconstructSidecars(t.Context(), signedBlk)
		require.ErrorContains(t, err, tc.wantErr, tc.name)
	}
}

// A block committing to no blobs never touches the EL.
func TestReconstructor_NoCommitmentsSkipsEL(t *testing.T) {
	t.Parallel()
	el := &fakeELBlobs{}
	reconstructor, signedBlk, _ := reconstructorEnv(t, el)
	signedBlk.GetBeaconBlock().GetBody().SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash]{})

	sidecars, err := reconstructor.ReconstructSidecars(t.Context(), signedBlk)
	require.NoError(t, err)
	require.Empty(t, sidecars)
	require.Nil(t, el.gotHashes, "the EL must not be queried")
}
