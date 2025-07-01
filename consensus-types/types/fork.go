// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// ForkSize is the size of the Fork object in bytes.
// 4 bytes for PreviousVersion + 4 bytes for CurrentVersion + 8 bytes for Epoch.
const ForkSize = 16

var (
	_ constraints.SSZMarshallableRootable = (*Fork)(nil)
)

// Fork as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#fork
type Fork struct {
	// PreviousVersion is the last version before the fork.
	PreviousVersion common.Version `json:"previous_version"`
	// CurrentVersion is the first version after the fork.
	CurrentVersion common.Version `json:"current_version"`
	// Epoch is the epoch at which the fork occurred.
	Epoch math.Epoch `json:"epoch"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

// New creates a new fork.
func NewFork(
	previousVersion common.Version,
	currentVersion common.Version,
	epoch math.Epoch,
) *Fork {
	return &Fork{
		PreviousVersion: previousVersion,
		CurrentVersion:  currentVersion,
		Epoch:           epoch,
	}
}

func NewEmptyFork() *Fork {
	return &Fork{}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// MarshalSSZ marshals the Fork object to SSZ format.
func (f *Fork) MarshalSSZ() ([]byte, error) {
	return f.MarshalSSZTo(nil)
}

// UnmarshalSSZ unmarshals the Fork object from SSZ format.
func (f *Fork) UnmarshalSSZ(buf []byte) error {
	if len(buf) < ForkSize {
		return fastssz.ErrSize
	}

	// Field (0) 'PreviousVersion'
	copy(f.PreviousVersion[:], buf[0:4])

	// Field (1) 'CurrentVersion'
	copy(f.CurrentVersion[:], buf[4:8])

	// Field (2) 'Epoch'
	f.Epoch = math.Epoch(fastssz.UnmarshallUint64(buf[8:16]))

	return nil
}

// SizeSSZ returns the SSZ encoded size of the Fork object in bytes.
func (f *Fork) SizeSSZ(*ssz.Sizer) uint32 {
	return ForkSize
}

// DefineSSZ defines the SSZ encoding (required for SSZMarshallable interface).
// Provides backward compatibility for code that still uses karalabe/ssz for unmarshaling.
func (f *Fork) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &f.PreviousVersion)
	ssz.DefineStaticBytes(codec, &f.CurrentVersion)
	ssz.DefineUint64(codec, &f.Epoch)
}

func (*Fork) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the Fork object.
func (f *Fork) HashTreeRoot() common.Root {
	hh := fastssz.NewHasher()
	if err := f.HashTreeRootWith(hh); err != nil {
		panic(err)
	}
	return common.Root(hh.Hash())
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the Fork object to a target array.
func (f *Fork) MarshalSSZTo(buf []byte) ([]byte, error) {
	// Field (0) 'PreviousVersion'
	buf = append(buf, f.PreviousVersion[:]...)

	// Field (1) 'CurrentVersion'
	buf = append(buf, f.CurrentVersion[:]...)

	// Field (2) 'Epoch'
	buf = fastssz.MarshalUint64(buf, uint64(f.Epoch))

	return buf, nil
}

// HashTreeRootWith ssz hashes the Fork object with a hasher.
func (f *Fork) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'PreviousVersion'
	hh.PutBytes(f.PreviousVersion[:])

	// Field (1) 'CurrentVersion'
	hh.PutBytes(f.CurrentVersion[:])

	// Field (2) 'Epoch'
	hh.PutUint64(uint64(f.Epoch))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Fork object.
func (f *Fork) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(f)
}
