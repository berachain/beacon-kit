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
	"testing"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/blob"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// noopSink is a no-op TelemetrySink implementation for tests.
type noopSink struct{}

func (noopSink) MeasureSince(string, time.Time, ...string) {}

// TestBuildSidecars_MismatchedBundle ensures a blobs bundle whose blobs,
// commitments and proofs have differing lengths is rejected with an error
// rather than triggering an out-of-bounds panic inside the build goroutines.
func TestBuildSidecars_MismatchedBundle(t *testing.T) {
	t.Parallel()

	factory := blob.NewSidecarFactory(noopSink{})
	signedBlk, err := ctypes.NewEmptySignedBeaconBlockWithVersion(version.Deneb())
	require.NoError(t, err)

	tests := []struct {
		name   string
		bundle *engineprimitives.BlobsBundleV1
	}{
		{
			name: "fewer commitments than blobs",
			bundle: &engineprimitives.BlobsBundleV1{
				Blobs:       []*eip4844.Blob{{}, {}},
				Commitments: []eip4844.KZGCommitment{{}},
				Proofs:      []eip4844.KZGProof{{}, {}},
			},
		},
		{
			name: "fewer proofs than blobs",
			bundle: &engineprimitives.BlobsBundleV1{
				Blobs:       []*eip4844.Blob{{}, {}},
				Commitments: []eip4844.KZGCommitment{{}, {}},
				Proofs:      []eip4844.KZGProof{{}},
			},
		},
		{
			name: "more commitments and proofs than blobs",
			bundle: &engineprimitives.BlobsBundleV1{
				Blobs:       []*eip4844.Blob{{}},
				Commitments: []eip4844.KZGCommitment{{}, {}},
				Proofs:      []eip4844.KZGProof{{}, {}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// The guard must turn the mismatch into a returned error, never a panic.
			require.NotPanics(t, func() {
				sidecars, bErr := factory.BuildSidecars(signedBlk, tt.bundle)
				require.ErrorContains(t, bErr, "mismatched blobs bundle")
				require.Nil(t, sidecars)
			})
		})
	}
}

// TestBuildSidecars_EmptyBundle ensures the length guard does not falsely
// reject a well-formed empty bundle (no blobs, commitments or proofs).
func TestBuildSidecars_EmptyBundle(t *testing.T) {
	t.Parallel()

	factory := blob.NewSidecarFactory(noopSink{})
	signedBlk, err := ctypes.NewEmptySignedBeaconBlockWithVersion(version.Deneb())
	require.NoError(t, err)

	sidecars, err := factory.BuildSidecars(signedBlk, &engineprimitives.BlobsBundleV1{})
	require.NoError(t, err)
	require.Empty(t, sidecars)
}
