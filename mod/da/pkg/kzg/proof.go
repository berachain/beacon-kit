// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package kzg

import (
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/gokzg"
	kzgtypes "github.com/berachain/beacon-kit/mod/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/da/kzg"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/da/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// NewBlobProofVerifier creates a new BlobVerifier with the given
// implementation.
func NewBlobProofVerifier(
	impl string,
	ts *gokzg4844.JSONTrustedSetup,
) (kzg.BlobProofVerifier[*kzgtypes.BlobProofArgs], error) {
	switch impl {
	case gokzg.Implementation:
		return gokzg.NewVerifier(ts)
	case ckzg.Implementation:
		return ckzg.NewVerifier(ts)
	default:
		return nil, errors.Wrapf(
			ErrUnsupportedKzgImplementation,
			"supplied: %s, supported: %s, %s",
			impl, gokzg.Implementation, ckzg.Implementation,
		)
	}
}

// ArgsFromSidecars converts a BlobSidecars to a slice of BlobProofArgs.
func ArgsFromSidecars[
	BeaconBlockHeaderT any,
	BlobSidecarT types.BlobSidecar[BlobSidecarT, BeaconBlockHeaderT],
	BlobSidecarsT types.BlobSidecars[BlobSidecarsT, BlobSidecarT],
](
	scs BlobSidecarsT,
) *kzgtypes.BlobProofArgs {
	proofArgs := &kzgtypes.BlobProofArgs{
		Blobs:       make([]*eip4844.Blob, scs.Len()),
		Proofs:      make([]eip4844.KZGProof, scs.Len()),
		Commitments: make([]eip4844.KZGCommitment, scs.Len()),
	}
	for i := range uint32(scs.Len()) {
		//nosec: G701 // definitively no error here
		sidecar, _ := scs.Get(i)
		proofArgs.Blobs[i] = sidecar.GetBlob()
		proofArgs.Proofs[i] = sidecar.GetProof()
		proofArgs.Commitments[i] = sidecar.GetCommitment()
	}
	return proofArgs
}
