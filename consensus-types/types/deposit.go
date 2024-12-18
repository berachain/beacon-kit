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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/karalabe/ssz"
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer. The proof is attached by the consensus layer.
type Deposit struct {
	// Proof is the proof of the deposit data.
	Proof [constants.DepositContractDepth + 1]common.Root

	// Data is the underlying deposit data.
	Data *DepositData
}

// NewDeposit creates a new Deposit instance.
func NewDeposit(
	proof [constants.DepositContractDepth + 1]common.Root, depositData *DepositData,
) *Deposit {
	return &Deposit{Proof: proof, Data: depositData}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// DefineSSZ defines the SSZ encoding for the Deposit object.
func (d *Deposit) DefineSSZ(c *ssz.Codec) {
	ssz.DefineArrayOfStaticBytes[
		[constants.DepositContractDepth + 1]common.Root, common.Root,
	](c, &d.Proof)
	ssz.DefineStaticObject(c, &d.Data)
}

// MarshalSSZ marshals the Deposit object to SSZ format.
func (d *Deposit) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(d))
	return buf, ssz.EncodeToBytes(buf, d)
}

// UnmarshalSSZ unmarshals the Deposit object from SSZ format.
func (d *Deposit) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, d)
}

// SizeSSZ returns the SSZ encoded size of the Deposit object.
func (d *Deposit) SizeSSZ(*ssz.Sizer) uint32 {
	return 1248 //nolint:mnd // (32 + 1) * 32 + 192 = 1248.
}

// HashTreeRoot computes the Merkleization of the Deposit object.
func (d *Deposit) HashTreeRoot() common.Root {
	return ssz.HashSequential(d)
}
