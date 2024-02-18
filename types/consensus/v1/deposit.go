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

package consensusv1

import (
	"cosmossdk.io/collections/codec"
)

var _ codec.ValueCodec[*Deposit] = &Deposit{}

func (*Deposit) Encode(value *Deposit) ([]byte, error) {
	return value.MarshalSSZ()
}

func (*Deposit) Decode(b []byte) (*Deposit, error) {
	var v = &Deposit{}
	if err := v.UnmarshalSSZ(b); err != nil {
		return nil, err
	}
	return v, nil
}

func (*Deposit) EncodeJSON(_ *Deposit) ([]byte, error) {
	panic("not implemented")
}

func (*Deposit) DecodeJSON(_ []byte) (*Deposit, error) {
	panic("not implemented")
}

func (*Deposit) Stringify(value *Deposit) string {
	return value.String()
}

func (*Deposit) ValueType() string {
	return "Deposit"
}
