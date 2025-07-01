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

	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
)

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// Compile-time check to ensure DepositRequests implements the necessary interfaces.
var _ constraints.SSZMarshaler = (*DepositRequests)(nil)

// DepositRequests is used for SSZ unmarshalling a list of DepositRequest
type DepositRequests []*DepositRequest

// MarshalSSZ marshals the Deposits object to SSZ format by encoding each deposit individually.
func (dr DepositRequests) MarshalSSZ() ([]byte, error) {
	return sszutil.MarshalItemsEIP7685(dr)
}

// ValidateAfterDecodingSSZ validates the DepositRequests object after decoding.
func (dr DepositRequests) ValidateAfterDecodingSSZ() error {
	if len(dr) > constants.MaxDepositRequestsPerPayload {
		return fmt.Errorf(
			"invalid number of deposit requests, got %d max %d",
			len(dr), constants.MaxDepositRequestsPerPayload,
		)
	}
	return nil
}

// DecodeDepositRequests decodes SSZ data by decoding each request individually.
func DecodeDepositRequests(data []byte) (DepositRequests, error) {
	maxSize := constants.MaxDepositRequestsPerPayload * depositSize
	if len(data) > maxSize {
		return nil, fmt.Errorf(
			"invalid deposit requests SSZ size, requests should not be more than the max per "+
				"payload, got %d max %d", len(data), maxSize,
		)
	}
	if len(data) < depositSize {
		return nil, fmt.Errorf(
			"invalid deposit requests SSZ size, got %d expected at least %d",
			len(data), depositSize,
		)
	}

	// Use the EIP-7685 unmarshalItems helper.
	return sszutil.UnmarshalItemsEIP7685(
		data,
		depositSize,
		func() *DepositRequest { return new(DepositRequest) },
	)
}
