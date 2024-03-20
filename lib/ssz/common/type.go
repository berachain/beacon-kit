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

// Type is a SSZ type.
type Type interface {
	Kind() Kind
}

type BasicType interface {
	Type
	// Size returns the length, in bytes,
	// of the serialized form of the basic type.
	SizeOf() int
}

type TypeUint struct {
	// size in bits
	Size int
}

func (t TypeUint) Kind() Kind {
	return KindUint
}

func (t TypeUint) SizeOf() int {
	return t.Size / BitsPerByte
}

type TypeBool struct{}

func (t TypeBool) Kind() Kind {
	return KindBool
}

func (t TypeBool) SizeOf() int {
	return 1
}

type TypeVector struct {
	Size     int
	ElemType Type
}

func (t TypeVector) Kind() Kind {
	return KindVector
}

type TypeList struct {
	MaxSize  int
	ElemType Type
}

func (t TypeList) Kind() Kind {
	return KindList
}

type TypeContainer struct {
	FieldTypes []Type
}

func (t TypeContainer) Kind() Kind {
	return KindContainer
}

func IsBasicType(t Type) bool {
	k := t.Kind()
	return k == KindUint || k == KindBool
}

// IsVariableSize returns true if the object is variable-size.
// A variable-size types to be lists, (unsupported unions, Bitlist)
// and all types that contain a variable-size type.
// All other types are said to be fixed-size.
func IsVariableSize(t Type) bool {
	switch t.Kind() {
	case KindUint, KindBool:
		return false
	case KindList:
		return true
	case KindVector:
		return IsVariableSize(t.(TypeVector).ElemType)
	case KindContainer:
		for _, ft := range t.(TypeContainer).FieldTypes {
			if IsVariableSize(ft) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func IsList(t Type) bool {
	return t.Kind() == KindList
}
