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

package codec

import (
	"cosmossdk.io/collections/codec"

	"github.com/itsdevbear/bolaris/lib/store/types"
)

type Uint256Key struct{}

var _ codec.KeyCodec[*types.Uint256] = Uint256Key{}

func (Uint256Key) Encode(buffer []byte, key *types.Uint256) (int, error) {
	key.WriteToSlice(buffer)
	return key.ByteLen(), nil
}

func (Uint256Key) Decode(buffer []byte) (int, *types.Uint256, error) {
	return len(buffer), new(types.Uint256).SetBytes(buffer), nil
}

func (Uint256Key) Size(key *types.Uint256) int {
	return key.ByteLen()
}

func (Uint256Key) EncodeJSON(value *types.Uint256) ([]byte, error) {
	return value.MarshalJSON()
}

func (Uint256Key) DecodeJSON(b []byte) (*types.Uint256, error) {
	key := new(types.Uint256)
	err := key.UnmarshalJSON(b)
	return key, err
}

func (Uint256Key) Stringify(key *types.Uint256) string {
	return key.String()
}

func (Uint256Key) KeyType() string {
	return "uint256"
}

func (k Uint256Key) EncodeNonTerminal(buffer []byte, key *types.Uint256) (int, error) {
	return k.Encode(buffer, key)
}

func (k Uint256Key) DecodeNonTerminal(buffer []byte) (int, *types.Uint256, error) {
	return k.Decode(buffer)
}

func (k Uint256Key) SizeNonTerminal(key *types.Uint256) int {
	return k.Size(key)
}
