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

package enginetypes

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Withdrawal represents a validator withdrawal from the consensus layer.
//
//go:generate go run github.com/fjl/gencodec -type Withdrawal -field-override withdrawalJSONMarshaling -out withdrawal.json.go
type Withdrawal struct {
	Index     uint64                      `json:"index"          ssz-size:"8"`
	Validator primitives.ValidatorIndex   `json:"validatorIndex" ssz-size:"8"`
	Address   primitives.ExecutionAddress `json:"address"        ssz-size:"20"`
	Amount    uint64                      `json:"amount"         ssz-size:"8"`
}

// String returns a string representation of the Withdrawal.
func (w *Withdrawal) String() string {
	//#nosec:G703 // ignore potential marshalling failure.
	output, _ := json.Marshal(w)
	return string(output)
}

// field type overrides for gencodec.
type withdrawalJSONMarshaling struct {
	Index     hexutil.Uint64
	Validator hexutil.Uint64
	Amount    hexutil.Uint64
}
