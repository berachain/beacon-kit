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

	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

const maxDepositRequestsPerPayload = 8192

// DepositRequest is introduced in EIP6110 which is currently not processed.
type DepositRequest = Deposit

// DepositRequests is used for SSZ unmarshalling a list of DepositRequest
type DepositRequests []*DepositRequest

// MarshalSSZ marshals the Deposits object to SSZ format by encoding each deposit individually.
func (dr *DepositRequests) MarshalSSZ() ([]byte, error) {
	return constraints.MarshalItems[*DepositRequest](*dr)
}

// NewFromSSZ decodes SSZ data into a Deposits object by decoding each deposit individually.
func (dr *DepositRequests) NewFromSSZ(data []byte) (*DepositRequests, error) {
	maxSize := maxDepositRequestsPerPayload * DepositSize
	if len(data) > maxSize {
		return nil, fmt.Errorf("invalid deposit requests SSZ size, requests should not be more than the max per payload, got %d max %d", len(data), maxSize)
	}
	depositSize := int(ssz.Size(&Deposit{}))
	if len(data) < depositSize {
		return nil, fmt.Errorf("invalid deposit requests SSZ size, got %d expected at least %d", len(data), depositSize)
	}
	// Use the generic unmarshalItems helper.
	items, err := constraints.UnmarshalItems[*DepositRequest](data, depositSize, func() *Deposit { return new(DepositRequest) })
	if err != nil {
		return nil, err
	}
	deposits := DepositRequests(items)
	return &deposits, nil
}
