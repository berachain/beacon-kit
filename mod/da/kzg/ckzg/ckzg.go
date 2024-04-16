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

package ckzg

import (
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Verifier is a verifier that utilizies the CKZG library.
type Verifier struct{}

// NewVerifier creates a new CKZG verifier.
//
//nolint:gomnd // lots of random numbers because cryptography.
func NewVerifier(ts *gokzg4844.JSONTrustedSetup) (*Verifier, error) {
	if err := gokzg4844.CheckTrustedSetupIsWellFormed(ts); err != nil {
		return nil, err
	}
	g1s := make(
		[]byte,
		len(ts.SetupG1Lagrange)*(len(ts.SetupG1Lagrange[0])-2)/2,
	)
	for i, g1 := range ts.SetupG1Lagrange {
		copy(g1s[i*(len(g1)-2)/2:], hexutil.MustDecode(g1))
	}
	g2s := make([]byte, len(ts.SetupG2)*(len(ts.SetupG2[0])-2)/2)
	for i, g2 := range ts.SetupG2 {
		copy(g2s[i*(len(g2)-2)/2:], hexutil.MustDecode(g2))
	}
	if err := ckzg4844.LoadTrustedSetup(g1s, g2s); err != nil {
		return nil, err
	}
	return &Verifier{}, nil
}
