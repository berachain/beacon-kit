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

//go:build ckzg

package ckzg

import (
	"unsafe"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v Verifier) VerifyBlobProof(
	blob *eip4844.Blob,
	proof eip4844.KZGProof,
	commitment eip4844.KZGCommitment,
) error {
	if valid, err := ckzg4844.VerifyBlobKZGProof(
		(*ckzg4844.Blob)(blob),
		(ckzg4844.Bytes48)(commitment),
		(ckzg4844.Bytes48)(proof),
	); err != nil {
		return err
	} else if !valid {
		return ErrInvalidProof
	}
	return nil
}

// VerifyBlobProofBatch verifies the KZG proof that the polynomial represented
// by the blob evaluated at the given point is the claimed value.
// It is more efficient than VerifyBlobProof when verifying multiple proofs.
func (v Verifier) VerifyBlobProofBatch(
	args *types.BlobProofArgs,
) error {
	blobs := make([]ckzg4844.Blob, len(args.Blobs))
	for i := range args.Blobs {
		blobs[i] = *(*ckzg4844.Blob)(args.Blobs[i])
	}

	ok, err := ckzg4844.VerifyBlobKZGProofBatch(blobs,
		*(*[]ckzg4844.Bytes48)(unsafe.Pointer(&args.Commitments)),
		*(*[]ckzg4844.Bytes48)(unsafe.Pointer(&args.Proofs)))
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidProof
	}
	return nil
}
