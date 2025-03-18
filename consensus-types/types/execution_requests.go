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

const sszDynamicObjectOffset = 4
const maxDepositRequestsPerPayload = 8192
const maxWithdrawalRequestsPerPayload = 16
const maxConsolidationRequestsPerPayload = 2
const sszWithdrawRequestSize = 76          // ExecutionAddress = 20, ValidatorPubKey = 48, Amount = 8
const sszConsolidationRequestSize = 116    // ExecutionAddress = 20, PubKey = 48, Pubkey = 48
const dynamicFieldsInExecutionRequests = 3 // 3 since three dynamic objects (Deposits, Withdrawals, Consolidations)

type ExecutionRequests struct {
	Deposits       []*DepositRequest
	Withdrawals    []*WithdrawalRequest
	Consolidations []*ConsolidationRequest
}

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// WithdrawalRequest is introduced in EIP7002 which we use for withdrawals.
type WithdrawalRequest struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

// ConsolidationRequest is introduced in Pectra but not used by us.
// We keep it so we can maintain parity tests with other SSZ implementations.
type ConsolidationRequest struct {
	SourceAddress common.ExecutionAddress
	SourcePubKey  crypto.BLSPubkey
	TargetPubKey  crypto.BLSPubkey
}

/* -------------------------------------------------------------------------- */
/*                       Execution Requests SSZ                               */
/* -------------------------------------------------------------------------- */

func (e *ExecutionRequests) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Deposits, maxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Withdrawals, maxWithdrawalRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Consolidations, maxConsolidationRequestsPerPayload)

	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Deposits, maxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Withdrawals, maxWithdrawalRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Consolidations, maxConsolidationRequestsPerPayload)
}

func (e *ExecutionRequests) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	size := uint32(sszDynamicObjectOffset * dynamicFieldsInExecutionRequests)
	if fixed {
		return size
	}
	size += ssz.SizeSliceOfStaticObjects(siz, e.Deposits)
	size += ssz.SizeSliceOfStaticObjects(siz, e.Withdrawals)
	size += ssz.SizeSliceOfStaticObjects(siz, e.Consolidations)
	return size
}

func (e *ExecutionRequests) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(e))
	return buf, ssz.EncodeToBytes(buf, e)
}

func (e *ExecutionRequests) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, e)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (e *ExecutionRequests) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(e)
}

/* -------------------------------------------------------------------------- */
/*                       Withdrawal Requests SSZ                              */
/* -------------------------------------------------------------------------- */

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

/* -------------------------------------------------------------------------- */
/*                       Consolidation Requests SSZ                           */
/* -------------------------------------------------------------------------- */

func (c *ConsolidationRequest) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &c.SourceAddress)
	ssz.DefineStaticBytes(codec, &c.SourcePubKey)
	ssz.DefineStaticBytes(codec, &c.TargetPubKey)
}

func (c *ConsolidationRequest) SizeSSZ(_ *ssz.Sizer) uint32 {
	return sszConsolidationRequestSize
}

func (c *ConsolidationRequest) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(c))
	return buf, ssz.EncodeToBytes(buf, c)
}

func (c *ConsolidationRequest) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, c)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (c *ConsolidationRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(c)
}
