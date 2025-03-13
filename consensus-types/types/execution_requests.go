// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Berachain Foundation
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
// WdeHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
)

const maxDepositRequestsPerPayload = 8192
const maxWithdrawalRequestsPerPayload = 16
const sszDynamicObjectOffset = 4
const sszDepositRequestSize = 192          // Pubkey = 48, WithdrawalCredentials = 32, Amount = 8, Signature = 96, Index = 8.
const sszWithdrawRequestSize = 76          // ExecutionAddress = 20, ValidatorPubKey = 48, Amount = 8
const dynamicFieldsInExecutionRequests = 2 // 2 since two dynamic objects (Deposits, Withdrawals)

type ExecutionRequests struct {
	Deposits    []*DepositRequest
	Withdrawals []*WithdrawalRequest
}

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest struct {
	Pubkey                crypto.BLSPubkey
	WithdrawalCredentials WithdrawalCredentials
	Amount                math.Gwei
	Signature             crypto.BLSSignature
	Index                 math.U64
}

// WithdrawalRequest is introduced in EIP7002 which we use for withdrawals.
type WithdrawalRequest struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

func (e *ExecutionRequests) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Deposits, maxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Withdrawals, maxWithdrawalRequestsPerPayload)

	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Deposits, maxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Withdrawals, maxWithdrawalRequestsPerPayload)
}

func (e *ExecutionRequests) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	// Multiply by 2 since two dynamic objects (Deposits, Withdrawals)
	size := uint32(sszDynamicObjectOffset * dynamicFieldsInExecutionRequests)
	if fixed {
		return size
	}
	size += ssz.SizeSliceOfStaticObjects(siz, e.Deposits)
	size += ssz.SizeSliceOfStaticObjects(siz, e.Withdrawals)
	return size
}

func (d *DepositRequest) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &d.Pubkey)
	ssz.DefineStaticBytes(codec, &d.WithdrawalCredentials)
	ssz.DefineUint64(codec, &d.Amount)
	ssz.DefineStaticBytes(codec, &d.Signature)
	ssz.DefineUint64(codec, &d.Index)
}

func (d *DepositRequest) SizeSSZ(_ *ssz.Sizer) uint32 {
	return sszDepositRequestSize
}

func (d *DepositRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(d))
	return buf, ssz.EncodeToBytes(buf, d)
}

func (d *DepositRequest) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, d)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (d *DepositRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(d)
}

func (w *WithdrawalRequest) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &w.SourceAddress)
	ssz.DefineStaticBytes(codec, &w.ValidatorPubKey)
	ssz.DefineUint64(codec, &w.Amount)
}

func (w *WithdrawalRequest) SizeSSZ(_ *ssz.Sizer) uint32 {
	return sszWithdrawRequestSize
}

func (w *WithdrawalRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(w))
	return buf, ssz.EncodeToBytes(buf, w)
}

func (w *WithdrawalRequest) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, w)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (w *WithdrawalRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(w)
}
