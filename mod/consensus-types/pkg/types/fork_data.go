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

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// ForkData as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#forkdata
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path fork_data.go -objs ForkData -include ../../../primitives/pkg/bytes,../../../primitives/pkg/common -output fork_data.ssz.go
//nolint:lll
type ForkData struct {
	// CurrentVersion is the current version of the fork.
	CurrentVersion common.Version `ssz-size:"4"`
	// GenesisValidatorsRoot is the root of the genesis validators.
	GenesisValidatorsRoot common.Root `ssz-size:"32"`
}

// NewForkData creates a new ForkData struct.
func NewForkData(
	currentVersion common.Version, genesisValidatorsRoot common.Root,
) *ForkData {
	return &ForkData{
		CurrentVersion:        currentVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// New creates a new ForkData struct.
func (fd *ForkData) New(
	currentVersion common.Version, genesisValidatorsRoot common.Root,
) *ForkData {
	return NewForkData(currentVersion, genesisValidatorsRoot)
}

// ComputeDomain as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_domain
//
//nolint:lll
func (fd *ForkData) ComputeDomain(
	domainType common.DomainType,
) (common.Domain, error) {
	forkDataRoot, err := fd.HashTreeRoot()
	if err != nil {
		return common.Domain{}, err
	}

	return common.Domain(
		append(
			domainType[:],
			forkDataRoot[:28]...),
	), nil
}

// ComputeRandaoSigningRoot computes the randao signing root.
func (fd *ForkData) ComputeRandaoSigningRoot(
	domainType common.DomainType,
	epoch math.Epoch,
) (common.Root, error) {
	signingDomain, err := fd.ComputeDomain(domainType)
	if err != nil {
		return primitives.Root{}, err
	}

	signingRoot, err := ssz.ComputeSigningRootUInt64(
		uint64(epoch),
		signingDomain,
	)

	if err != nil {
		return primitives.Root{},
			errors.Newf("failed to compute signing root: %w", err)
	}
	return signingRoot, nil
}
