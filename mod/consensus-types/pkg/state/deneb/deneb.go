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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

//go:generate go run github.com/ferranbt/fastssz/sszgen -path deneb.go -objs BeaconState -include ../../../../primitives/pkg/crypto,../../../../primitives/pkg/common,../../../../primitives/pkg/bytes,../../../../primitives/mod.go,../../../../consensus-types/pkg/types,../../../../primitives-engine,../../../../primitives/mod.go,../../../../primitives/pkg/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deneb.ssz.go
//nolint:lll // various json tags.
type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot primitives.Root `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  math.Slot       `json:"slot"`
	Fork                  *types.Fork     `json:"fork"`

	// History
	LatestBlockHeader *types.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        []primitives.Root        `json:"blockRoots"        ssz-size:"?,32" ssz-max:"8192"`
	StateRoots        []primitives.Root        `json:"stateRoots"        ssz-size:"?,32" ssz-max:"8192"`

	// Eth1
	Eth1Data                     *types.Eth1Data                    `json:"eth1Data"`
	Eth1DepositIndex             uint64                             `json:"eth1DepositIndex"`
	LatestExecutionPayloadHeader *types.ExecutionPayloadHeaderDeneb `json:"latestExecutionPayloadHeader"`

	// Registry
	Validators []*types.Validator `json:"validators" ssz-max:"1099511627776"`
	Balances   []uint64           `json:"balances"   ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes []primitives.Bytes32 `json:"randaoMixes" ssz-size:"?,32" ssz-max:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64  `json:"slashings"     ssz-max:"1099511627776"`
	TotalSlashing math.Gwei `json:"totalSlashing"`
}
