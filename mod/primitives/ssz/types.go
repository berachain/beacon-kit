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

// Basic defines an interface for SSZ basic types which includes methods for
// determining the size of the SSZ encoding and computing the hash tree root.
type Basic[SpecT any, RootT ~[32]byte] interface {
	// SizeSSZ returns the size in bytes of the SSZ-encoded data.
	SizeSSZ() int
	// HashTreeRoot computes and returns the hash tree root of the data as RootT
	// and an error if the computation fails.
	HashTreeRoot( /*...args*/ ) (RootT, error)
}

// Composite is an interface that embeds the Basic interface. It is used for
// types that are composed of other SSZ encodable values.
type Composite[SpecT any, RootT ~[32]byte] interface {
	Basic[SpecT, RootT]
	Elements() []Value
}

// Container is an interface for SSZ container types that can be marshaled and
// unmarshaled.
type Container[SpecT any, RootT ~[32]byte] interface {
	Composite[SpecT, RootT]
	FieldTypes() []Type
	Kind() Kind
}
type (
	// Kind represents different types of SSZ'able data.
	Kind uint8

	// Type represents a SSZ type.
	Type interface {
		Kind() Kind
	}

	// BasicType represents a basic SSZ type.
	BasicType interface {
		Type
		// Size returns the length, in bytes,
		// of the serialized form of the basic type.
		SizeOf() int
	}

	// ArrayType represents a SSZ vector or list type.
	ArrayType interface {
		Type
		// ElemType returns the type of the array elements.
		ElemType() Type
		// Size returns the size of the composite type.
		Size() int
	}

	ContainerType interface {
		Type
		// FieldTypes returns the types of the container fields.
		FieldTypes() []Type
	}
)

const (
	// KindUndefined is a sentinel zero value.
	KindUndefined Kind = iota
	// KindUInt is a SSZ int type, include byte.
	KindUInt
	// KindBool is a SSZ bool type.
	KindBool
	// KindBytes is a SSZ fixed or dynamic bytes type.
	KindBytes
	// KindVector is a SSZ vector.
	KindVector
	// KindList is a SSZ list.
	KindList
	// KindContainer is a SSZ container.
	KindContainer
)

// Vector represents the SSZ vector type.
type Vector struct {
	size     int
	elemType Type
}

// NewVector creates a new vector type.
func NewVector(size int, elemType Type) Vector {
	return Vector{
		size:     size,
		elemType: elemType,
	}
}

// Kind returns the type kind.
func (t Vector) Kind() Kind {
	return KindVector
}

// Size returns the fixed size of the vector.
func (t Vector) Size() int {
	return t.size
}

// ElemType returns the type of the vector elements.
func (t Vector) ElemType() Type {
	return t.elemType
}

// Value is a common interface for SSZ values.
type Value interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Type() Type
}
