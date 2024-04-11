// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package types

import (
	primitives "github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
)

// BlobSidecar as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/p2p-interface.md?ref=bankless.ghost.io#blobsidecar
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen -path ./sidecar.go -objs BlobSidecar -include ../../primitives/kzg,../../primitives,$GETH_PKG_INCLUDE/common -output sidecar.ssz.go
type BlobSidecar struct {
	// Index represents the index of the blob in the block.
	Index uint64
	// Blob represents the blob data.
	// TODO: Wrangle fastssz to allow us to use a pointer here to avoid
	// copying the blob around all the time. Benchmark this as well.
	Blob kzg.Blob
	// KzgCommitment is the KZG commitment of the blob.
	KzgCommitment kzg.Commitment
	// Kzg proof allows folr the verification of the KZG commitment.
	KzgProof kzg.Proof
	// BeaconBlockHeader represents the beacon block header for which this blob
	// is being included.
	BeaconBlockHeader *primitives.BeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block.
	InclusionProof [][32]byte `ssz-size:"8,32"`
}
