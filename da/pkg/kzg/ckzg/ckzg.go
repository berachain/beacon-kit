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

package ckzg

import (
	"github.com/berachain/beacon-kit/primitives/pkg/encoding/hex"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

// Implementation is the ethereum/c-kzg-4844 implementation.
const Implementation = "ethereum/c-kzg-4844"

// Verifier is a verifier that utilizies the CKZG library.
type Verifier struct{}

// GetImplementation returns the implementation of the verifier.
func (v Verifier) GetImplementation() string {
	return Implementation
}

// NewVerifier creates a new CKZG verifier.
//
//nolint:mnd // lots of random numbers because cryptography.
func NewVerifier(ts *gokzg4844.JSONTrustedSetup) (*Verifier, error) {
	if err := gokzg4844.CheckTrustedSetupIsWellFormed(ts); err != nil {
		return nil, err
	}
	g1s := make(
		[]byte,
		len(ts.SetupG1Lagrange)*(len(ts.SetupG1Lagrange[0])-2)/2,
	)
	for i, g1 := range ts.SetupG1Lagrange {
		copy(g1s[i*(len(g1)-2)/2:], hex.MustToBytes(g1))
	}
	g2s := make([]byte, len(ts.SetupG2)*(len(ts.SetupG2[0])-2)/2)
	for i, g2 := range ts.SetupG2 {
		copy(g2s[i*(len(g2)-2)/2:], hex.MustToBytes(g2))
	}
	if err := ckzg4844.LoadTrustedSetup(g1s, g2s); err != nil {
		return nil, err
	}
	return &Verifier{}, nil
}
