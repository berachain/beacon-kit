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

package randao

import (
	"cosmossdk.io/log"
	crypto "github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
)

type Option func(*Processor) error

// WithBeaconStateProvider sets the beacon state provider.
func WithBeaconStateProvider(beaconStateProvider BeaconStateProvider) Option {
	return func(p *Processor) error {
		p.BeaconStateProvider = beaconStateProvider
		return nil
	}
}

// WithSigner sets the signer.
func WithSigner(
	signer crypto.Signer[[bls12381.SignatureLength]byte],
) Option {
	return func(p *Processor) error {
		p.signer = signer
		return nil
	}
}

// WithLogger sets the logger.
func WithLogger(logger log.Logger) Option {
	return func(p *Processor) error {
		p.logger = logger
		return nil
	}
}
