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

// SimulationClient calls `eth_simulateV1` on a Geth node.
type SimulationClient struct {
	engineClient *client.EngineClient
}

// TransactionArgs represents the fields needed to construct a dynamic-fee transaction.
type TransactionArgs struct {
	From                 common.Address  `json:"from"`
	To                   common.Address  `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`
	Input                *hexutil.Bytes  `json:"input"`

	AccessList types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big     `json:"chainId,omitempty"`

	// BlobTxType fields.
	BlobFeeCap  *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes  []common.Hash `json:"blobVersionedHashes,omitempty"`
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

// BlockOverrides holds optional block-level overrides for simulation.
type BlockOverrides struct {
	Number        *hexutil.Big    `json:"number,omitempty"`
	Difficulty    *hexutil.Big    `json:"difficulty,omitempty"`
	Time          *hexutil.Uint64 `json:"time,omitempty"`
	GasLimit      *hexutil.Uint64 `json:"gasLimit,omitempty"`
	FeeRecipient  *common.Address `json:"feeRecipient,omitempty"`
	PrevRandao    *common.Hash    `json:"prevRandao,omitempty"`
	BaseFeePerGas *hexutil.Big    `json:"baseFeePerGas,omitempty"`
	BlobBaseFee   *hexutil.Big    `json:"blobBaseFee,omitempty"`
	BeaconRoot    *common.Hash    `json:"beaconRoot,omitempty"`
	Withdrawals   gethprimitives.Withdrawals
}

// SimBlock is a block containing calls and optional overrides for simulation.
type SimBlock struct {
	Calls          []TransactionArgs `json:"calls"`
	BlockOverrides *BlockOverrides   `json:"blockOverrides"`
}

// SimOpts groups all parameters for `eth_simulateV1`.
type SimOpts struct {
	BlockStateCalls []*SimBlock `json:"blockStateCalls"`
	Validation      bool        `json:"validation"`
	TraceTransfers  bool        `json:"traceTransfers"`
}

// CallResult describes the outcome of an individual simulated transaction.
type CallResult struct {
	ReturnData hexutil.Bytes   `json:"returnData"`
	Logs       json.RawMessage `json:"logs"`
	GasUsed    *hexutil.Uint64 `json:"gasUsed"`
	Status     hexutil.Uint64  `json:"status"`
}

// SimulatedBlock is the response structure returned by `eth_simulateV1`.
type SimulatedBlock struct {
	BaseFeePerGas         *hexutil.Big      `json:"baseFeePerGas"`
	BlobGasUsed           *hexutil.Big      `json:"blobGasUsed"`
	Calls                 []CallResult      `json:"calls"`
	Difficulty            *hexutil.Big      `json:"difficulty"`
	ExcessBlobGas         *hexutil.Big      `json:"excessBlobGas"`
	ExtraData             hexutil.Bytes     `json:"extraData"`
	GasLimit              *hexutil.Uint64   `json:"gasLimit"`
	GasUsed               *hexutil.Uint64   `json:"gasUsed"`
	Hash                  common.Hash       `json:"hash"`
	LogsBloom             hexutil.Bytes     `json:"logsBloom"`
	Miner                 common.Address    `json:"miner"`
	MixHash               common.Hash       `json:"mixHash"`
	Nonce                 hexutil.Bytes     `json:"nonce"`
	Number                *hexutil.Big      `json:"number"`
	ParentBeaconBlockRoot common.Hash       `json:"parentBeaconBlockRoot"`
	ParentHash            common.Hash       `json:"parentHash"`
	ReceiptsRoot          common.Hash       `json:"receiptsRoot"`
	Sha3Uncles            common.Hash       `json:"sha3Uncles"`
	Size                  *hexutil.Uint64   `json:"size"`
	StateRoot             common.Hash       `json:"stateRoot"`
	Timestamp             *hexutil.Uint64   `json:"timestamp"`
	Transactions          []common.Hash     `json:"transactions"`
	TransactionsRoot      common.Hash       `json:"transactionsRoot"`
	Uncles                []common.Hash     `json:"uncles"`
	Withdrawals           types.Withdrawals `json:"withdrawals"`
	WithdrawalsRoot       common.Hash       `json:"withdrawalsRoot"`
}

// NewSimulationClient returns a client for eth_simulateV1 calls.
func NewSimulationClient(cli *client.EngineClient) *SimulationClient {
	return &SimulationClient{cli}
}

// Simulate calls `eth_simulateV1` with the provided block number and options.
func (c *SimulationClient) Simulate(
	ctx context.Context,
	blockNumber int64,
	opts *SimOpts,
) ([]*SimulatedBlock, error) {

	var result []*SimulatedBlock
	blockNumberInput := hexutil.Uint64(blockNumber)

	if err := c.engineClient.Call(ctx, &result, "eth_simulateV1", opts, blockNumberInput); err != nil {
		return nil, err
	}
	return result, nil
}

// TxsToTransactionArgs converts a slice of Geth transactions to TransactionArgs suitable for simulation.
// The transactions must be dynamic-fee (EIP-1559 or EIP-4844) type.
// TODO: use the LatestSigner based on a Geth ChainConfig parsed from an EL Genesis File as have currently hardcoded Cancun Signer.
func TxsToTransactionArgs(chainID uint64, txs []*gethprimitives.Transaction) ([]TransactionArgs, error) {
	signer := types.NewCancunSigner(big.NewInt(int64(chainID)))
	args := make([]TransactionArgs, len(txs))

	for i, tx := range txs {
		sender, err := signer.Sender(tx)
		if err != nil {
			return nil, err
		}

		gas := hexutil.Uint64(tx.Gas())
		feeCap := hexutil.Big(*tx.GasFeeCap())
		tipCap := hexutil.Big(*tx.GasTipCap())
		val := hexutil.Big(*tx.Value())
		nonce := hexutil.Uint64(tx.Nonce())
		data := hexutil.Bytes(tx.Data())
		chainIDHex := hexutil.Big(*big.NewInt(int64(chainID)))

		call := TransactionArgs{
			From:                 sender,
			To:                   *tx.To(),
			Gas:                  &gas,
			MaxFeePerGas:         &feeCap,
			MaxPriorityFeePerGas: &tipCap,
			Value:                &val,
			Nonce:                &nonce,
			Input:                &data,
			AccessList:           tx.AccessList(),
			ChainID:              &chainIDHex,
		}

		if sidecar := tx.BlobTxSidecar(); sidecar != nil {
			blobCap := hexutil.Big(*tx.BlobGasFeeCap())
			call.BlobHashes = tx.BlobHashes()
			call.BlobFeeCap = &blobCap
			call.Blobs = sidecar.Blobs
			call.Commitments = sidecar.Commitments
			call.Proofs = sidecar.Proofs
		}

		args[i] = call
	}
	return args, nil
}
