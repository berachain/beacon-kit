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

package consensus

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

//go:generate go run github.com/ferranbt/fastssz/sszgen -path eth1data.go -objs Eth1Data -include ../common,../../primitives.go,../bytes,$GETH_PKG_INCLUDE/common -output eth1data.ssz.go
type Eth1Data struct {
	// DepositRoot is the root of the deposit tree.
	DepositRoot primitives.Root `json:"depositRoot"  ssz-size:"32"`
	// DepositCount is the number of deposits in the deposit tree.
	DepositCount uint64 `json:"depositCount"`
	// BlockHash is the hash of the block corresponding to the Eth1Data.
	BlockHash common.ExecutionHash `json:"blockHash"    ssz-size:"32"`
}
