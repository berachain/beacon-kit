//go:build simulated

//
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

package execution

import (
	"context"

	"github.com/berachain/beacon-kit/execution/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

type SimulationClient struct {
	engineClient *client.EngineClient
}

// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
// Taken from https://github.com/ethereum/go-ethereum/blob/e6f3ce7b168b8f346de621a8f60d2fa57c2ebfb0/internal/ethapi/transaction_args.go#L42
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs"`
	Commitments []kzg4844.Commitment `json:"commitments"`
	Proofs      []kzg4844.Proof      `json:"proofs"`
}

type SimBlock struct {
	Calls       []TransactionArgs
	BlockNumber int64
}

func NewSimulationClient(client *client.EngineClient) *SimulationClient {
	return &SimulationClient{client}
}

func (c *SimulationClient) Simulate(ctx context.Context, simBlock SimBlock) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := c.engineClient.Call(ctx, &result, "eth_simulateV1", simBlock)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//func TxAndSidecarsToTransactionArgs(txs []*gethprimitives.Transaction, sidecars []*types.BlobTxSidecar) TransactionArgs {
//	calls := make([]TransactionArgs, len(txs))
//	for i, tx := range txs {
//		TransactionArgs{
//			From:                 nil,
//			To:                   tx.To(),
//			Gas:                  tx.Gas(),
//			GasPrice:             nil,
//			MaxFeePerGas:         nil,
//			MaxPriorityFeePerGas: nil,
//			Value:                nil,
//			Nonce:                nil,
//			Input:                nil,
//			AccessList:           nil,
//			ChainID:              nil,
//			BlobFeeCap:           nil,
//			BlobHashes:           nil,
//			Blobs:                nil,
//			Commitments:          nil,
//			Proofs:               nil,
//		}
//	}
//}
