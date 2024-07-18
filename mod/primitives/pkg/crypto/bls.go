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

package crypto

import "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"

const (
	// CometBLSType is the BLS curve type used in the Comet BFT consensus
	// algorithm.
	CometBLSType = "bls12_381"

	// CometBLSPower is the voting power given to a validator when they
	// are in the active set.
	// TODO: Move this, it doesn't really belong here.
	CometBLSPower = 1e10
)

//nolint:lll // link.
type (
	// BLSPubkey as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	BLSPubkey = bytes.B48

	// BLSSignature as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	BLSSignature = bytes.B96
)

// BLSSigner defines an interface for cryptographic signing operations.
// It uses generic type parameters Signature and Pubkey, both of which are
// slices of bytes.
type BLSSigner interface {
	// PublicKey returns the public key of the signer.
	PublicKey() BLSPubkey

	// Sign takes a message as a slice of bytes and returns a signature as a
	// slice of bytes and an error.
	Sign([]byte) (BLSSignature, error)

	// VerifySignature verifies a signature against a message and a public key.
	VerifySignature(pubKey BLSPubkey, msg []byte, signature BLSSignature) error
}
