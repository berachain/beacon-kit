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

package types

import (
	"fmt"
	"math/big"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	"github.com/berachain/beacon-kit/primitives/math"
)

const (
	defaultGasLimit      = math.U64(30000000)
	defaultBaseFeePerGas = int64(3906250)
)

// Genesis is a struct that contains the genesis information
// need to start the beacon chain.
type Genesis struct {
	// ForkVersion is the fork version of the genesis slot.
	ForkVersion common.Version `json:"fork_version"`

	// Deposits represents the deposits in the genesis. Deposits are
	// used to initialize the validator set.
	Deposits []*Deposit `json:"deposits"`

	// ExecutionPayloadHeader is the header of the execution payload
	// in the genesis.
	ExecutionPayloadHeader *ExecutionPayloadHeader `json:"execution_payload_header"`
}

// GetForkVersion returns the fork version in the genesis.
func (g *Genesis) GetForkVersion() common.Version {
	return g.ForkVersion
}

// GetDeposits returns the deposits in the genesis.
func (g *Genesis) GetDeposits() []*Deposit {
	return g.Deposits
}

// GetExecutionPayloadHeader returns the execution payload header.
func (g *Genesis) GetExecutionPayloadHeader() *ExecutionPayloadHeader {
	return g.ExecutionPayloadHeader
}

// UnmarshalJSON for Genesis.
func (g *Genesis) UnmarshalJSON(
	data []byte,
) error {
	type genesisMarshalable[Deposit any] struct {
		ForkVersion            common.Version  `json:"fork_version"`
		Deposits               []*Deposit      `json:"deposits"`
		ExecutionPayloadHeader json.RawMessage `json:"execution_payload_header"`
	}
	var g2 genesisMarshalable[Deposit]
	if err := json.Unmarshal(data, &g2); err != nil {
		return err
	}

	payloadHeader := NewEmptyExecutionPayloadHeaderWithVersion(g2.ForkVersion)
	if err := json.Unmarshal(g2.ExecutionPayloadHeader, payloadHeader); err != nil {
		return err
	}

	g.Deposits = g2.Deposits
	g.ForkVersion = g2.ForkVersion
	g.ExecutionPayloadHeader = payloadHeader
	return nil
}

// DefaultGenesis returns the default genesis.
func DefaultGenesis(v common.Version) *Genesis {
	defaultHeader, err := DefaultGenesisExecutionPayloadHeader(v)
	if err != nil {
		panic(err)
	}

	return &Genesis{
		ForkVersion:            v,
		Deposits:               make([]*Deposit, 0),
		ExecutionPayloadHeader: defaultHeader,
	}
}

// DefaultGenesisExecutionPayloadHeader returns a default ExecutionPayloadHeader.
func DefaultGenesisExecutionPayloadHeader(v common.Version) (*ExecutionPayloadHeader, error) {
	stateRoot, err := bytes.ToBytes32(
		hex.MustToBytes(
			"0x12965ab9cbe2d2203f61d23636eb7e998f167cb79d02e452f532535641e35bcc",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed generating state root: %w", err)
	}

	receiptsRoot, err := bytes.ToBytes32(
		hex.MustToBytes(
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
		Versionable:   NewVersionable(v),
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
