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
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
)

// WithdrawalRequest is introduced in EIP7002 which we use for withdrawals.
type WithdrawalRequest struct {
	SourceAddress   common.ExecutionAddress
	ValidatorPubKey crypto.BLSPubkey
	Amount          math.Gwei
}

// WithdrawalRequests is used for SSZ unmarshalling a list of WithdrawalRequest
type WithdrawalRequests []*WithdrawalRequest

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

func (w *WithdrawalRequest) NewFromSSZ(buf []byte) (*WithdrawalRequest, error) {
	if w == nil {
		w = &WithdrawalRequest{}
	}
	return w, ssz.DecodeFromBytes(buf, w)
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (w *WithdrawalRequest) HashTreeRoot() common.Root {
	return ssz.HashSequential(w)
}

// MarshalSSZ marshals the WithdrawalRequests object to SSZ format by encoding each deposit individually.
func (wr *WithdrawalRequests) MarshalSSZ() ([]byte, error) {
	return constraints.MarshalItems[*WithdrawalRequest](*wr)
}

// NewFromSSZ decodes SSZ data into a Deposits object by decoding each deposit individually.
func (wr *WithdrawalRequests) NewFromSSZ(data []byte) (*WithdrawalRequests, error) {
	maxSize := maxWithdrawalRequestsPerPayload * sszWithdrawRequestSize
	if len(data) > maxSize {
		return nil, fmt.Errorf(
			"invalid withdrawal requests SSZ size, requests should not be more "+
				"than the max per payload, got %d max %d", len(data), maxSize,
		)
	}
	withdrawalSize := int(ssz.Size(&WithdrawalRequest{}))
	if len(data) < withdrawalSize {
		return nil, fmt.Errorf("invalid withdrawal requests SSZ size, got %d expected at least %d", len(data), withdrawalSize)
	}
	// Use the generic unmarshalItems helper.
	items, err := constraints.UnmarshalItems[*WithdrawalRequest](
		data,
		withdrawalSize,
		func() *WithdrawalRequest { return new(WithdrawalRequest) },
	)
	if err != nil {
		return nil, err
	}
	deposits := WithdrawalRequests(items)
	return &deposits, nil
}
