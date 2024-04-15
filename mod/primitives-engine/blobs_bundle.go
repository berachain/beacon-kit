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

package engineprimitives

import "github.com/berachain/beacon-kit/mod/primitives/kzg"

// BlobsBundleV1 represents a collection of commitments, proofs, and blobs.
// Each field is a slice of bytes that are serialized for transmission or
// storage.
type BlobsBundleV1 struct {
	// Commitments are the KZG commitments included in the bundle.
	Commitments []kzg.Commitment `json:"commitments"`
	// Proofs are the KZG proofs corresponding to the commitments.
	Proofs []kzg.Proof `json:"proofs"`
	// Blobs are arbitrary data blobs included in the bundle.
	Blobs []*kzg.Blob `json:"blobs"`
}

// GetCommitments returns the slice of commitments in the bundle.
func (b *BlobsBundleV1) GetCommitments() []kzg.Commitment {
	return b.Commitments
}

// GetProofs returns the slice of proofs in the bundle.
func (b *BlobsBundleV1) GetProofs() []kzg.Proof {
	return b.Proofs
}

// GetBlobs returns the slice of data blobs in the bundle.
func (b *BlobsBundleV1) GetBlobs() []*kzg.Blob {
	return b.Blobs
}
