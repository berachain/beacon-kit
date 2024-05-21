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

package crypto

import "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"

const (
	// CometBLSType is the BLS curve type used in the Comet BFT consensus
	// algorithm.
	CometBLSType = "bls12_381"

	// TODO: Move this, it doesn't really belong here.
	//
	// CometBLSPower is the voting power given to a validator when they
	// are in the active set.
	CometBLSPower = 100
)

//nolint:lll // links.
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
