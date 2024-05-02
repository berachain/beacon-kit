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

package ssz

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
)

// SigningData as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#signingdata
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen -path signing_data.go -objs SigningData -include ../bytes,../common -output signing_data.ssz.go
type SigningData struct {
	ObjectRoot common.Root   `ssz-size:"32"`
	Domain     common.Domain `ssz-size:"32"`
}

// ComputeSigningRoot as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_signing_root
//
//nolint:lll
func ComputeSigningRoot(
	sszObject interface{ HashTreeRoot() ([32]byte, error) },
	domain common.Domain,
) (common.Root, error) {
	objectRoot, err := sszObject.HashTreeRoot()
	if err != nil {
		return common.Root{}, err
	}
	return (&SigningData{
		ObjectRoot: objectRoot,
		Domain:     domain,
	}).HashTreeRoot()
}

// ComputeSigningRootUInt64 computes the signing root of a uint64 value.
func ComputeSigningRootUInt64(
	value uint64,
	domain common.Domain,
) (common.Root, error) {
	bz := make([]byte, constants.RootLength)
	binary.LittleEndian.PutUint64(bz, value)
	return (&SigningData{
		ObjectRoot: common.Root(bz),
		Domain:     domain,
	}).HashTreeRoot()
}
