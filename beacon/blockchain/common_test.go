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

package blockchain_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/crypto"
	testutils "github.com/berachain/beacon-kit/testing/utils"
	"github.com/stretchr/testify/require"
)

// fakeABCIRequest implements encoding.ABCIRequest.
type fakeABCIRequest struct {
	txs    [][]byte
	height int64
}

func (r *fakeABCIRequest) GetTxs() [][]byte   { return r.txs }
func (r *fakeABCIRequest) GetTime() time.Time { return time.Unix(1, 0) }
func (r *fakeABCIRequest) GetHeight() int64   { return r.height }

// newParseTestService builds a Service with just enough wiring to exercise
// ParseBeaconBlock. The devnet spec enables blob consensus at height 2.
func newParseTestService(t *testing.T) (*blockchain.Service, chain.Spec) {
	t.Helper()
	cs, err := chain.NewSpec(spec.DevnetChainSpecData())
	require.NoError(t, err)
	require.Equal(t, int64(2), cs.BlobConsensusEnableHeight())
	svc := blockchain.NewService(
		nil, // storage backend unused by ParseBeaconBlock
		nil, // blob processor unused
		nil, // blob requester unused
		nil, // blob reconstructor unused
		noopBlobFetcher{},
		nil, // deposit contract unused
		log.NewTestLogger(t),
		cs,
		nil, // execution engine unused
		nil, // local builder unused
		nil, // state processor unused
		metrics.NewNoOpTelemetrySink(),
	)
	return svc, cs
}

// makeTestProposal returns the SSZ bytes of a valid signed block and an empty
// sidecars tx for the devnet's active fork version.
func makeTestProposal(t *testing.T, cs chain.Spec) ([]byte, []byte) {
	t.Helper()

	forkVersion := cs.ActiveForkVersionForTimestamp(1)
	blk := testutils.GenerateValidBeaconBlock(t, forkVersion)
	signedBlk := &ctypes.SignedBeaconBlock{
		BeaconBlock: blk,
		Signature:   crypto.BLSSignature{0x42},
	}

	blkBz, err := signedBlk.MarshalSSZ()
	require.NoError(t, err)

	sidecars := generateTestSidecarsBytes(t)
	return blkBz, sidecars
}

func generateTestSidecarsBytes(t *testing.T) []byte {
	t.Helper()
	sidecars := datypes.BlobSidecars{}
	bz, err := sidecars.MarshalSSZ()
	require.NoError(t, err)
	return bz
}

// ParseBeaconBlock honors the two-tx layout below the blob consensus enable
// height.
func TestParseBeaconBlock_LegacyLayout(t *testing.T) {
	t.Parallel()
	s, cs := newParseTestService(t)
	blkBz, sidecarsBz := makeTestProposal(t, cs)

	// Height 1 is below the devnet enable height (2): two txs expected.
	signedBlk, sidecars, err := s.ParseBeaconBlock(&fakeABCIRequest{
		txs:    [][]byte{blkBz, sidecarsBz},
		height: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, signedBlk)
	require.NotNil(t, sidecars)

	// A single tx is not a valid legacy proposal.
	_, _, err = s.ParseBeaconBlock(&fakeABCIRequest{
		txs:    [][]byte{blkBz},
		height: 1,
	})
	require.Error(t, err)

	// Three txs are never valid.
	_, _, err = s.ParseBeaconBlock(&fakeABCIRequest{
		txs:    [][]byte{blkBz, sidecarsBz, sidecarsBz},
		height: 1,
	})
	require.ErrorIs(t, err, blockchain.ErrTooManyConsensusTxs)
}

// ParseBeaconBlock enforces the single-tx layout at and above the blob
// consensus enable height; the sidecars slot is nil (distributed via the
// blob reactor).
func TestParseBeaconBlock_BlobConsensusLayout(t *testing.T) {
	t.Parallel()
	s, cs := newParseTestService(t)
	blkBz, sidecarsBz := makeTestProposal(t, cs)

	signedBlk, sidecars, err := s.ParseBeaconBlock(&fakeABCIRequest{
		txs:    [][]byte{blkBz},
		height: 2,
	})
	require.NoError(t, err)
	require.NotNil(t, signedBlk)
	require.Nil(t, sidecars)

	// A proposer stuffing a second tx after the transition must be rejected.
	_, _, err = s.ParseBeaconBlock(&fakeABCIRequest{
		txs:    [][]byte{blkBz, sidecarsBz},
		height: 2,
	})
	require.ErrorIs(t, err, blockchain.ErrTooManyConsensusTxs)
}
