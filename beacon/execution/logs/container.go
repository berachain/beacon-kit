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

package logs

import (
	"reflect"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

var _ LogValueContainer = (*Container)(nil)

type Container struct {
	value       *reflect.Value
	index       uint64
	sig         ethcommon.Hash
	blockNumber uint64
	blockHash   ethcommon.Hash
}

// BlockHash returns the block hash of the log.
func (c *Container) BlockHash() ethcommon.Hash {
	return c.blockHash
}

// BlockNumber returns the block number of the log.
func (c *Container) BlockNumber() uint64 {
	return c.blockNumber
}

// LogIndex returns the index of the log.
func (c *Container) LogIndex() uint64 {
	return c.index
}

// Value returns the value of the log.
func (c *Container) Value() *reflect.Value {
	return c.value
}

// Signature returns the signature of the log.
func (c *Container) Signature() ethcommon.Hash {
	return c.sig
}
