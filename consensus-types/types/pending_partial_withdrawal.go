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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

// TODO: PendingPartialWithdrawal ready for full fastssz migration (already has all needed methods)
// go:generate sszgen -path . -objs PendingPartialWithdrawal -output pending_partial_withdrawal_sszgen.go -include ../../primitives/common,../../primitives/math

package types

import (
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// sszPendingPartialWithdrawalSize defines the total SSZ serialized size for
// PendingPartialWithdrawal. The fields are assumed to be encoded as follows:
// - ValidatorIndex: 8 bytes (uint64)
// - Amount:         8 bytes (math.Gwei)
// - WithdrawableEpoch: 8 bytes (uint64)
// Total = 8 + 8 + 8 = 24 bytes.
const sszPendingPartialWithdrawalSize = 24

// Compile-time check to ensure PendingPartialWithdrawal and PendingPartialWithdrawals implements the necessary interfaces.
var (
	_ ssz.StaticObject            = (*PendingPartialWithdrawal)(nil)
	_ constraints.SSZMarshallable = (*PendingPartialWithdrawal)(nil)

	_ ssz.DynamicObject           = (*PendingPartialWithdrawals)(nil)
	_ constraints.SSZMarshallable = (*PendingPartialWithdrawals)(nil)
)

// PendingPartialWithdrawal reflects the following spec:
//
//	class PendingPartialWithdrawal(Container):
//	    validator_index: ValidatorIndex
//	    amount: Gwei
//	    withdrawable_epoch: Epoch
type PendingPartialWithdrawal struct {
	ValidatorIndex    math.ValidatorIndex
	Amount            math.Gwei
	WithdrawableEpoch math.Epoch
}

/* -------------------------------------------------------------------------- */
/*                      PendingPartialWithdrawal SSZ                          */
/* -------------------------------------------------------------------------- */

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawal object
// after decoding from SSZ. Customize further validation as needed.
func (p *PendingPartialWithdrawal) ValidateAfterDecodingSSZ() error {
	return nil
}

// DefineSSZ registers the SSZ encoding for each field in PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &p.ValidatorIndex)
	ssz.DefineUint64(codec, &p.Amount)
	ssz.DefineUint64(codec, &p.WithdrawableEpoch)
}

// SizeSSZ returns the fixed size of the SSZ serialization for PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) SizeSSZ(_ *ssz.Sizer) uint32 {
	return sszPendingPartialWithdrawalSize
}

// MarshalSSZ returns the SSZ encoding of the PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(p))
	return buf, ssz.EncodeToBytes(buf, p)
}

// HashTreeRoot computes and returns the hash tree root for the PendingPartialWithdrawal.
func (p *PendingPartialWithdrawal) HashTreeRoot() common.Root {
	return ssz.HashSequential(p)
}

// HashTreeRootWith SSZ hashes the Deposit object with a hasher. Needed for BeaconState SSZ.
func (p *PendingPartialWithdrawal) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'ValidatorIndex'
	hh.PutUint64(uint64(p.ValidatorIndex))

	// Field (1) 'Amount'
	hh.PutUint64(uint64(p.Amount))

	// Field (2) 'WithdrawableEpoch'
	hh.PutUint64(uint64(p.WithdrawableEpoch))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the PendingPartialWithdrawal object.
func (p *PendingPartialWithdrawal) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(p)
}

// MarshalSSZTo marshals the PendingPartialWithdrawal object into a pre-allocated byte slice.
func (p *PendingPartialWithdrawal) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := p.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// PendingPartialWithdrawals is a SSZ list of PendingPartialWithdrawal containers.
type PendingPartialWithdrawals []*PendingPartialWithdrawal

// NewEmptyPendingPartialWithdrawals returns a new empty PendingPartialWithdrawals list.
func NewEmptyPendingPartialWithdrawals() *PendingPartialWithdrawals {
	return &PendingPartialWithdrawals{}
}

// DefineSSZ defines the SSZ encoding for the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, (*[]*PendingPartialWithdrawal)(p), constants.PendingPartialWithdrawalsLimit)
	ssz.DefineSliceOfStaticObjectsContent(codec, (*[]*PendingPartialWithdrawal)(p), constants.PendingPartialWithdrawalsLimit)
}

// SizeSSZ returns the size of the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	if fixed {
		return constants.SSZOffsetSize
	}
	return constants.SSZOffsetSize + ssz.SizeSliceOfStaticObjects(siz, *p)
}

// MarshalSSZ returns the SSZ encoding of the PendingPartialWithdrawals list.
func (p *PendingPartialWithdrawals) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(p))
	return buf, ssz.EncodeToBytes(buf, p)
}

// ValidateAfterDecodingSSZ validates the PendingPartialWithdrawals list after decoding from SSZ.
func (p *PendingPartialWithdrawals) ValidateAfterDecodingSSZ() error {
	if p == nil {
		return errors.New("nil PendingPartialWithdrawals")
	}
	if len(*p) > constants.PendingPartialWithdrawalsLimit {
		return errors.New("pending partial withdrawals too large")
	}
	return nil
}

// PendingBalanceToWithdraw implements get_pending_balance_to_withdraw from the ETH2.0 spec.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_pending_balance_to_withdraw
func (p PendingPartialWithdrawals) PendingBalanceToWithdraw(validatorIndex math.ValidatorIndex) math.Gwei {
	var total math.Gwei
	for _, withdrawal := range p {
		if withdrawal.ValidatorIndex == validatorIndex {
			total += withdrawal.Amount
		}
	}
	return total
}
