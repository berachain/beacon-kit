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

package staking

import (
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// Name of the Deposit event
	// in the deposit contract.
	DepositEventName = "Deposit"

	// Name the Withdrawal event
	// in the deposit contract.
	WithdrawalEventName = "Withdrawal"
)

//nolint:gochecknoglobals // Avoid re-allocating these variables.
var (
	// Signature and type of the Deposit event
	// in the deposit contract.
	DepositEventSig = crypto.Keccak256Hash(
		[]byte(DepositEventName + "(bytes,bytes,uint64,bytes,uint64)"),
	)

	// Signature and type of the Withdraw event
	// in the deposit contract.
	WithdrawalEventSig = crypto.Keccak256Hash(
		[]byte(WithdrawalEventName + "(bytes,bytes,bytes,uint64,uint64)"),
	)
)
