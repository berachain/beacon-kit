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

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	sszutil "github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/karalabe/ssz"
)

// 3 since three dynamic objects (Deposits, Withdrawals, Consolidations)
const dynamicFieldsInExecutionRequests = 3

// Compile-time check to ensure ExecutionRequests implements the necessary interfaces.
var (
	_ ssz.DynamicObject                   = (*ExecutionRequests)(nil)
	_ constraints.SSZMarshallableRootable = (*ExecutionRequests)(nil)
)

// EncodedExecutionRequest is the result of GetExecutionRequestsList which is spec defined.
type EncodedExecutionRequest = bytes.Bytes

type ExecutionRequests struct {
	Deposits       []*DepositRequest
	Withdrawals    []*WithdrawalRequest
	Consolidations []*ConsolidationRequest
}

func (e *ExecutionRequests) ValidateAfterDecodingSSZ() error {
	return errors.Join(
		DepositRequests(e.Deposits).ValidateAfterDecodingSSZ(),
		WithdrawalRequests(e.Withdrawals).ValidateAfterDecodingSSZ(),
		ConsolidationRequests(e.Consolidations).ValidateAfterDecodingSSZ(),
	)
}

// GetExecutionRequestsList introduced in pectra from the consensus spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_execution_requests_list
func GetExecutionRequestsList(er *ExecutionRequests) ([]EncodedExecutionRequest, error) {
	if er == nil {
		return nil, errors.New("nil execution requests")
	}
	result := make([]EncodedExecutionRequest, 0)

	// Process deposit requests if non-empty.
	if len(er.Deposits) > 0 {
		depositBytes, err := sszutil.MarshalItemsEIP7685(er.Deposits)
		if err != nil {
			return nil, err
		}
		combined := append([]byte{constants.DepositRequestType}, depositBytes...)
		result = append(result, combined)
	}

	// Process withdrawal requests if non-empty.
	if len(er.Withdrawals) > 0 {
		withdrawalBytes, err := sszutil.MarshalItemsEIP7685(er.Withdrawals)
		if err != nil {
			return nil, err
		}
		combined := append([]byte{constants.WithdrawalRequestType}, withdrawalBytes...)
		result = append(result, combined)
	}

	// Process consolidation requests if non-empty.
	if len(er.Consolidations) > 0 {
		consolidationBytes, err := sszutil.MarshalItemsEIP7685(er.Consolidations)
		if err != nil {
			return nil, err
		}
		combined := append([]byte{constants.ConsolidationRequestType}, consolidationBytes...)
		result = append(result, combined)
	}

	return result, nil
}

// DecodeExecutionRequests is used to decode the result from GetPayload into an ExecutionRequests.
// TODO(pectra): Change this to use []EncodedExecutionRequest as input and fix tests.
func DecodeExecutionRequests(encodedRequests [][]byte) (*ExecutionRequests, error) {
	var (
		result   ExecutionRequests
		prevType *uint8
		err      error
	)

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
		case constants.DepositRequestType:
			if result.Deposits, err = DecodeDepositRequests(data); err != nil {
				return nil, err
			}
		case constants.WithdrawalRequestType:
			if result.Withdrawals, err = DecodeWithdrawalRequests(data); err != nil {
				return nil, err
			}
		case constants.ConsolidationRequestType:
			if result.Consolidations, err = DecodeConsolidationRequests(data); err != nil {
				return nil, err
			}
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
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Deposits, constants.MaxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Withdrawals, constants.MaxWithdrawalRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsOffset(codec, &e.Consolidations, constants.MaxConsolidationRequestsPerPayload)

	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Deposits, constants.MaxDepositRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Withdrawals, constants.MaxWithdrawalRequestsPerPayload)
	ssz.DefineSliceOfStaticObjectsContent(codec, &e.Consolidations, constants.MaxConsolidationRequestsPerPayload)
}

func (e *ExecutionRequests) SizeSSZ(siz *ssz.Sizer, fixed bool) uint32 {
	size := constants.SSZOffsetSize * dynamicFieldsInExecutionRequests
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

// HashTreeRoot returns the hash tree root of the Deposits.
func (e *ExecutionRequests) HashTreeRoot() common.Root {
	return ssz.HashSequential(e)
}
