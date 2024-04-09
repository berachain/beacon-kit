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

package forks

import (
	"github.com/berachain/beacon-kit/mod/primitives"
)

// ForkData as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#forkdata
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path data.go -objs ForkData -include ../primitives -output data.ssz.go
//nolint:lll
type ForkData struct {
	// CurrentVersion is the current version of the fork.
	CurrentVersion primitives.Version `ssz-size:"4"`
	// GenesisValidatorsRoot is the root of the genesis validators.
	GenesisValidatorsRoot primitives.Root `ssz-size:"32"`
}

// NewForkData creates a new ForkData struct.
func NewForkData(
	currentVersion primitives.Version, genesisValidatorsRoot primitives.Root,
) *ForkData {
	return &ForkData{
		CurrentVersion:        currentVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// ComputeDomain as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_domain
//
//nolint:lll
func (fv ForkData) ComputeDomain(
	domainType primitives.DomainType,
) (primitives.Domain, error) {
	forkDataRoot, err := fv.HashTreeRoot()
	if err != nil {
		return primitives.Domain{}, err
	}

	return primitives.Domain(
		append(
			domainType[:],
			forkDataRoot[:28]...),
	), nil
}
