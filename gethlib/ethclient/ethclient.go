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
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/berachain/beacon-kit/gethlib/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	coretypes "github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client is a wrapper around go-ethereum's ethclient which handles unmarhsalling
// Berachain blocks and transactions.
type Client struct {
	*gethclient.Client
}

// Wrap wraps a go-ethereum's ethclient and returns a Berachain-specific ethclient.
func Wrap(c *gethclient.Client) *Client {
	return &Client{c}
}

// BlockByNumber overrides the original method to unmarshal the block.
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.getBlock(ctx, "eth_getBlockByNumber", toBlockNumArg(number), true)
}

// BlockByHash overrides the original method to unmarshal the block.
func (c *Client) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	return c.getBlock(ctx, "eth_getBlockByHash", hash, true)
}

// TransactionByHash overrides the original method to unmarshal the transaction.
func (c *Client) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (tx *types.Transaction, isPending bool, err error) {
	var jsonTx *rpcTransaction
	err = c.Client.Client().CallContext(ctx, &jsonTx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	}
	if jsonTx == nil {
		return nil, false, ethereum.NotFound
	}
	if err := ensureTransactionHasRequiredSignature(jsonTx.tx); err != nil {
		return nil, false, err
	}
	return jsonTx.tx, jsonTx.BlockNumber == nil, nil
}

// TransactionInBlock overrides the original method to unmarshal the transaction.
func (c *Client) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*types.Transaction, error) {
	var jsonTx *rpcTransaction
	err := c.Client.Client().CallContext(ctx, &jsonTx, "eth_getTransactionByBlockHashAndIndex", blockHash, hexutil.Uint64(index))
	if err != nil {
		return nil, err
	}
	if jsonTx == nil {
		return nil, ethereum.NotFound
	}
	if err := ensureTransactionHasRequiredSignature(jsonTx.tx); err != nil {
		return nil, err
	}
	return jsonTx.tx, nil
}

