// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package ethclient

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ChainID retrieves the current chain ID.
func (ec *Client[ExecutionPayloadT, _]) ChainID(
	ctx context.Context,
) (math.U64, error) {
	var result math.U64
	if err := ec.Call(ctx, &result, "eth_chainId"); err != nil {
		return 0, err
	}
	return result, nil
}

// GetLogsAtBlockNumber retrieves logs for a block number, contract address,
// and optional topics.
//
// The topics list restricts matches to particular event topics. Topics matches
// a prefix of that list. An empty element slice matches any topic. Non-empty
// elements represent an alternative that matches any of the contained topics.
//
// Examples:
// {} or nil          matches any topic list
// {{A}}              matches topic A in first position
// {{}, {B}}          matches any topic in first position AND B in second
//
//	position
//
// {{A}, {B}}         matches topic A in first position AND B in second
//
//	position
//
// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in
//
//	second position
func (ec *Client[ExecutionPayloadT, LogT]) GetLogsAtBlockNumber(
	ctx context.Context,
	number math.U64,
	address common.ExecutionAddress,
	topics [][]common.ExecutionHash,
) ([]LogT, error) {
	var result []LogT
	return result,
		ec.Call(ctx,
			&result, "eth_getLogs", map[string]interface{}{
				"fromBlock": number.Hex(),
				"toBlock":   number.Hex(),
				"address":   address,
				"topics":    topics,
			})
}
