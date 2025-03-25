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
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package eip7002

import (
	"bytes"
	"context"
	"encoding/binary"
	"math/big"

	"github.com/berachain/beacon-kit/errors"
	beaconbytes "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/ethereum/go-ethereum/params"
)

type feeOpts struct {
	To string `json:"to"`
}

// GetWithdrawalFee returns the withdrawal fee in wei. See https://eips.ethereum.org/EIPS/eip-7002 for more.
func GetWithdrawalFee(ctx context.Context, client rpcClient) (*big.Int, error) {
	var result string
	feeInput := &feeOpts{
		To: params.WithdrawalQueueAddress.String(),
	}
	err := client.Call(ctx, &result, "eth_call", feeInput)
	if err != nil {
		return nil, err
	}
	n, ok := new(big.Int).SetString(result, 0)
	if !ok {
		return nil, errors.New("error converting hex string to big.Int")
	}
	return n, nil
}

// CreateWithdrawalRequestData returns the request body formatted as defined by the EIP-7002 specification.
func CreateWithdrawalRequestData(blsPubKey crypto.BLSPubkey, withdrawAmount math.Gwei) (beaconbytes.Bytes, error) {
	// Create a buffer to hold the packed encoding.
	var packed bytes.Buffer
	if _, err := packed.Write(blsPubKey[:]); err != nil {
		return nil, err
	}
	// Write the uint64 value in big-endian order.
	if err := binary.Write(&packed, binary.BigEndian, withdrawAmount); err != nil {
		return nil, err
	}
	return packed.Bytes(), nil
}
