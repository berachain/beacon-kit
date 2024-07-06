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

package genesis

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// Genesis is a struct that contains the genesis information
// need to start the beacon chain.
//
//nolint:lll
type Genesis[
	DepositT any,
	ExecutionPayloadHeaderT interface {
		NewFromJSON([]byte, uint32) (ExecutionPayloadHeaderT, error)
	},
] struct {
	// ForkVersion is the fork version of the genesis slot.
	ForkVersion common.Version `json:"fork_version"`

	// Deposits represents the deposits in the genesis. Deposits are
	// used to initialize the validator set.
	Deposits []DepositT `json:"deposits"`

	// ExecutionPayloadHeader is the header of the execution payload
	// in the genesis.
	ExecutionPayloadHeader ExecutionPayloadHeaderT `json:"execution_payload_header"`
}

// GetForkVersion returns the fork version in the genesis.
func (g *Genesis[
	DepositT, ExecutionPayloadHeaderT,
]) GetForkVersion() common.Version {
	return g.ForkVersion
}

// GetDeposits returns the deposits in the genesis.
func (g *Genesis[DepositT, ExecutionPayloadHeaderT]) GetDeposits() []DepositT {
	return g.Deposits
}

// GetExecutionPayloadHeader returns the execution payload header.
func (g *Genesis[
	DepositT, ExecutionPayloadHeaderT,
]) GetExecutionPayloadHeader() ExecutionPayloadHeaderT {
	return g.ExecutionPayloadHeader
}

// UnmarshalJSON for Genesis.
func (g *Genesis[DepositT, ExecutionPayloadHeaderT]) UnmarshalJSON(
	data []byte,
) error {
	type genesisMarshalable[Deposit any] struct {
		ForkVersion            common.Version  `json:"fork_version"`
		Deposits               []DepositT      `json:"deposits"`
		ExecutionPayloadHeader json.RawMessage `json:"execution_payload_header"`
	}
	var g2 genesisMarshalable[DepositT]
	if err := json.Unmarshal(data, &g2); err != nil {
		return err
	}

	var (
		payloadHeader ExecutionPayloadHeaderT
		err           error
	)
	payloadHeader, err = payloadHeader.NewFromJSON(
		g2.ExecutionPayloadHeader,
		version.ToUint32(g2.ForkVersion),
	)
	if err != nil {
		return err
	}

	g.Deposits = g2.Deposits
	g.ForkVersion = g2.ForkVersion
	g.ExecutionPayloadHeader = payloadHeader
	return nil
}

// DefaultGenesisDeneb returns a the default genesis.
func DefaultGenesisDeneb() *Genesis[
	*types.Deposit, *types.ExecutionPayloadHeader,
] {
	defaultHeader, err :=
		DefaultGenesisExecutionPayloadHeaderDeneb()
	if err != nil {
		panic(err)
	}

	// TODO: Uncouple from deneb.
	return &Genesis[*types.Deposit, *types.ExecutionPayloadHeader]{
		ForkVersion: version.FromUint32[common.Version](
			version.Deneb,
		),
		Deposits: make([]*types.Deposit, 0),
		ExecutionPayloadHeader: &types.ExecutionPayloadHeader{
			InnerExecutionPayloadHeader: defaultHeader,
		},
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
		emptyTxsRoot         common.Root
		emptyWithdrawalsRoot common.Root
	)

	g.Go(func() error {
		var err error
		emptyTxsRoot, err = engineprimitives.Transactions{}.HashTreeRoot()
		return err
	})

	g.Go(func() error {
		var err error
		wds := ssz.ListFromElements(
			spec.DevnetChainSpec().MaxWithdrawalsPerPayload(),
			[]*engineprimitives.Withdrawal{}...,
		)
		emptyWithdrawalsRoot, err = wds.HashTreeRoot()
		return err
	})

	// If deriving either of the roots fails, return the error.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &types.ExecutionPayloadHeaderDeneb{
		ParentHash:   gethprimitives.ZeroHash,
		FeeRecipient: gethprimitives.ZeroAddress,
		StateRoot: common.Bytes32(gethprimitives.Hex2BytesFixed(
			"0x12965ab9cbe2d2203f61d23636eb7e998f167cb79d02e452f532535641e35bcc",
			constants.RootLength,
		)),
		ReceiptsRoot: common.Bytes32(gethprimitives.Hex2BytesFixed(
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			constants.RootLength,
		)),
		LogsBloom: make([]byte, constants.LogsBloomLength),
		Random:    common.Bytes32{},
		Number:    0,
		//nolint:mnd // default value.
		GasLimit:  math.U64(30000000),
		GasUsed:   0,
		Timestamp: 0,
		ExtraData: make([]byte, constants.ExtraDataLength),
		//nolint:mnd // default value.
		BaseFeePerGas: math.MustNewU256LFromBigInt(big.NewInt(3906250)),
		BlockHash: gethprimitives.HexToHash(
			"0xcfff92cd918a186029a847b59aca4f83d3941df5946b06bca8de0861fc5d0850",
		),
		TransactionsRoot: emptyTxsRoot,
		WithdrawalsRoot:  emptyWithdrawalsRoot,
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}, nil
}
