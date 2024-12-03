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

package types

import (
	"github.com/berachain/beacon-kit/primitives/pkg/common"
	"github.com/berachain/beacon-kit/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

// ForkData as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#forkdata
//
//nolint:lll
type ForkData struct {
	// CurrentVersion is the current version of the fork.
	CurrentVersion common.Version
	// GenesisValidatorsRoot is the root of the genesis validators.
	GenesisValidatorsRoot common.Root
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

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the SigningData object in SSZ encoding.
func (*ForkData) SizeSSZ(*ssz.Sizer) uint32 {
	//nolint:mnd // 32+4 = 36.
	return 36
}

// DefineSSZ defines the SSZ encoding for the ForkData object.
func (fd *ForkData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &fd.CurrentVersion)
	ssz.DefineStaticBytes(codec, &fd.GenesisValidatorsRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the ForkData object.
func (fd *ForkData) HashTreeRoot() common.Root {
	return ssz.HashSequential(fd)
}

// MarshalSSZTo marshals the ForkData object to SSZ format into the provided
// buffer.
func (fd *ForkData) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, fd)
}

// MarshalSSZ marshals the ForkData object to SSZ format.
func (fd *ForkData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(fd))
	return fd.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the ForkData object from SSZ format.
func (fd *ForkData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, fd)
}

// ComputeDomain as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#compute_domain
//
//nolint:lll
func (fd *ForkData) ComputeDomain(
	domainType common.DomainType,
) common.Domain {
	forkDataRoot := fd.HashTreeRoot()
	return common.Domain(
		append(
			domainType[:],
			forkDataRoot[:28]...),
	)
}

// ComputeRandaoSigningRoot computes the randao signing root.
func (fd *ForkData) ComputeRandaoSigningRoot(
	domainType common.DomainType,
	epoch math.Epoch,
) common.Root {
	return ComputeSigningRootUInt64(
		epoch.Unwrap(),
		fd.ComputeDomain(domainType),
	)
}
