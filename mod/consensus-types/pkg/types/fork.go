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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// ForkSize is the size of the Fork object in bytes.
// 4 bytes for PreviousVersion + 4 bytes for CurrentVersion + 8 bytes for Epoch.
const ForkSize = 16

// Fork as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#fork
//
//nolint:lll
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
func (f *Fork) New(
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

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size of the Fork object in bytes.
func (f *Fork) SizeSSZ() uint32 {
	return ForkSize
}

// DefineSSZ defines the SSZ encoding for the Fork object.
func (f *Fork) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &f.PreviousVersion)
	ssz.DefineStaticBytes(codec, &f.CurrentVersion)
	ssz.DefineUint64(codec, &f.Epoch)
}

// MarshalSSZ marshals the Fork object to SSZ format.
func (f *Fork) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, f.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, f)
}

// UnmarshalSSZ unmarshals the Fork object from SSZ format.
func (f *Fork) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, f)
}

// HashTreeRoot computes the SSZ hash tree root of the Fork object.
func (f *Fork) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(f), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the Fork object to a target array.
func (f *Fork) MarshalSSZTo(buf []byte) ([]byte, error) {
	dst := buf

	// Field (0) 'PreviousVersion'
	dst = append(dst, f.PreviousVersion[:]...)

	// Field (1) 'CurrentVersion'
	dst = append(dst, f.CurrentVersion[:]...)

	// Field (2) 'Epoch'
	dst = fastssz.MarshalUint64(dst, uint64(f.Epoch))

	return dst, nil
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
