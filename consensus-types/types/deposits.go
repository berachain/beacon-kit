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
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/karalabe/ssz"
)

// Deposits is a typealias for a list of Deposits.
type Deposits []*Deposit

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (ds Deposits) SizeSSZ(fixed bool) uint32 {
	return uint32(len(ds) * DepositSize)
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
// TODO: get from accessible chainspec field params.
func (ds Deposits) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*Deposit)(&ds), constants.MaxDeposits)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*Deposit)(&ds), constants.MaxDeposits)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*Deposit)(&ds), constants.MaxDeposits)
	})
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (ds Deposits) HashTreeRoot() common.Root {
	// TODO: determine if using HashConcurrent optimizes performance.
	return ssz.HashSequential(ds)
}

// MarshalSSZ marshals the Deposits object to SSZ format by encoding each deposit individually.
func (dr *Deposits) MarshalSSZ() ([]byte, error) {
	var buf []byte
	depositSize := ssz.Size(&Deposit{})
	// Loop through each deposit in the slice.
	for _, deposit := range *dr {
		// Allocate a buffer for the current deposit.
		depositBuf := make([]byte, depositSize)
		// Encode the deposit into the depositBuf.
		if err := ssz.EncodeToBytes(depositBuf, deposit); err != nil {
			return nil, err
		}
		// Append the encoded bytes to the main buffer.
		buf = append(buf, depositBuf...)
	}
	return buf, nil
}

func (dr *Deposits) NewFromSSZ(data []byte) (*Deposits, error) {
	if dr == nil {
		dr = &Deposits{}
	}
	// Get the size of a single deposit (assumed constant).
	depositSize := int(ssz.Size(&Deposit{}))
	// Ensure that the data length is a multiple of depositSize.
	if len(data)%depositSize != 0 {
		return nil, fmt.Errorf("invalid data length: %d is not a multiple of deposit size %d", len(data), depositSize)
	}
	// Reset the slice.
	*dr = make(Deposits, 0, len(data)/depositSize)

	// Loop through the data in chunks of depositSize.
	for i := 0; i < len(data); i += depositSize {
		chunk := data[i : i+depositSize]
		deposit := new(Deposit)
		if err := ssz.DecodeFromBytes(chunk, deposit); err != nil {
			return nil, err
		}
		*dr = append(*dr, deposit)
	}
	return dr, nil
}
