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

package eip7685

import (
	"bytes"
	"context"
	"encoding/binary"

	beaconbytes "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// WithdrawalRequestPredeployAddress is a spec defined address for the withdrawal contract
const WithdrawalRequestPredeployAddress = "0x00000961Ef480Eb55e80D19ad83579A64c007002"
const ethCall = "eth_call"

// GetWithdrawalFee returns the withdrawal fee in wei. See https://eips.ethereum.org/EIPS/eip-7002 for more.
func GetWithdrawalFee(ctx context.Context, client rpcClient) (math.U64, error) {
	var result math.U64
	err := client.Call(ctx, &result, ethCall, WithdrawalRequestPredeployAddress)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// CreateWithdrawalRequestData returns the request body formatted as defined by the EIP-7002 specification.
func CreateWithdrawalRequestData(blsPubKey crypto.BLSPubkey, withdrawAmount math.U64) (beaconbytes.Bytes, error) {
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
