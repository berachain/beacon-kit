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

package ethclient

import (
	"context"
	"math/big"

	"github.com/berachain/beacon-kit/gethlib/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethclient "github.com/ethereum/go-ethereum/ethclient"
)

// Client is a wrapper around go-ethereum's ethclient which handles unmarhsalling
// Berachain blocks and transactions.
type Client struct {
	c *gethclient.Client
}

// Wrap wraps a go-ethereum's ethclient and returns a Berachain-specific ethclient.
func Wrap(c *gethclient.Client) *Client {
	return &Client{c}
}

// BlockByNumber overrides the original method to unmarshal the block.
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	panic("TODO: implement by `getBlock`")
}

// BlockByHash overrides the original method to unmarshal the block.
func (c *Client) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	panic("TODO: implement by `getBlock`")
}

// TransactionByHash overrides the original method to unmarshal the transaction.
func (c *Client) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (tx *types.Transaction, isPending bool, err error) {
	panic("TODO: implement")
}

// TransactionInBlock overrides the original method to unmarshal the transaction.
func (c *Client) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*types.Transaction, error) {
	panic("TODO: implement")
}
