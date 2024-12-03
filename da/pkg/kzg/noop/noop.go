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

package noop

import (
	"github.com/berachain/beacon-kit/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/primitives/pkg/eip4844"
)

// Verifier is a no-op KZG proof verifier.
type Verifier struct{}

// NewVerifier creates a new GoKZGVerifier.
func NewVerifier() *Verifier {
	return &Verifier{}
}

// GetImplementation returns the implementation of the verifier.
func (v Verifier) GetImplementation() string {
	return "noop"
}

// VerifyBlobProof is a no-op.
func (v Verifier) VerifyBlobProof(
	*eip4844.Blob,
	eip4844.KZGProof,
	eip4844.KZGCommitment,
) error {
	return nil
}

// VerifyBlobProofBatch is a no-op.
func (v Verifier) VerifyBlobProofBatch(
	*types.BlobProofArgs,
) error {
	return nil
}
