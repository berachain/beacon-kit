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
	"math/big"

	"github.com/berachain/beacon-kit/execution/client"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
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
	From                 common.Address  `json:"from"`
	To                   common.Address  `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	Input *hexutil.Bytes `json:"input"`
	// Introduced by AccessListTxType transaction.
	AccessList types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big     `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs"`
	Commitments []kzg4844.Commitment `json:"commitments"`
	Proofs      []kzg4844.Proof      `json:"proofs"`
}

type SimBlock struct {
	Calls          []TransactionArgs `json:"calls"`
	BlockOverrides *BlockOverrides   `json:"blockOverrides"`
	// TODO: in the future we could add state and block overrides here to do more complex EVM simulations
}

type BlockOverrides struct {
	Number        *hexutil.Big
	Difficulty    *hexutil.Big // No-op if we're simulating post-merge calls.
	Time          *hexutil.Uint64
	GasLimit      *hexutil.Uint64
	FeeRecipient  *common.Address
	PrevRandao    *common.Hash
	BaseFeePerGas *hexutil.Big
	BlobBaseFee   *hexutil.Big
}

type SimulateInputs struct {
	BlockStateCalls []*SimBlock `json:"blockStateCalls"`
	Validation      bool        `json:"validation"`
	TraceTransfers  bool        `json:"traceTransfers"`
}

// CallResult represents the result of an individual call in the simulated block.
type CallResult struct {
	ReturnData hexutil.Bytes   `json:"returnData"`
	Logs       json.RawMessage `json:"logs"`    // if logs structure is unknown, use RawMessage
	GasUsed    *hexutil.Uint64 `json:"gasUsed"` // pointer so we can detect absence if needed
	Status     hexutil.Uint64  `json:"status"`
}

// SimulatedBlock represents the simulated block header (with extra fields) returned
// by the eth_simulateV1 method.
type SimulatedBlock struct {
	BaseFeePerGas         *hexutil.Big    `json:"baseFeePerGas"`
	BlobGasUsed           *hexutil.Big    `json:"blobGasUsed"`
	Calls                 []CallResult    `json:"calls"`
	Difficulty            *hexutil.Big    `json:"difficulty"`
	ExcessBlobGas         *hexutil.Big    `json:"excessBlobGas"`
	ExtraData             hexutil.Bytes   `json:"extraData"`
	GasLimit              *hexutil.Uint64 `json:"gasLimit"`
	GasUsed               *hexutil.Uint64 `json:"gasUsed"`
	Hash                  common.Hash     `json:"hash"`
	LogsBloom             hexutil.Bytes   `json:"logsBloom"`
	Miner                 common.Address  `json:"miner"`
	MixHash               common.Hash     `json:"mixHash"`
	Nonce                 hexutil.Bytes   `json:"nonce"`
	Number                *hexutil.Big    `json:"number"`
	ParentBeaconBlockRoot common.Hash     `json:"parentBeaconBlockRoot"`
	ParentHash            common.Hash     `json:"parentHash"`
	ReceiptsRoot          common.Hash     `json:"receiptsRoot"`
	Sha3Uncles            common.Hash     `json:"sha3Uncles"`
	Size                  *hexutil.Uint64 `json:"size"`
	StateRoot             common.Hash     `json:"stateRoot"`
	Timestamp             *hexutil.Uint64 `json:"timestamp"`
	Transactions          []common.Hash   `json:"transactions"`
	TransactionsRoot      common.Hash     `json:"transactionsRoot"`
	Uncles                []common.Hash   `json:"uncles"`
	Withdrawals           json.RawMessage `json:"withdrawals"` // use RawMessage for now
	WithdrawalsRoot       common.Hash     `json:"withdrawalsRoot"`
}

func NewSimulationClient(client *client.EngineClient) *SimulationClient {
	return &SimulationClient{client}
}

func (c *SimulationClient) Simulate(ctx context.Context, blockNumber int64, inputs *SimulateInputs) ([]*SimulatedBlock, error) {
	var result []*SimulatedBlock
	blockNumberInput := hexutil.Uint64(blockNumber)
	err := c.engineClient.Call(ctx, &result, "eth_simulateV1", inputs, blockNumberInput)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// TxsToSimBlock Transactions must use Dynamic Fee values, i.e. non legacy txs
func TxsToSimBlock(chainId uint64, txs []*gethprimitives.Transaction) ([]TransactionArgs, error) {
	// TODO: use the LatestSigner based on a Geth ChainConfig parsed from an EL Genesis File.
	// Transactions must use Dynamic Fee values, i.e. non legacy txs
	signer := types.NewCancunSigner(big.NewInt(int64(chainId)))
	calls := make([]TransactionArgs, len(txs))
	for i, tx := range txs {
		sender, err := signer.Sender(tx)
		if err != nil {
			return nil, err
		}
		gas := hexutil.Uint64(tx.Gas())
		gasFeeCap := hexutil.Big(*tx.GasFeeCap())
		gasTipCap := hexutil.Big(*tx.GasTipCap())
		value := hexutil.Big(*tx.Value())
		nonce := hexutil.Uint64(tx.Nonce())
		data := hexutil.Bytes(tx.Data())
		chainIdHex := hexutil.Big(*big.NewInt(int64(chainId)))
		blobGasFeeCap := hexutil.Big(*tx.GasFeeCap())
		call := TransactionArgs{
			From:                 sender,
			To:                   *tx.To(),
			Gas:                  &gas,
			MaxFeePerGas:         &gasFeeCap,
			MaxPriorityFeePerGas: &gasTipCap,
			Value:                &value,
			Nonce:                &nonce,
			Input:                &data,
			AccessList:           tx.AccessList(),
			ChainID:              &chainIdHex,
			BlobFeeCap:           &blobGasFeeCap,
			BlobHashes:           tx.BlobHashes(),
		}
		if tx.BlobTxSidecar() != nil {
			call.Blobs = tx.BlobTxSidecar().Blobs
			call.Commitments = tx.BlobTxSidecar().Commitments
			call.Proofs = tx.BlobTxSidecar().Proofs
		}
		calls[i] = call
	}
	return calls, nil
}
