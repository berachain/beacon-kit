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

package engineprimitives

import "github.com/berachain/beacon-kit/primitives/eip4844"

// Compile-time assertion to ensure BlobsBundleV1 implements BlobsBundle.
var _ BlobsBundle = (*BlobsBundleV1)(nil)

// BlobsBundle is an interface for the blobs bundle.
type BlobsBundle interface {
	// GetCommitments returns the commitments in the blobs bundle.
	GetCommitments() []eip4844.KZGCommitment
	// GetProofs returns the proofs in the blobs bundle.
	GetProofs() []eip4844.KZGProof
	// GetBlobs returns the blobs in the blobs bundle.
	GetBlobs() []*eip4844.Blob
}

// BlobsBundleV1 represents a collection of commitments, proofs, and blobs.
// Each field is a slice of bytes that are serialized for transmission or
// storage.
type BlobsBundleV1 struct {
	// Commitments are the KZG commitments included in the bundle.
	Commitments []eip4844.KZGCommitment `json:"commitments"`
	// Proofs are the KZG proofs corresponding to the commitments.
	Proofs []eip4844.KZGProof `json:"proofs"`
	// Blobs are arbitrary data blobs included in the bundle.
	Blobs []*eip4844.Blob `json:"blobs"`
}

// GetCommitments returns the slice of commitments in the bundle.
func (b *BlobsBundleV1) GetCommitments() []eip4844.KZGCommitment {
	return b.Commitments
}

// GetProofs returns the slice of proofs in the bundle.
func (b *BlobsBundleV1) GetProofs() []eip4844.KZGProof {
	return b.Proofs
}

// GetBlobs returns the slice of data blobs in the bundle.
func (b *BlobsBundleV1) GetBlobs() []*eip4844.Blob {
	return b.Blobs
}
