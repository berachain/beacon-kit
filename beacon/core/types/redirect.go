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

import "strconv"

type Redirect struct {
	// Public key of the validator specified in the deposit.
	Pubkey []byte `json:"pubkey"      ssz-max:"48"`
	// Public key of the validator specified in the deposit.
	NewPubkey []byte `json:"newPubkey"   ssz-max:"48"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials DepositCredentials `json:"credentials"              ssz-size:"32"`
	// Deposit amount in gwei.
	Amount uint64 `json:"amount"`
	// Signature of the deposit data.
	Signature []byte `json:"signature"   ssz-max:"96"`
	// Index of the redirect in the deposit contract.
	Index uint64
}

func (r *Redirect) String() string {
	return "Redirect{" +
		"Pubkey: " + string(r.Pubkey) +
		", Credentials: " + string(r.Credentials[:]) +
		", Amount: " + strconv.FormatUint(r.Amount, 10) +
		", Signature: " + string(r.Signature) +
		"}"
}
