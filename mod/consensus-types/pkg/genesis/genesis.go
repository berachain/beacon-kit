// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package genesis

import (
	"context"
	"math/big"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// Genesis is a struct that contains the genesis information
// need to start the beacon chain.
type Genesis[
	DepositT any,
	ExecutonPayloadHeaderT engineprimitives.ExecutionPayloadHeader,
] struct {
	// ForkVersion is the fork version of the genesis slot.
	ForkVersion primitives.Version `json:"fork_version"`

	// Deposits represents the deposits in the genesis. Deposits are
	// used to initialize the validator set.
	Deposits []DepositT `json:"deposits"`

	// ExecutionPayloadHeader is the header of the execution payload
	// in the genesis.
	ExecutionPayloadHeader ExecutonPayloadHeaderT `json:"execution_payload_header"`
}

// DefaultGenesis returns a the default genesis.
func DefaultGenesisDeneb() *Genesis[
	*types.Deposit, *types.ExecutionPayloadHeaderDeneb,
] {
	defaultHeader, err :=
		DefaultGenesisExecutionPayloadHeaderDeneb()
	if err != nil {
		panic(err)
	}

	return &Genesis[*types.Deposit, *types.ExecutionPayloadHeaderDeneb]{
		ForkVersion: version.FromUint32[primitives.Version](
			version.Deneb,
		),
		Deposits:               make([]*types.Deposit, 0),
		ExecutionPayloadHeader: defaultHeader,
	}
}

// DefaultGenesisExecutionPayloadHeaderDeneb returns a default
// ExecutionPayloadHeaderDeneb.
func DefaultGenesisExecutionPayloadHeaderDeneb() (
	*types.ExecutionPayloadHeaderDeneb, error,
) {
	// Get the merkle roots of empty transactions and withdrawals in parallel.
	var (
		g, _                 = errgroup.WithContext(context.Background())
		emptyTxsRoot         primitives.Root
		emptyWithdrawalsRoot primitives.Root
	)

	g.Go(func() error {
		var err error
		emptyTxsRoot, err = engineprimitives.Transactions{}.HashTreeRoot()
		return err
	})

	g.Go(func() error {
		var err error
		emptyWithdrawalsRoot, err = engineprimitives.Withdrawals{}.HashTreeRoot()
		return err
	})

	// If deriving either of the roots fails, return the error.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &types.ExecutionPayloadHeaderDeneb{
		ParentHash:   common.ZeroHash,
		FeeRecipient: common.ZeroAddress,
		StateRoot: primitives.Bytes32(common.Hex2BytesFixed(
			"0x12965ab9cbe2d2203f61d23636eb7e998f167cb79d02e452f532535641e35bcc",
			constants.RootLength,
		)),
		ReceiptsRoot: primitives.Bytes32(common.Hex2BytesFixed(
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			constants.RootLength,
		)),
		LogsBloom: make([]byte, constants.LogsBloomLength),
		Random:    primitives.Bytes32{},
		Number:    0,
		//nolint:mnd // default value.
		GasLimit:  math.U64(30000000),
		GasUsed:   0,
		Timestamp: 0,
		ExtraData: make([]byte, constants.ExtraDataLength),
		//nolint:mnd // default value.
		BaseFeePerGas: math.MustNewU256LFromBigInt(big.NewInt(3906250)),
		BlockHash: common.HexToHash(
			"0xcfff92cd918a186029a847b59aca4f83d3941df5946b06bca8de0861fc5d0850",
		),
		TransactionsRoot: emptyTxsRoot,
		WithdrawalsRoot:  emptyWithdrawalsRoot,
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}, nil
}
