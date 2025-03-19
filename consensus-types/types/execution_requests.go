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
	"fmt"

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

// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#execution-layer-triggered-requests
var (
	DepositRequestType       = []byte{0x00}
	WithdrawalRequestType    = []byte{0x01}
	ConsolidationRequestType = []byte{0x02}
)

type ExecutionRequests struct {
	Deposits       []*DepositRequest
	Withdrawals    []*WithdrawalRequest
	Consolidations []*ConsolidationRequest
}

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// DepositRequests is used for SSZ unmarshalling a list of DepositRequest
type DepositRequests []*DepositRequest

// WithdrawalRequest is introduced in EIP7002 which we use for withdrawals.
type WithdrawalRequest struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

// WithdrawalRequests is used for SSZ unmarshalling a list of WithdrawalRequest
type WithdrawalRequests []*WithdrawalRequest

// ConsolidationRequest is introduced in Pectra but not used by us.
// We keep it so we can maintain parity tests with other SSZ implementations.
type ConsolidationRequest struct {
	SourceAddress common.ExecutionAddress
	SourcePubKey  crypto.BLSPubkey
	TargetPubKey  crypto.BLSPubkey
}

// ConsolidationRequests is used for SSZ unmarshalling a list of ConsolidationRequest
type ConsolidationRequests []*ConsolidationRequest

// GetExecutionRequestsList introduced in pectra from the consensus spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_execution_requests_list
// TODO(pectra): Test this
func GetExecutionRequestsList(er *ExecutionRequests) ([][]byte, error) {
	var result [][]byte

	// Process deposit requests if non-empty.
	if len(er.Deposits) > 0 {
		depositBytes, err := marshalSSZDeposits(er.Deposits)
		if err != nil {
			return nil, err
		}
		combined := append(DepositRequestType, depositBytes...)
		result = append(result, combined)
	}

	// Process withdrawal requests if non-empty.
	if len(er.Withdrawals) > 0 {
		withdrawalBytes, err := marshalSSZWithdrawals(er.Withdrawals)
		if err != nil {
			return nil, err
		}
		combined := append(WithdrawalRequestType, withdrawalBytes...)
		result = append(result, combined)
	}

	// Process consolidation requests if non-empty.
	if len(er.Consolidations) > 0 {
		consolidationBytes, err := marshalSSZConsolidations(er.Consolidations)
		if err != nil {
			return nil, err
		}
		combined := append(ConsolidationRequestType, consolidationBytes...)
		result = append(result, combined)
	}

	return result, nil
}

// DecodeExecutionRequests is used to decode the result from GetPayload into an ExecutionRequests.
func DecodeExecutionRequests(encodedRequests [][]byte) (*ExecutionRequests, error) {
	var result ExecutionRequests
	var prevType *uint8

	// Iterate over each encoded request group.
	for i, encoded := range encodedRequests {
		if len(encoded) < 1 {
			return nil, fmt.Errorf("encoded request group %d is empty", i)
		}
		// The first byte indicates the request type.
		reqType := encoded[0]

		// Enforce that request types are in strictly increasing order.
		if prevType != nil && *prevType >= reqType {
			return nil, fmt.Errorf("invalid request type order or duplicate at group %d: got %d after %d", i, reqType, *prevType)
		}
		prevType = &reqType

		// The remaining bytes are the SSZ serialization for this group.
		data := encoded[1:]

		// Switch based on the request type.
		switch reqType {
		case DepositRequestType[0]:
			drs, err := unmarshalSSZDeposits(data)
			if err != nil {
				return nil, err
			}
			result.Deposits = drs
		case WithdrawalRequestType[0]:
			wrs, err := unmarshalSSZWithdrawals(data)
			if err != nil {
				return nil, err
			}
			result.Withdrawals = wrs
		case ConsolidationRequestType[0]:
			crs, err := unmarshalSSZConsolidations(data)
			if err != nil {
				return nil, err
			}
			result.Consolidations = crs
		default:
			return nil, fmt.Errorf("unsupported request type %d", reqType)
		}
	}

	return &result, nil
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

func (e *ExecutionRequests) NewFromSSZ(buf []byte) (*ExecutionRequests, error) {
	if e == nil {
		e = &ExecutionRequests{}
	}
	return e, ssz.DecodeFromBytes(buf, e)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (e *ExecutionRequests) HashTreeRoot() common.Root {
	return ssz.HashSequential(e)
}

/* -------------------------------------------------------------------------- */
/*                       Deposit    Requests SSZ                              */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (dr DepositRequests) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*DepositRequest)(dr))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
func (dr DepositRequests) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*DepositRequest)(&dr), maxDepositRequestsPerPayload)
	})
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

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (wr WithdrawalRequests) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*WithdrawalRequest)(wr))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
func (wr WithdrawalRequests) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*WithdrawalRequest)(&wr), maxWithdrawalRequestsPerPayload)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*WithdrawalRequest)(&wr), maxWithdrawalRequestsPerPayload)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*WithdrawalRequest)(&wr), maxWithdrawalRequestsPerPayload)
	})
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (wr WithdrawalRequests) HashTreeRoot() common.Root {
	return ssz.HashSequential(wr)
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

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (cr ConsolidationRequests) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(siz, ([]*ConsolidationRequest)(cr))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
// TODO: get from accessible chainspec field params.
func (cr ConsolidationRequests) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(c, (*[]*ConsolidationRequest)(&cr), maxConsolidationRequestsPerPayload)
	})
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (cr ConsolidationRequests) HashTreeRoot() common.Root {
	return ssz.HashSequential(cr)
}

/* -------------------------------------------------------------------------- */
/*                       SSZ Utilities                                        */
/* -------------------------------------------------------------------------- */

func marshalSSZDeposits(deposits []*DepositRequest) ([]byte, error) {
	d := DepositRequests(deposits)
	buf := make([]byte, ssz.Size(&d))
	err := ssz.EncodeToBytes(buf, &d)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func unmarshalSSZDeposits(data []byte) ([]*DepositRequest, error) {
	var deps DepositRequests
	err := ssz.DecodeFromBytes(data, &deps)
	return deps, err
}

func marshalSSZWithdrawals(withdrawals []*WithdrawalRequest) ([]byte, error) {
	w := WithdrawalRequests(withdrawals)
	buf := make([]byte, ssz.Size(w))
	err := ssz.EncodeToBytes(buf, w)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func unmarshalSSZWithdrawals(data []byte) ([]*WithdrawalRequest, error) {
	var withdrawals WithdrawalRequests
	err := ssz.DecodeFromBytes(data, &withdrawals)
	return withdrawals, err
}

func marshalSSZConsolidations(consolidations []*ConsolidationRequest) ([]byte, error) {
	c := ConsolidationRequests(consolidations)
	buf := make([]byte, ssz.Size(c))
	err := ssz.EncodeToBytes(buf, c)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func unmarshalSSZConsolidations(data []byte) ([]*ConsolidationRequest, error) {
	var consolidations ConsolidationRequests
	err := ssz.DecodeFromBytes(data, &consolidations)
	return consolidations, err
}
