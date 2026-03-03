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

package ethclient_test

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	beraclient "github.com/berachain/beacon-kit/gethlib/ethclient"
	"github.com/berachain/beacon-kit/gethlib/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	gethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	testBlockHash    = gethcommon.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	testParentHash   = gethcommon.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	testStateRoot    = gethcommon.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	testReceiptsRoot = gethcommon.HexToHash("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	testTxRoot       = gethcommon.HexToHash("0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	testTxHash       = gethcommon.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")
	testFrom         = gethcommon.HexToAddress("0x1000000000000000000000000000000000000001")
	testTo           = gethcommon.HexToAddress("0x2000000000000000000000000000000000000002")
)

func TestBlockByNumberWithPoLTransaction(t *testing.T) {
	t.Parallel()

	client := newMockClient(t)

	block, err := client.BlockByNumber(t.Context(), big.NewInt(2))
	if err != nil {
		t.Fatalf("BlockByNumber returned error: %v", err)
	}
	if block == nil {
		t.Fatal("BlockByNumber returned nil block")
	}
	if got := len(block.Transactions()); got != 1 {
		t.Fatalf("unexpected tx count: got %d want %d", got, 1)
	}
	if txType := block.Transactions()[0].Type(); txType != types.PoLTxType {
		t.Fatalf("unexpected tx type: got 0x%x want 0x%x", txType, types.PoLTxType)
	}
}

func TestBlockByHashWithPoLTransaction(t *testing.T) {
	t.Parallel()

	client := newMockClient(t)

	block, err := client.BlockByHash(t.Context(), testBlockHash)
	if err != nil {
		t.Fatalf("BlockByHash returned error: %v", err)
	}
	if block == nil {
		t.Fatal("BlockByHash returned nil block")
	}
	if got := len(block.Transactions()); got != 1 {
		t.Fatalf("unexpected tx count: got %d want %d", got, 1)
	}
	if txType := block.Transactions()[0].Type(); txType != types.PoLTxType {
		t.Fatalf("unexpected tx type: got 0x%x want 0x%x", txType, types.PoLTxType)
	}
}

func TestTransactionByHashWithPoLTransaction(t *testing.T) {
	t.Parallel()

	client := newMockClient(t)

	tx, isPending, err := client.TransactionByHash(t.Context(), testTxHash)
	if err != nil {
		t.Fatalf("TransactionByHash returned error: %v", err)
	}
	if tx == nil {
		t.Fatal("TransactionByHash returned nil transaction")
	}
	if isPending {
		t.Fatal("TransactionByHash unexpectedly marked tx as pending")
	}
	if txType := tx.Type(); txType != types.PoLTxType {
		t.Fatalf("unexpected tx type: got 0x%x want 0x%x", txType, types.PoLTxType)
	}
}

func TestTransactionInBlockWithPoLTransaction(t *testing.T) {
	t.Parallel()

	client := newMockClient(t)

	tx, err := client.TransactionInBlock(t.Context(), testBlockHash, 0)
	if err != nil {
		t.Fatalf("TransactionInBlock returned error: %v", err)
	}
	if tx == nil {
		t.Fatal("TransactionInBlock returned nil transaction")
	}
	if txType := tx.Type(); txType != types.PoLTxType {
		t.Fatalf("unexpected tx type: got 0x%x want 0x%x", txType, types.PoLTxType)
	}
}

type mockEthService struct{}

func (s *mockEthService) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (json.RawMessage, error) {
	_ = ctx
	_ = number
	_ = fullTx
	return mockBlockJSON(), nil
}

func (s *mockEthService) GetBlockByHash(ctx context.Context, hash gethcommon.Hash, fullTx bool) (json.RawMessage, error) {
	_ = ctx
	_ = hash
	_ = fullTx
	return mockBlockJSON(), nil
}

func (s *mockEthService) GetTransactionByHash(ctx context.Context, hash gethcommon.Hash) (json.RawMessage, error) {
	_ = ctx
	_ = hash
	return mockTransactionJSON(), nil
}

func (s *mockEthService) GetTransactionByBlockHashAndIndex(
	ctx context.Context,
	blockHash gethcommon.Hash,
	index string,
) (json.RawMessage, error) {
	_ = ctx
	_ = blockHash
	_ = index
	return mockTransactionJSON(), nil
}

func newMockClient(t *testing.T) *beraclient.Client {
	t.Helper()

	srv := rpc.NewServer()
	if err := srv.RegisterName("eth", &mockEthService{}); err != nil {
		t.Fatalf("failed to register mock service: %v", err)
	}
	rpcClient := rpc.DialInProc(srv)
	client := beraclient.Wrap(gethclient.NewClient(rpcClient))
	t.Cleanup(func() {
		client.Close()
		srv.Stop()
	})
	return client
}

func mockTransactionJSON() json.RawMessage {
	tx := mockTransactionObject()
	raw, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return raw
}

func mockTransactionObject() map[string]any {
	return map[string]any{
		"type":             "0x7e",
		"chainId":          "0x1",
		"from":             testFrom.Hex(),
		"to":               testTo.Hex(),
		"nonce":            "0x2",
		"gas":              "0x5208",
		"gasPrice":         "0x1",
		"input":            "0x1234",
		"hash":             testTxHash.Hex(),
		"blockHash":        testBlockHash.Hex(),
		"blockNumber":      "0x2",
		"transactionIndex": "0x0",
	}
}

func mockBlockJSON() json.RawMessage {
	block := map[string]any{
		"hash":             testBlockHash.Hex(),
		"parentHash":       testParentHash.Hex(),
		"sha3Uncles":       coretypes.EmptyUncleHash.Hex(),
		"miner":            gethcommon.Address{}.Hex(),
		"stateRoot":        testStateRoot.Hex(),
		"transactionsRoot": testTxRoot.Hex(),
		"receiptsRoot":     testReceiptsRoot.Hex(),
		"logsBloom":        "0x" + strings.Repeat("0", 512),
		"difficulty":       "0x1",
		"number":           "0x2",
		"gasLimit":         "0x1c9c380",
		"gasUsed":          "0x5208",
		"timestamp":        "0x1",
		"extraData":        "0x",
		"mixHash":          gethcommon.Hash{}.Hex(),
		"nonce":            "0x0000000000000000",
		"transactions":     []any{mockTransactionObject()},
		"uncles":           []any{},
	}
	raw, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return raw
}