type rpcBlock struct {
	Hash         *gethcommon.Hash        `json:"hash"`
	Transactions []rpcTransaction        `json:"transactions"`
	UncleHashes  []gethcommon.Hash       `json:"uncles"`
	Withdrawals  []*coretypes.Withdrawal `json:"withdrawals,omitempty"`
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string             `json:"blockNumber,omitempty"`
	BlockHash   *gethcommon.Hash    `json:"blockHash,omitempty"`
	From        *gethcommon.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

type rpcHeader struct {
	ParentHash           *gethcommon.Hash       `json:"parentHash"`
	UncleHash            *gethcommon.Hash       `json:"sha3Uncles"`
	Coinbase             *gethcommon.Address    `json:"miner"`
	Root                 *gethcommon.Hash       `json:"stateRoot"`
	TxHash               *gethcommon.Hash       `json:"transactionsRoot"`
	ReceiptHash          *gethcommon.Hash       `json:"receiptsRoot"`
	Bloom                *coretypes.Bloom       `json:"logsBloom"`
	Difficulty           *hexutil.Big           `json:"difficulty"`
	Number               *hexutil.Big           `json:"number"`
	GasLimit             *hexutil.Uint64        `json:"gasLimit"`
	GasUsed              *hexutil.Uint64        `json:"gasUsed"`
	Time                 *hexutil.Uint64        `json:"timestamp"`
	Extra                *hexutil.Bytes         `json:"extraData"`
	MixDigest            *gethcommon.Hash       `json:"mixHash"`
	Nonce                *coretypes.BlockNonce  `json:"nonce"`
	BaseFee              *hexutil.Big           `json:"baseFeePerGas,omitempty"`
	WithdrawalsHash      *gethcommon.Hash       `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed          *hexutil.Uint64        `json:"blobGasUsed,omitempty"`
	ExcessBlobGas        *hexutil.Uint64        `json:"excessBlobGas,omitempty"`
	ParentBeaconRoot     *gethcommon.Hash       `json:"parentBeaconBlockRoot,omitempty"`
	RequestsHash         *gethcommon.Hash       `json:"requestsHash,omitempty"`
	ParentProposerPubkey *types.ExecutionPubkey `json:"parentProposerPubkey,omitempty"`
}

func (h *rpcHeader) toHeader() (*types.Header, error) {
	if h == nil {
		return nil, nil
	}
	if h.ParentHash == nil {
		return nil, errors.New("missing required field 'parentHash' for Header")
	}
	if h.UncleHash == nil {
		return nil, errors.New("missing required field 'sha3Uncles' for Header")
	}
	if h.Root == nil {
		return nil, errors.New("missing required field 'stateRoot' for Header")
	}
	if h.TxHash == nil {
		return nil, errors.New("missing required field 'transactionsRoot' for Header")
	}
	if h.ReceiptHash == nil {
		return nil, errors.New("missing required field 'receiptsRoot' for Header")
	}
	if h.Bloom == nil {
		return nil, errors.New("missing required field 'logsBloom' for Header")
	}
	if h.Difficulty == nil {
		return nil, errors.New("missing required field 'difficulty' for Header")
	}
	if h.Number == nil {
		return nil, errors.New("missing required field 'number' for Header")
	}
	if h.GasLimit == nil {
		return nil, errors.New("missing required field 'gasLimit' for Header")
	}
	if h.GasUsed == nil {
		return nil, errors.New("missing required field 'gasUsed' for Header")
	}
	if h.Time == nil {
		return nil, errors.New("missing required field 'timestamp' for Header")
	}
	if h.Extra == nil {
		return nil, errors.New("missing required field 'extraData' for Header")
	}

	head := &types.Header{
		ParentHash:  *h.ParentHash,
		UncleHash:   *h.UncleHash,
		Root:        *h.Root,
		TxHash:      *h.TxHash,
		ReceiptHash: *h.ReceiptHash,
		Bloom:       *h.Bloom,
		Difficulty:  (*big.Int)(h.Difficulty),
		Number:      (*big.Int)(h.Number),
		GasLimit:    uint64(*h.GasLimit),
		GasUsed:     uint64(*h.GasUsed),
		Time:        uint64(*h.Time),
		Extra:       *h.Extra,
	}
	if h.Coinbase != nil {
		head.Coinbase = *h.Coinbase
	}
	if h.MixDigest != nil {
		head.MixDigest = *h.MixDigest
	}
	if h.Nonce != nil {
		head.Nonce = *h.Nonce
	}
	if h.BaseFee != nil {
		head.BaseFee = (*big.Int)(h.BaseFee)
	}
	if h.WithdrawalsHash != nil {
		hash := *h.WithdrawalsHash
		head.WithdrawalsHash = &hash
	}
	if h.BlobGasUsed != nil {
		blobGasUsed := uint64(*h.BlobGasUsed)
		head.BlobGasUsed = &blobGasUsed
	}
	if h.ExcessBlobGas != nil {
		excessBlobGas := uint64(*h.ExcessBlobGas)
		head.ExcessBlobGas = &excessBlobGas
	}
	if h.ParentBeaconRoot != nil {
		root := *h.ParentBeaconRoot
		head.ParentBeaconRoot = &root
	}
	if h.RequestsHash != nil {
		hash := *h.RequestsHash
		head.RequestsHash = &hash
	}
	if h.ParentProposerPubkey != nil {
		pubkey := *h.ParentProposerPubkey
		head.ParentProposerPubkey = &pubkey
	}
	return head, nil
}

func (c *Client) getBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	if err := c.Client.Client().CallContext(ctx, &raw, method, args...); err != nil {
		return nil, err
	}

	// Decode header and transactions.
	var rawHead *rpcHeader
	if err := json.Unmarshal(raw, &rawHead); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if rawHead == nil {
		return nil, ethereum.NotFound
	}
	head, err := rawHead.toHeader()
	if err != nil {
		return nil, err
	}

	var body rpcBlock
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Pending blocks don't return a block hash. Compute it for consistency.
	if body.Hash == nil {
		tmp := head.Hash()
		body.Hash = &tmp
	}

	// Quick-verify transaction and uncle lists.
	if head.UncleHash == coretypes.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != coretypes.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == coretypes.EmptyTxsHash && len(body.Transactions) > 0 {
		return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != coretypes.EmptyTxsHash && len(body.Transactions) == 0 {
		return nil, errors.New("server returned empty transaction list but block header indicates transactions")
	}

	// Load uncles because they are not included in the block response.
	var uncles []*types.Header
	if len(body.UncleHashes) > 0 {
		uncles = make([]*types.Header, len(body.UncleHashes))
		rawUncles := make([]*rpcHeader, len(body.UncleHashes))
		reqs := make([]rpc.BatchElem, len(body.UncleHashes))
		for i := range reqs {
			reqs[i] = rpc.BatchElem{
				Method: "eth_getUncleByBlockHashAndIndex",
				Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
				Result: &rawUncles[i],
			}
		}
		if err := c.Client.Client().BatchCallContext(ctx, reqs); err != nil {
			return nil, err
		}
		for i := range reqs {
			if reqs[i].Error != nil {
				return nil, reqs[i].Error
			}
			if rawUncles[i] == nil {
				return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
			}
			uncle, err := rawUncles[i].toHeader()
			if err != nil {
				return nil, err
			}
			uncles[i] = uncle
		}
	}

	// Fill transaction list.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if err := ensureTransactionHasRequiredSignature(tx.tx); err != nil {
			return nil, err
		}
		txs[i] = tx.tx
	}

	return types.NewBlockWithHeader(head).WithBody(
		types.Body{
			Transactions: txs,
			Uncles:       uncles,
			Withdrawals:  body.Withdrawals,
		}), nil
}

func ensureTransactionHasRequiredSignature(tx *types.Transaction) error {
	if tx == nil {
		return errors.New("server returned null transaction")
	}
	if tx.Type() == types.PoLTxType {
		return nil
	}
	_, r, _ := tx.RawSignatureValues()
	if r == nil {
		return errors.New("server returned transaction without signature")
	}
	return nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}
