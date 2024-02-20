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

package encoding

import (
	"cosmossdk.io/collections/codec"
	"github.com/ethereum/go-ethereum/common"
)

type Hash struct{}

// This is an assertion that Deposit implements the codec.ValueCodec interface.
var _ codec.ValueCodec[common.Hash] = Hash{}

func (Hash) Decode(b []byte) (common.Hash, error) {
	return common.BytesToHash(b), nil
}

// DecodeJSON implements codec.ValueCodec.
func (Hash) DecodeJSON(_ []byte) (common.Hash, error) {
	panic("unimplemented")
}

// Encode implements codec.ValueCodec.
func (Hash) Encode(value common.Hash) ([]byte, error) {
	return value[:], nil
}

// EncodeJSON implements codec.ValueCodec.
func (Hash) EncodeJSON(_ common.Hash) ([]byte, error) {
	panic("unimplemented")
}

// Stringify implements codec.ValueCodec.
func (Hash) Stringify(value common.Hash) string {
	return value.Hex()
}

// ValueType implements codec.ValueCodec.
func (Hash) ValueType() string {
	return "common.Hash"
}
