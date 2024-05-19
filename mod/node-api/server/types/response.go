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

package types

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type ErrorResponse struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type GenesisData struct {
	GenesisTime           string             `json:"genesis_time"`
	GenesisValidatorsRoot primitives.Bytes32 `json:"genesis_validators_root"`
	GenesisForkVersion    string             `json:"genesis_fork_version"`
}

type RootData struct {
	Root primitives.Root `json:"root"`
}

type ValidatorStateResponse struct {
	ExecutionOptimistic bool         `json:"execution_optimistic"`
	Finalized           bool         `json:"finalized"`
	Data                any `json:"data"`
}


type ValidatorData struct {
	Index     math.U64         `json:"index"`
	Balance   math.U64         `json:"balance"`
	Status    string           `json:"status"`
	Validator *types.Validator `json:"validator"`
}


type ValidatorBalanceData struct {
	Index   math.U64 `json:"index"`
	Balance math.U64 `json:"balance"`
}

type CommitteeData struct {
	Index      math.U64   `json:"index"`
	Slot       math.U64   `json:"slot"`
	Validators []math.U64 `json:"validators"`
}
