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

package noop_test

import (
	"testing"

	"github.com/berachain/beacon-kit/da/kzg/noop"
	"github.com/berachain/beacon-kit/da/kzg/types"
	"github.com/berachain/beacon-kit/primitives/eip4844"
)

func TestVerifyBlobProof(t *testing.T) {
	t.Parallel()
	verifier := noop.NewVerifier()
	err := verifier.VerifyBlobProof(
		&eip4844.Blob{},
		eip4844.KZGProof{},
		eip4844.KZGCommitment{},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVerifyBlobProofBatch(t *testing.T) {
	t.Parallel()
	verifier := noop.NewVerifier()
	args := &types.BlobProofArgs{
		Blobs:       []*eip4844.Blob{{}},
		Proofs:      []eip4844.KZGProof{{}},
		Commitments: []eip4844.KZGCommitment{{}},
	}
	err := verifier.VerifyBlobProofBatch(args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
