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

package types

import (
	"fmt"
	"math/big"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

const (
	defaultGasLimit      = math.U64(30000000)
	defaultBaseFeePerGas = int64(3906250)
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
	*Deposit, *ExecutionPayloadHeader,
] {
	defaultHeader, err :=
		DefaultGenesisExecutionPayloadHeaderDeneb()
	if err != nil {
		panic(err)
	}

	// TODO: Uncouple from deneb.
	return &Genesis[*Deposit, *ExecutionPayloadHeader]{
		ForkVersion: version.FromUint32[common.Version](
			version.Deneb,
		),
		Deposits:               make([]*Deposit, 0),
		ExecutionPayloadHeader: defaultHeader,
	}
}

// DefaultGenesisExecutionPayloadHeaderDeneb returns a default
// ExecutionPayloadHeaderDeneb.
func DefaultGenesisExecutionPayloadHeaderDeneb() (
	*ExecutionPayloadHeader, error,
) {
	stateRoot, err := byteslib.ToBytes32(
		hex.ToBytesSafe(
			"0x12965ab9cbe2d2203f61d23636eb7e998f167cb79d02e452f532535641e35bcc",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed generating state root: %w", err)
	}

	receiptsRoot, err := byteslib.ToBytes32(
		hex.ToBytesSafe(
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed generating receipts root: %w", err)
	}

	baseFeePerGas, err := math.NewU256FromBigInt(
		big.NewInt(defaultBaseFeePerGas),
	)
	if err != nil {
		return nil, fmt.Errorf("failed setting base fee per gas: %w", err)
	}

	return &ExecutionPayloadHeader{
		ParentHash:    common.ExecutionHash{},
		FeeRecipient:  common.ExecutionAddress{},
		StateRoot:     stateRoot,
		ReceiptsRoot:  receiptsRoot,
		LogsBloom:     [256]byte{},
		Random:        common.Bytes32{},
		Number:        0,
		GasLimit:      defaultGasLimit,
		GasUsed:       0,
		Timestamp:     0,
		ExtraData:     make([]byte, constants.ExtraDataLength),
		BaseFeePerGas: baseFeePerGas,
		BlockHash: common.NewExecutionHashFromHex(
			"0xcfff92cd918a186029a847b59aca4f83d3941df5946b06bca8de0861fc5d0850",
		),
		TransactionsRoot: engineprimitives.Transactions(nil).
			HashTreeRoot(),
		WithdrawalsRoot: engineprimitives.Withdrawals(nil).HashTreeRoot(),
		BlobGasUsed:     0,
		ExcessBlobGas:   0,
	}, nil
}
