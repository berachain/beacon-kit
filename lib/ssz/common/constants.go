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

package common

const (
	BitsPerByte          = 8
	BytesPerChunk        = 32
	BytesPerLengthOffset = 4
)

// Type is a SSZ type.
type Type int

const (
	// TypeUndefined is a sentinel zero value.
	TypeUndefined Type = iota
	// TypeUint is a SSZ int type, include byte.
	TypeUint
	// TypeBool is a SSZ bool type.
	TypeBool
	// TypeBytes is a SSZ fixed or dynamic bytes type.
	TypeBytes
	// TypeVector is a SSZ vector.
	TypeVector
	// TypeList is a SSZ list.
	TypeList
	// TypeContainer is a SSZ container.
	TypeContainer
)

func IsBasicType(t Type) bool {
	return t == TypeUint || t == TypeBool
}

func IsVariableSize(t Type) bool {
	switch t {
	case TypeList, TypeContainer:
		return true
	default:
		return false
	}
}
