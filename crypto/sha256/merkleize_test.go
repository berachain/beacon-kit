// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package sha256_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
	bitfield "github.com/prysmaticlabs/go-bitfield"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"
)

func Test_MerkleizeVectorSSZ(t *testing.T) {
	t.Run("empty vector", func(t *testing.T) {
		attList := make([]*ethpb.Attestation, 0)
		expected := tree.Root{
			83, 109, 152, 131, 127, 45, 209, 101, 165, 93, 94,
			234, 233, 20, 133, 149, 68, 114, 213, 111, 36, 109,
			242, 86, 191, 60, 174, 25, 53, 42, 18, 60}
		length := uint64(16)
		root, err := sha256.BuildMerkleTree(attList, length)
		require.NoError(t, err)
		require.Equal(t, expected, root)
	})
	t.Run("non empty vector", func(t *testing.T) {
		sig := make([]byte, 96)
		br := make([]byte, 32)
		attList := make([]*ethpb.Attestation, 1)
		attList[0] = &ethpb.Attestation{
			AggregationBits: bitfield.Bitlist{0x01},
			Data: &ethpb.AttestationData{
				BeaconBlockRoot: br,
				Source: &ethpb.Checkpoint{
					Root: br,
				},
				Target: &ethpb.Checkpoint{
					Root: br,
				},
			},
			Signature: sig,
		}
		expected := tree.Root{
			199, 186, 55, 142, 200, 75, 219, 191, 66, 153, 100,
			181, 200, 15, 143, 160, 25, 133, 105, 26, 183, 107,
			10, 198, 232, 231, 107, 162, 243, 243, 56, 20}
		length := uint64(16)
		root, err := sha256.BuildMerkleTree(attList, length)
		require.NoError(t, err)
		require.Equal(t, expected, root)
	})
}
