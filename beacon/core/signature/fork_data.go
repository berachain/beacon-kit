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

package signature

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/primitives"
)

// Fork as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#fork
//
//nolint:lll
type Fork struct {
	// PreviousVersion is the last version before the fork.
	PreviousVersion primitives.Version `ssz-size:"32"`
	// CurrentVersion is the first version after the fork.
	CurrentVersion primitives.Version `ssz-size:"32"`
	// Epoch is the epoch at which the fork occurred.
	Epoch primitives.Epoch
}

// ForkData as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#forkdata
//
//nolint:lll
type ForkData struct {
	// CurrentVersion is the current version of the fork.
	CurrentVersion primitives.Version `ssz-size:"4"`
	// GenesisValidatorsRoot is the root of the genesis validators.
	GenesisValidatorsRoot primitives.Root `ssz-size:"32"`
}

// VersionFromUint returns a primitives.Version from a uint32.
func VersionFromUint32(version uint32) primitives.Version {
	versionBz := primitives.Version{}
	binary.LittleEndian.PutUint32(versionBz[:], version)
	return versionBz
}

// ComputeForkDataRoot as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_fork_data_root
//
//nolint:lll
func ComputeForkDataRoot(
	currentVersion primitives.Version,
	genesisValidatorsRoot primitives.Root,
) (primitives.Root, error) {
	forkData := ForkData{
		CurrentVersion:        currentVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
	return forkData.HashTreeRoot()
}

// ComputeForkDigest as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_fork_digest
//
//nolint:lll
func ComputeForkDigest(
	currentVersion primitives.Version,
	genesisValidatorsRoot primitives.Root,
) (primitives.ForkDigest, error) {
	forkDataRoot, err := ComputeForkDataRoot(
		currentVersion,
		genesisValidatorsRoot,
	)
	if err != nil {
		return primitives.ForkDigest{}, err
	}
	return primitives.ForkDigest(forkDataRoot[:4]), nil
}
