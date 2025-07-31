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

//go:generate sszgen -path withdrawal_request.go -objs WithdrawalRequest -output withdrawal_request_sszgen.go -include ../../primitives/common,../../primitives/crypto,../../primitives/math,../../primitives/bytes

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
	"github.com/berachain/beacon-kit/primitives/math"
)

const sszWithdrawRequestSize = 76 // ExecutionAddress = 20, ValidatorPubKey = 48, Amount = 8

// WithdrawalRequest is introduced in EIP7002 which we use for withdrawals.
type WithdrawalRequest struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

// ValidateAfterDecodingSSZ validates the withdrawal request after decoding.
func (w *WithdrawalRequest) ValidateAfterDecodingSSZ() error {
	return nil
}

/* -------------------------------------------------------------------------- */
/*                       Withdrawal Requests SSZ                              */
/* -------------------------------------------------------------------------- */

// Compile-time check to ensure WithdrawalRequests implements the necessary interfaces.
var _ constraints.SSZMarshaler = (*WithdrawalRequests)(nil)

// WithdrawalRequests is used for SSZ unmarshalling a list of WithdrawalRequest
type WithdrawalRequests []*WithdrawalRequest

// MarshalSSZ marshals the WithdrawalRequests object to SSZ format by encoding each deposit individually.
func (wr WithdrawalRequests) MarshalSSZ() ([]byte, error) {
	return sszutil.MarshalItemsEIP7685(wr)
}

// ValidateAfterDecodingSSZ validates the WithdrawalRequests object after decoding.
func (wr WithdrawalRequests) ValidateAfterDecodingSSZ() error {
	if len(wr) > constants.MaxWithdrawalRequestsPerPayload {
		return fmt.Errorf(
			"invalid number of withdrawal requests, got %d max %d",
			len(wr), constants.MaxWithdrawalRequestsPerPayload,
		)
	}
	return nil
}

// DecodeWithdrawalRequests decodes SSZ data by decoding each withdrawal individually.
func DecodeWithdrawalRequests(data []byte) (WithdrawalRequests, error) {
	maxSize := constants.MaxWithdrawalRequestsPerPayload * sszWithdrawRequestSize
	if len(data) > maxSize {
		return nil, fmt.Errorf(
			"invalid withdrawal requests SSZ size, requests should not be more "+
				"than the max per payload, got %d max %d", len(data), maxSize,
		)
	}
	if len(data) < sszWithdrawRequestSize {
		return nil, fmt.Errorf(
			"invalid withdrawal requests SSZ size, got %d expected at least %d",
			len(data), sszWithdrawRequestSize,
		)
	}

	// Use the EIP-7685 unmarshalItems helper.
	return sszutil.UnmarshalItemsEIP7685(
		data,
		sszWithdrawRequestSize,
		func() *WithdrawalRequest { return new(WithdrawalRequest) },
	)
}
