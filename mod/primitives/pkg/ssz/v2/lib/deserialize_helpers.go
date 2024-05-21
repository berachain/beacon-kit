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

package ssz

import (
	"reflect"

	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

func RouteUintUnmarshal(k reflect.Kind, buf []byte) reflect.Value {
	switch k {
	case reflect.Uint8:
		return reflect.ValueOf(ssz.UnmarshalU8[uint8](buf))
	case reflect.Uint16:
		return reflect.ValueOf(ssz.UnmarshalU16[uint16](buf))
	case reflect.Uint32:
		return reflect.ValueOf(ssz.UnmarshalU32[uint32](buf))
	case reflect.Uint64:
		// handle native
		// if data, ok := val.Interface().([]byte); ok {
		// 	u64Val := ssz.UnmarshalU64(data)
		// 	return ssz.MarshalU64(u64Val)
		// }
		return reflect.ValueOf(ssz.UnmarshalU64[uint64](buf))

	// TODO(Chibera): Handle numbers over 64bit?
	// case reflect.Uint128:
	// 	return UnmarshalU128(val.Interface().([]byte))
	// case reflect.Uint256:
	// 	return UnmarshalU256(val.Interface().([]byte))
	default:
		return reflect.ValueOf(make([]byte, 0))
	}
}
