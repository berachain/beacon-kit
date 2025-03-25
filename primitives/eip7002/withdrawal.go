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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
