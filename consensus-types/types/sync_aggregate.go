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

// TODO: Enable once full migration to fastssz is complete
// go:generate sszgen --path . --include ../../primitives/common,../../primitives/bytes,../../primitives/crypto --objs SyncAggregate --output sync_aggregate_sszgen.go

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	fastssz "github.com/ferranbt/fastssz"
)

// TODO: Re-enable interface assertions once constraints are updated
// var (
// 	_ common.UnusedEnforcer               = (*SyncAggregate)(nil)
// )

const (
	syncCommitteeSize       = 512
	syncCommitteeBitsLength = syncCommitteeSize / 8
)

type SyncAggregate struct {
	SyncCommitteeBits      [64]byte            `ssz-size:"64"`
	SyncCommitteeSignature crypto.BLSSignature `ssz-size:"96"`
}

// SizeSSZ returns the SSZ encoded size in bytes for the SyncAggregate.
func (s *SyncAggregate) SizeSSZ() int {
	return 160 // syncCommitteeBitsLength + 96
}


func (s *SyncAggregate) ValidateAfterDecodingSSZ() error { return s.EnforceUnused() }

// MarshalSSZ marshals the SyncAggregate into SSZ format.
func (s *SyncAggregate) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, 160)
	return s.MarshalSSZTo(buf)
}

// HashTreeRoot returns the SSZ hash tree root of the SyncAggregate.
func (s *SyncAggregate) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	s.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

// EnforceUnused return true if the SyncAggregate contains all zero values.
// As long as this type remains unused and unvalidated by consensus,
// we must enforce that it contains no data.
func (s *SyncAggregate) EnforceUnused() error {
	if (s != nil && *s != SyncAggregate{}) {
		return errors.New("SyncAggregate must be unused")
	}
	return nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the SyncAggregate object to a target array.
func (s *SyncAggregate) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Field (0) 'SyncCommitteeBits'
	dst = append(dst, s.SyncCommitteeBits[:]...)

	// Field (1) 'SyncCommitteeSignature'
	dst = append(dst, s.SyncCommitteeSignature[:]...)

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the SyncAggregate object.
func (s *SyncAggregate) UnmarshalSSZ(buf []byte) error {
	if len(buf) != 160 {
		return errors.Wrapf(fastssz.ErrSize, "expected buffer of length 160, received %d", len(buf))
	}
	copy(s.SyncCommitteeBits[:], buf[0:64])
	copy(s.SyncCommitteeSignature[:], buf[64:160])
	return s.ValidateAfterDecodingSSZ()
}


// HashTreeRootWith ssz hashes the SyncAggregate object with a hasher.
func (s *SyncAggregate) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'SyncCommitteeBits'
	hh.PutBytes(s.SyncCommitteeBits[:])

	// Field (1) 'SyncCommitteeSignature'
	hh.PutBytes(s.SyncCommitteeSignature[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the SyncAggregate object.
func (s *SyncAggregate) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(s)
}
