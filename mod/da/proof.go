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

package da

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/da/proof"
	"github.com/berachain/beacon-kit/mod/da/proof/ckzg"
	"github.com/berachain/beacon-kit/mod/da/proof/gokzg"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// NewBlobProofVerifier creates a new BlobVerifier with the given
// implementation.
func NewBlobProofVerifier(
	impl string,
	ts *gokzg4844.JSONTrustedSetup,
) (proof.BlobProofVerifier, error) {
	switch impl {
	case "crate-crypto/go-kzg-4844":
		return gokzg.NewVerifier(ts)
	case "ethereum/c-kzg-4844":
		return ckzg.NewVerifier(ts)
	default:
		return nil, errors.New("unsupported KZG implementation")
	}
}
