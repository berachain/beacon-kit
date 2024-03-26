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

package signing

import "github.com/berachain/beacon-kit/primitives"

// SigningData is a struct used to compute
// hash(root_hash(object), domain_hash).
// Spec:
// https://github.com/ethereum/annotated-spec/blob/master/phase0/beacon-chain.md#signingdata.
//
//nolint:lll // Urls are long.
type Data struct {
	ObjectRoot primitives.Root `ssz-size:"32"`
	Domain     Domain          `ssz-size:"32"`
}

// ComputeSigningRoot computes the signing root for a given object and domain.
// Spec:
// def compute_signing_root(ssz_object: SSZObject, domain: Domain) -> Root:
//
//	return hash_tree_root(SigningData(
//	    object_root=hash_tree_root(ssz_object),
//	    domain=domain,
//	))
func ComputeSigningRoot(
	sszObject SSZObject,
	domain Domain,
) (primitives.Root, error) {
	objectRoot, err := sszObject.HashTreeRoot()
	if err != nil {
		return primitives.Root{}, err
	}
	data := Data{
		ObjectRoot: objectRoot,
		Domain:     domain,
	}
	return data.HashTreeRoot()
}
