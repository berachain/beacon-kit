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

import "github.com/ethereum/go-ethereum/common/hexutil"

// Validator is a struct that represents a validator in the beacon chain.
// Validator represents a participant in the beacon chain consensus mechanism.
// It holds the validator's public key, withdrawal credentials, effective
// balance, and slashing status.
//

//go:generate go run github.com/fjl/gencodec -type Validator -field-override validatorJSONMarshaling -out validator.json.go
type Validator struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey [48]byte `json:"pubkey"           ssz-size:"48"`
	// Credentials are an address that controls the validator.
	Credentials [32]byte `json:"credentials"      ssz-size:"32"`
	// EffectiveBalance is the validator's current effective balance in gwei.
	EffectiveBalance uint64 `json:"effectiveBalance"`
	// Slashed indicates whether the validator has been slashed.
	Slashed bool `json:"slashed"`
}

// JSON type overrides for ExecutionPayloadEnvelope.
type validatorJSONMarshaling struct {
	Pubkey      hexutil.Bytes
	Credentials hexutil.Bytes
}

// String returns a string representation of the Validator.
func (v Validator) String() string {
	return string(v.Pubkey[:])
}
