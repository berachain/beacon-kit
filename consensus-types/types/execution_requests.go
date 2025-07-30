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
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
	ssz "github.com/ferranbt/fastssz"
)

// 3 since three dynamic objects (Deposits, Withdrawals, Consolidations)
const dynamicFieldsInExecutionRequests = 3

// SizeSSZ returns the SSZ encoded size in bytes.
func (e *ExecutionRequests) SizeSSZ() int {
	size := int(constants.SSZOffsetSize * dynamicFieldsInExecutionRequests)
	size += len(e.Deposits) * 192       // Each deposit request is 192 bytes
	size += len(e.Withdrawals) * 76     // Each withdrawal request is 76 bytes
	size += len(e.Consolidations) * 116 // Each consolidation request is 116 bytes
	return size
}

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

// HashTreeRoot returns the hash tree root of the ExecutionRequests.
func (e *ExecutionRequests) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	defer ssz.DefaultHasherPool.Put(hh)
	if err := e.HashTreeRootWith(hh); err != nil {
		return [32]byte{}, err
	}
	return hh.HashRoot()

}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZ marshals the ExecutionRequests object.
func (e *ExecutionRequests) MarshalSSZ() ([]byte, error) {
	// Initialize empty slices if nil
	if e.Deposits == nil {
		e.Deposits = make([]*DepositRequest, 0)
	}
	if e.Withdrawals == nil {
		e.Withdrawals = make([]*WithdrawalRequest, 0)
	}
	if e.Consolidations == nil {
		e.Consolidations = make([]*ConsolidationRequest, 0)
	}

	// Calculate size
	size := 12 // 3 fields * 4 bytes offset each

	// Add dynamic content sizes
	for _, d := range e.Deposits {
		depositBytes, err := d.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		size += len(depositBytes)
	}
	for _, w := range e.Withdrawals {
		withdrawalBytes, err := w.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		size += len(withdrawalBytes)
	}
	for _, c := range e.Consolidations {
		consolidationBytes, err := c.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		size += len(consolidationBytes)
	}

	// Create buffer
	buf := make([]byte, size)
	offset := 12

	// Write offsets
	// Deposits offset
	binary.LittleEndian.PutUint32(buf[0:4], uint32(offset))
	for _, d := range e.Deposits {
		depositBytes, _ := d.MarshalSSZ()
		offset += len(depositBytes)
	}

	// Withdrawals offset
	binary.LittleEndian.PutUint32(buf[4:8], uint32(offset))
	for _, w := range e.Withdrawals {
		withdrawalBytes, _ := w.MarshalSSZ()
		offset += len(withdrawalBytes)
	}

	// Consolidations offset
	binary.LittleEndian.PutUint32(buf[8:12], uint32(offset))

	// Write content
	offset = 12
	for _, d := range e.Deposits {
		depositBytes, _ := d.MarshalSSZ()
		copy(buf[offset:], depositBytes)
		offset += len(depositBytes)
	}
	for _, w := range e.Withdrawals {
		withdrawalBytes, _ := w.MarshalSSZ()
		copy(buf[offset:], withdrawalBytes)
		offset += len(withdrawalBytes)
	}
	for _, c := range e.Consolidations {
		consolidationBytes, _ := c.MarshalSSZ()
		copy(buf[offset:], consolidationBytes)
		offset += len(consolidationBytes)
	}

	return buf, nil
}

// MarshalSSZTo ssz marshals the ExecutionRequests object to a target array.
func (e *ExecutionRequests) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := e.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the ExecutionRequests object.
func (e *ExecutionRequests) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 12 {
		return errors.New("buffer too short for ExecutionRequests")
	}

	// Initialize empty slices
	e.Deposits = make([]*DepositRequest, 0)
	e.Withdrawals = make([]*WithdrawalRequest, 0)
	e.Consolidations = make([]*ConsolidationRequest, 0)

	// Read offsets
	depositsOffset := binary.LittleEndian.Uint32(buf[0:4])
	withdrawalsOffset := binary.LittleEndian.Uint32(buf[4:8])
	consolidationsOffset := binary.LittleEndian.Uint32(buf[8:12])

	// Validate offsets
	if depositsOffset < 12 || depositsOffset > uint32(len(buf)) {
		return errors.New("invalid deposits offset")
	}
	if withdrawalsOffset < depositsOffset || withdrawalsOffset > uint32(len(buf)) {
		return errors.New("invalid withdrawals offset")
	}
	if consolidationsOffset < withdrawalsOffset || consolidationsOffset > uint32(len(buf)) {
		return errors.New("invalid consolidations offset")
	}

	// Unmarshal deposits
	if depositsOffset < withdrawalsOffset {
		depositsData := buf[depositsOffset:withdrawalsOffset]
		if len(depositsData) > 0 {
			deposits, err := DecodeDepositRequests(depositsData)
			if err != nil {
				return err
			}
			e.Deposits = deposits
		}
	}

	// Unmarshal withdrawals
	if withdrawalsOffset < consolidationsOffset {
		withdrawalsData := buf[withdrawalsOffset:consolidationsOffset]
		if len(withdrawalsData) > 0 {
			withdrawals, err := DecodeWithdrawalRequests(withdrawalsData)
			if err != nil {
				return err
			}
			e.Withdrawals = withdrawals
		}
	}

	// Unmarshal consolidations
	if consolidationsOffset < uint32(len(buf)) {
		consolidationsData := buf[consolidationsOffset:]
		if len(consolidationsData) > 0 {
			consolidations, err := DecodeConsolidationRequests(consolidationsData)
			if err != nil {
				return err
			}
			e.Consolidations = consolidations
		}
	}

	return nil
}

// SizeSSZFastSSZ returns the ssz encoded size in bytes for the ExecutionRequests (fastssz).
// TODO: Rename to SizeSSZ() once karalabe/ssz is fully removed.
func (e *ExecutionRequests) SizeSSZFastSSZ() (size int) {
	size = 12 // 3 fields * 4 bytes offset each

	// Add dynamic sizes
	for range e.Deposits {
		size += 192 // deposit size
	}
	for range e.Withdrawals {
		size += 76 // withdrawal request size
	}
	for range e.Consolidations {
		size += 116 // consolidation request size
	}

	return
}

// HashTreeRootWith ssz hashes the ExecutionRequests object with a hasher.
func (e *ExecutionRequests) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Deposits'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Deposits))
		if num > constants.MaxDepositRequestsPerPayload {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range e.Deposits {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxDepositRequestsPerPayload)
	}

	// Field (1) 'Withdrawals'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Withdrawals))
		if num > constants.MaxWithdrawalRequestsPerPayload {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range e.Withdrawals {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxWithdrawalRequestsPerPayload)
	}

	// Field (2) 'Consolidations'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Consolidations))
		if num > constants.MaxConsolidationRequestsPerPayload {
			return ssz.ErrIncorrectListSize
		}
		for _, elem := range e.Consolidations {
			if err := elem.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, constants.MaxConsolidationRequestsPerPayload)
	}

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the ExecutionRequests object.
func (e *ExecutionRequests) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(e)
}
