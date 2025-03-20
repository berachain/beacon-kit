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

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/karalabe/ssz"
)

const sszDynamicObjectOffset = 4
const maxWithdrawalRequestsPerPayload = 16
const maxConsolidationRequestsPerPayload = 2
const sszWithdrawRequestSize = 76          // ExecutionAddress = 20, ValidatorPubKey = 48, Amount = 8
const sszConsolidationRequestSize = 116    // ExecutionAddress = 20, PubKey = 48, Pubkey = 48
const dynamicFieldsInExecutionRequests = 3 // 3 since three dynamic objects (Deposits, Withdrawals, Consolidations)

// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#execution-layer-triggered-requests

type ExecutionRequests struct {
	Deposits       []*DepositRequest
	Withdrawals    []*WithdrawalRequest
	Consolidations []*ConsolidationRequest
}

// GetExecutionRequestsList introduced in pectra from the consensus spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_execution_requests_list
func GetExecutionRequestsList(er *ExecutionRequests) ([][]byte, error) {
	result := make([][]byte, 0)

	// Process deposit requests if non-empty.
	if len(er.Deposits) > 0 {
		requests := DepositRequests(er.Deposits)
		depositBytes, err := requests.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		combined := append([]byte{}, depositRequestType()...)
		combined = append(combined, depositBytes...)
		result = append(result, combined)
	}

	// Process withdrawal requests if non-empty.
	if len(er.Withdrawals) > 0 {
		requests := WithdrawalRequests(er.Withdrawals)
		withdrawalBytes, err := requests.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		combined := append([]byte{}, withdrawalRequestType()...)
		combined = append(combined, withdrawalBytes...)
		result = append(result, combined)
	}

	// Process consolidation requests if non-empty.
	if len(er.Consolidations) > 0 {
		requests := ConsolidationRequests(er.Consolidations)
		consolidationBytes, err := requests.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		combined := append([]byte{}, consolidationRequestType()...)
		combined = append(combined, consolidationBytes...)
		result = append(result, combined)
	}

	return result, nil
}

// DecodeExecutionRequests is used to decode the result from GetPayload into an ExecutionRequests.
func DecodeExecutionRequests(encodedRequests [][]byte) (*ExecutionRequests, error) {
	var result ExecutionRequests
	var prevType *uint8

	// Iterate over each encoded request group.
	for _, encoded := range encodedRequests {
		if len(encoded) < 1 {
			return nil, errors.New("invalid execution request, length less than 1")
		}
		// The first byte indicates the request type.
		reqType := encoded[0]

		// Enforce that request types are in strictly increasing order.
		if prevType != nil && *prevType >= reqType {
			return nil, errors.New("requests should be in sorted order and unique")
		}
		prevType = &reqType

		// The remaining bytes are the SSZ serialization for this group.
		data := encoded[1:]

		// Switch based on the request type.
		switch reqType {
		case depositRequestType()[0]:
			var req *DepositRequests
			req, err := req.NewFromSSZ(data)
			if err != nil {
				return nil, err
			}
			result.Deposits = *req
		case withdrawalRequestType()[0]:
			var req *WithdrawalRequests
			req, err := req.NewFromSSZ(data)
			if err != nil {
				return nil, err
			}
			result.Withdrawals = *req
		case consolidationRequestType()[0]:
			var req *ConsolidationRequests
			req, err := req.NewFromSSZ(data)
			if err != nil {
				return nil, err
			}
			result.Consolidations = *req
		default:
			return nil, fmt.Errorf("unsupported request type %d", reqType)
		}
	}

	return &result, nil
}

func depositRequestType() []byte {
	return []byte{0x00}
}

func withdrawalRequestType() []byte {
	return []byte{0x01}
}

func consolidationRequestType() []byte {
	return []byte{0x02}
}

/* -------------------------------------------------------------------------- */
/*                       Execution Requests SSZ                               */
/* -------------------------------------------------------------------------- */

func (e *ExecutionRequests) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Deposits, MaxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Withdrawals, maxWithdrawalRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Consolidations, maxConsolidationRequestsPerPayload)

	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Deposits, MaxDepositRequestsPerPayload)
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
