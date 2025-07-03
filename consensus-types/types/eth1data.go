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

// TODO: Eth1Data needs manual fastssz migration to handle dual interface compatibility
// go:generate sszgen -path . -objs Eth1Data -output eth1data_sszgen.go -include ../../primitives/common,../../primitives/math

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
)

// Eth1DataSize is the size of the Eth1Data object in bytes.
// 32 bytes for DepositRoot + 8 bytes for DepositCount + 8 bytes for BlockHash.
const Eth1DataSize = 72

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZMarshallableRootable = (*Eth1Data)(nil)
// )

type Eth1Data struct {
	// DepositRoot is the root of the deposit tree.
	DepositRoot common.Root `json:"depositRoot"`
	// DepositCount is the number of deposits in the deposit tree.
	DepositCount math.U64 `json:"depositCount"`
	// BlockHash is the hash of the block corresponding to the Eth1Data.
	BlockHash common.ExecutionHash `json:"blockHash"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

func NewEth1Data(depositRoot common.Root) *Eth1Data {
	return &Eth1Data{
		DepositRoot: depositRoot,
	}
}

func NewEmptyEth1Data() *Eth1Data {
	return &Eth1Data{}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the Eth1Data object in SSZ encoding.
func (*Eth1Data) SizeSSZ() int {
	return Eth1DataSize
}


// HashTreeRoot computes the SSZ hash tree root of the Eth1Data object.
func (e *Eth1Data) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	e.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

// MarshalSSZ marshals the Eth1Data object to SSZ format.
func (e *Eth1Data) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, Eth1DataSize)
	return e.MarshalSSZTo(buf)
}

func (*Eth1Data) ValidateAfterDecodingSSZ() error { return nil }

// MarshalSSZTo marshals the Eth1Data object into a pre-allocated byte slice.
func (e *Eth1Data) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Field (0) 'DepositRoot'
	dst = append(dst, e.DepositRoot[:]...)

	// Field (1) 'DepositCount'
	dst = fastssz.MarshalUint64(dst, uint64(e.DepositCount))

	// Field (2) 'BlockHash'
	dst = append(dst, e.BlockHash[:]...)

	return dst, nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the Eth1Data object with a hasher.
func (e *Eth1Data) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'DepositRoot'
	hh.PutBytes(e.DepositRoot[:])

	// Field (1) 'DepositCount'
	hh.PutUint64(uint64(e.DepositCount))

	// Field (2) 'BlockHash'
	hh.PutBytes(e.BlockHash[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Eth1Data object.
func (e *Eth1Data) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(e)
}

// UnmarshalSSZ ssz unmarshals the Eth1Data object.
func (e *Eth1Data) UnmarshalSSZ(buf []byte) error {
	if len(buf) != Eth1DataSize {
		return fastssz.ErrSize
	}

	// Field (0) 'DepositRoot'
	copy(e.DepositRoot[:], buf[0:32])

	// Field (1) 'DepositCount'
	e.DepositCount = math.U64(fastssz.UnmarshallUint64(buf[32:40]))

	// Field (2) 'BlockHash'
	copy(e.BlockHash[:], buf[40:72])

	return nil
}


// GetDepositCount returns the deposit count.
func (e *Eth1Data) GetDepositCount() math.U64 {
	return e.DepositCount
}
