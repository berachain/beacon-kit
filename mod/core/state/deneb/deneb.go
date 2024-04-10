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
	"github.com/berachain/beacon-kit/mod/config/version"
	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultBeaconState returns a default BeaconState.
//
// TODO: take in BeaconConfig params to determine the
// default length of the arrays, which we are currently
// and INCORRECTLY setting to 0.
func DefaultBeaconState() *BeaconState {
	//nolint:gomnd // default allocs.
	return &BeaconState{
		GenesisValidatorsRoot: primitives.Root{},
		Slot:                  0,
		Fork: &primitives.Fork{
			PreviousVersion: version.FromUint32(version.Deneb),
			CurrentVersion:  version.FromUint32(version.Deneb),
			Epoch:           0,
		},
		LatestBlockHeader: &primitives.BeaconBlockHeader{
			Slot:          0,
			ProposerIndex: 0,
			ParentRoot:    primitives.Root{},
			StateRoot:     primitives.Root{},
			BodyRoot:      primitives.Root{},
		},
		BlockRoots: make([][32]byte, 8),
		StateRoots: make([][32]byte, 8),
		Eth1BlockHash: common.HexToHash(
			"0xcfff92cd918a186029a847b59aca4f83d3941df5946b06bca8de0861fc5d0850",
		),
		Eth1Data: &primitives.Eth1Data{
			DepositRoot:  primitives.Root{},
			DepositCount: 0,
			BlockHash:    primitives.ExecutionHash{},
		},
		Eth1DepositIndex:             0,
		Validators:                   make([]*types.Validator, 0),
		Balances:                     make([]uint64, 0),
		NextWithdrawalIndex:          0,
		NextWithdrawalValidatorIndex: 0,
		RandaoMixes:                  make([][32]byte, 8),
		Slashings:                    make([]uint64, 0),
		TotalSlashing:                0,
	}
}

// TODO: should we replace ? in ssz-size with values to ensure we are hash tree
// rooting correctly?
//
//go:generate go run github.com/fjl/gencodec -type BeaconState -field-override BeaconStateJSONMarshaling -out deneb.json.go
//go:generate go run github.com/ferranbt/fastssz/sszgen -path deneb.go -objs BeaconState -include ../../types,../../../primitives,../../../execution/types,$GETH_PKG_INCLUDE/common -output deneb.ssz.go
//nolint:lll // various json tags.
type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot primitives.Root  `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  primitives.Slot  `json:"slot"`
	Fork                  *primitives.Fork `json:"fork"`

	// History
	LatestBlockHeader *primitives.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        [][32]byte                    `json:"blockRoots"        ssz-size:"?,32" ssz-max:"8192"`
	StateRoots        [][32]byte                    `json:"stateRoots"        ssz-size:"?,32" ssz-max:"8192"`

	// Eth1
	Eth1BlockHash    primitives.ExecutionHash `json:"eth1BlockHash" ssz-size:"32"`
	Eth1Data         *primitives.Eth1Data     `json:"eth1Data"`
	Eth1DepositIndex uint64                   `json:"eth1DepositIndex"`

	// Registry
	Validators []*types.Validator `json:"validators" ssz-max:"1099511627776"`
	Balances   []uint64           `json:"balances"   ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes [][32]byte `json:"randaoMixes" ssz-size:"?,32" ssz-max:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64                    `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex primitives.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64        `json:"slashings"     ssz-max:"1099511627776"`
	TotalSlashing primitives.Gwei `json:"totalSlashing"`
}

// String returns a string representation of BeaconState.
func (b *BeaconState) String() string {
	return spew.Sdump(b)
}

// BeaconStateJSONMarshaling is a type used to marshal/unmarshal
// BeaconState.
type BeaconStateJSONMarshaling struct {
	GenesisValidatorsRoot hexutil.Bytes
	BlockRoots            []primitives.Root
	StateRoots            []primitives.Root
	RandaoMixes           []primitives.Root
}
