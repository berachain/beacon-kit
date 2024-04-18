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

package deneb

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultBeaconState returns a default BeaconState.
//
// TODO: take in BeaconConfig params to determine the
// default length of the arrays, which we are currently
// and INCORRECTLY setting to 0.
func DefaultBeaconState() *BeaconState {
	//nolint:mnd // default allocs.
	return &BeaconState{
		GenesisValidatorsRoot: primitives.Root{},
		Slot:                  0,
		Fork: &primitives.Fork{
			PreviousVersion: version.FromUint32[primitives.Version](
				version.Deneb,
			),
			CurrentVersion: version.FromUint32[primitives.Version](
				version.Deneb,
			),
			Epoch: 0,
		},
		LatestBlockHeader: &primitives.BeaconBlockHeader{
			Slot:          0,
			ProposerIndex: 0,
			ParentRoot:    primitives.Root{},
			StateRoot:     primitives.Root{},
			BodyRoot:      primitives.Root{},
		},
		BlockRoots:             make([]primitives.Root, 8),
		StateRoots:             make([]primitives.Root, 8),
		LatestExecutionPayload: DefaultGenesisExecutionPayload(),
		Eth1Data: &primitives.Eth1Data{
			DepositRoot:  primitives.Root{},
			DepositCount: 0,
			BlockHash:    primitives.ExecutionHash{},
		},
		Eth1DepositIndex:             0,
		Validators:                   make([]*primitives.Validator, 0),
		Balances:                     make([]uint64, 0),
		NextWithdrawalIndex:          0,
		NextWithdrawalValidatorIndex: 0,
		RandaoMixes:                  make([]primitives.Bytes32, 8),
		Slashings:                    make([]uint64, 0),
		TotalSlashing:                0,
	}
}

// DefaultGenesisExecutionPayload returns a default ExecutableDataDeneb.
//
//nolint:mnd // default values pulled from current eth-genesis.json file.
func DefaultGenesisExecutionPayload() *engineprimitives.ExecutableDataDeneb {
	return &engineprimitives.ExecutableDataDeneb{
		ParentHash:   primitives.ExecutionHash{},
		FeeRecipient: primitives.ExecutionAddress{},
		StateRoot: common.HexToHash(
			"0x12965ab9cbe2d2203f61d23636eb7e998f167cb79d02e452f532535641e35bcc",
		),
		ReceiptsRoot: common.HexToHash(
			"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		),
		LogsBloom: make([]byte, 256),
		Random:    primitives.ExecutionHash{},
		Number:    0,
		GasLimit:  math.U64(hexutil.MustDecodeUint64("0x1c9c380")),
		GasUsed:   0,
		Timestamp: 0,
		ExtraData: make([]byte, 32),
		BaseFeePerGas: math.NewU256LFromBigEndian(
			hexutil.MustDecode("0x3b9aca"),
		),
		BlockHash: common.HexToHash(
			"0xcfff92cd918a186029a847b59aca4f83d3941df5946b06bca8de0861fc5d0850",
		),
		Transactions:  [][]byte{},
		Withdrawals:   []*engineprimitives.Withdrawal{},
		BlobGasUsed:   0,
		ExcessBlobGas: 0,
	}
}

//go:generate go run github.com/ferranbt/fastssz/sszgen -path deneb.go -objs BeaconState -include ../../types,../../../primitives-engine,../../../primitives,../../../primitives/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deneb.ssz.go
//nolint:lll // various json tags.
type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot primitives.Root  `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  math.Slot        `json:"slot"`
	Fork                  *primitives.Fork `json:"fork"`

	// History
	LatestBlockHeader *primitives.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        []primitives.Root             `json:"blockRoots"        ssz-size:"?,32" ssz-max:"8192"`
	StateRoots        []primitives.Root             `json:"stateRoots"        ssz-size:"?,32" ssz-max:"8192"`

	// Eth1
	LatestExecutionPayload *engineprimitives.ExecutableDataDeneb `json:"latestExecutionPayload"`
	Eth1Data               *primitives.Eth1Data                  `json:"eth1Data"`
	Eth1DepositIndex       uint64                                `json:"eth1DepositIndex"`

	// Registry
	Validators []*primitives.Validator `json:"validators" ssz-max:"1099511627776"`
	Balances   []uint64                `json:"balances"   ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes []primitives.Bytes32 `json:"randaoMixes" ssz-size:"?,32" ssz-max:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64  `json:"slashings"     ssz-max:"1099511627776"`
	TotalSlashing math.Gwei `json:"totalSlashing"`
}
