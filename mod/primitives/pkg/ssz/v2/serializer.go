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

//nolint:ineffassign,wastedassign,mnd // experimental
package ssz

import (
	"encoding/binary"
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// SerializeError represents serialization errors.
type SerializeError struct {
	Err error
}

func (e *SerializeError) Error() string {
	return e.Err.Error()
}

// NewSerializeErrorMaximumLengthReached
// creates a new SerializeError for maximum length reached.
func NewSerializeErrorMaximumLengthReached(size int) *SerializeError {
	return &SerializeError{
		Err: errors.Newf(
			"the encoded length is %d which meets or exceeds the maximum length %d",
			size,
			MaximumLength,
		),
	}
}

// NewSerializeErrorInvalidInstance creates
// a new SerializeError for invalid instances.
func NewSerializeErrorInvalidInstance(err error) *SerializeError {
	return &SerializeError{Err: errors.Newf("invalid instance: %w", err)}
}

// NewSerializeErrorInvalidType creates a new SerializeError for invalid types.
func NewSerializeErrorInvalidType(err error) *SerializeError {
	return &SerializeError{Err: errors.Newf("invalid type: %w", err)}
}

type Serializer struct {
	ISerializer
}

type ISerializer interface {
	Elements(s GenericSSZType) []GenericSSZType
	MarshalSSZ(s GenericSSZType) ([]byte, error)
}

func NewSerializer() Serializer {
	return Serializer{}
}

// MarshalSSZ takes a SSZ value, reflects on the type, and returns a buffer. 0
// indexed, of the encoded value.
func (s *Serializer) MarshalSSZ(c interface{}) ([]byte, error) {
	typ := reflect.TypeOf(c)
	val := reflect.ValueOf(c)
	k := typ.Kind()

	switch {
	case IsUintLike(k):
		return RouteUint(val, typ), nil
	case k == reflect.Bool:
		return ssz.MarshalBool(c.(bool)), nil
	case k == reflect.Slice:
		// 1 dimensional array of uint8s or bytearray []byte.
		if typ.Elem().Kind() == reflect.Uint8 {
			return s.MarshalToDefaultBuffer(val, typ, s.MarshalByteArray)
		}
		// We follow fastssz generated code samples in
		// bellatrix.ssz.go for these.
		if isBasicType(typ.Elem().Kind()) {
			return s.MarshalNDimensionalArray(val)
		}
		if IsNDimensionalSliceLike(typ) {
			return s.MarshalNDimensionalArray(val)
		}
		// Todo: Variable size handling.
		// if isVariableSizeType(typ.Elem()) {
		// composite arr.
		// return s.MarshalToDefaultBuffer(val, typ, s.MarshalComposite)
		// }
		fallthrough
	case k == reflect.Array:
		// 1 dimensional array of uint8s or bytearray []byte.
		if typ.Elem().Kind() == reflect.Uint8 {
			return s.MarshalToDefaultBuffer(val, typ, s.MarshalByteArray)
		}
		// We follow fastssz generated code samples in
		// bellatrix.ssz.go for these.
		if isBasicType(typ.Elem().Kind()) {
			return s.MarshalNDimensionalArray(val)
		}
		if IsNDimensionalArrayLike(typ) {
			return s.MarshalNDimensionalArray(val)
		}
		// Todo: Variable size handling.
		// if isVariableSizeType(typ.Elem()) {
		// composite arr.
		// return s.MarshalToDefaultBuffer(val, typ, s.MarshalComposite)
		// }
		fallthrough
	// TODO(Chibera): fix me!
	// Composite structs appear initially as pointers so we Look inside
	// case k == reflect.Struct || reflect.TypeOf(val.Elem()).Kind() ==
	// reflect.Struct:
	// Composite struct
	// buf := make([]byte, 0)
	// _, err := s.MarshalStruct(val, typ, buf, 0)
	// if err != nil {
	// 	return nil, err
	// }
	// return buf, nil
	// case k == reflect.Ptr:
	// 	return make([]byte, 0), nil
	// Composite struct? Look inside?
	// return s.MarshalSSZ(val.Elem())
	default:
		return make(
				[]byte,
				0,
			), errors.Newf(
				"type %v is not serializable",
				val.Type(),
			)
	}
}

// Marshal is the top level fn. it returns a properly encoded byte buffer. given
// a pre-existing buf and typ.
func (s *Serializer) Marshal(
	val reflect.Value,
	typ reflect.Type,
	input []byte,
	startOffset uint64,
) (uint64, error) {
	marshalled, err := s.MarshalSSZ(val.Interface())
	if err != nil {
		return startOffset, err
	}
	var size uint64
	if isVariableSizeType(typ) {
		size = determineVariableSize(val, typ)
	} else {
		size = determineFixedSize(val, typ)
	}
	offset := startOffset + size
	//nolint:wastedassign // the underlying passed in input buffer is read
	input = append(input[startOffset:], marshalled...)
	return offset, err
}

func (s *Serializer) MarshalToDefaultBuffer(
	val reflect.Value,
	typ reflect.Type,
	cb func(reflect.Value, reflect.Type, []byte, uint64) (uint64, error),
) ([]byte, error) {
	aLen := val.Len()
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		aLen = GetNestedArrayLength(val)
	}
	buf := make([]byte, aLen)
	_, err := cb(val, typ, buf, 0)
	return buf, err
}

func (s *Serializer) MarshalNDimensionalArray(
	val reflect.Value,
) ([]byte, error) {
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return nil, errors.New("input is not an array or slice")
	}

	dimensionality := GetArrayDimensionality(val)
	if dimensionality == 0 {
		return nil, errors.New("zero-dimensional array provided")
	}

	// Calculate the total number of elements across all dimensions.
	totalElements := GetNestedArrayLength(val)
	if totalElements == 0 {
		return make(
			[]byte,
			0,
		), nil // Return an empty byte slice for an empty array
	}

	// Create a buffer to hold all byte values
	var buffer []byte

	// Recursive function to traverse and serialize elements
	var serializeRecursive func(reflect.Value) error
	serializeRecursive = func(currentVal reflect.Value) error {
		if currentVal.Kind() == reflect.Array ||
			currentVal.Kind() == reflect.Slice {
			for i := range currentVal.Len() {
				if err := serializeRecursive(currentVal.Index(i)); err != nil {
					return err
				}
			}
		} else {
			// Serialize single element
			bytes, err := s.MarshalSSZ(currentVal.Interface())
			if err != nil {
				return err
			}
			buffer = append(buffer, bytes...)
		}
		return nil
	}

	// Start the recursive serialization
	if err := serializeRecursive(val); err != nil {
		return nil, err
	}

	if len(buffer) > 0 {
		return buffer, nil
	}
	return nil, errors.Newf("got empty buffer in MarshalNDimensionalArray")
}

func (s *Serializer) MarshalByteArray(
	val reflect.Value,
	typ reflect.Type,
	buf []byte,
	startOffset uint64,
) (uint64, error) {
	if val.Kind() == reflect.Array {
		for i := range val.Len() {
			//#nosec:G701 // int overflow should be caught earlier in the stack.
			buf[int(startOffset)+i] = uint8(val.Index(i).Uint())
		}
		//#nosec:G701 // int overflow should be caught earlier in the stack.
		return startOffset + uint64(val.Len()), nil
	}
	if val.IsNil() {
		item := make([]byte, typ.Len())
		copy(buf[startOffset:], item)
		//#nosec:G701 // int overflow should be caught earlier in the stack.
		return startOffset + uint64(typ.Len()), nil
	}
	copy(buf[startOffset:], val.Bytes())

	//#nosec:G701 // int overflow should be caught earlier in the stack.
	return startOffset + uint64(val.Len()), nil
}

func (s *Serializer) UnmarshalByteArray(
	val reflect.Value,
	_ reflect.Type,
	input []byte,
	startOffset uint64,
) (uint64, error) {
	offset := startOffset + uint64(len(input))
	val.SetBytes(input[startOffset:offset])
	return offset, nil
}

func (s *Serializer) MarshalComposite(
	val reflect.Value,
	typ reflect.Type,
	buf []byte,
	startOffset uint64,
) (uint64, error) {
	index := startOffset
	//nolint:ineffassign,wastedassign // its fine. we reuse the err
	err := errors.Newf("failed to MarshalComposite from %v of typ %v", val, typ)
	if val.Len() == 0 {
		return index, nil
	}
	if !isVariableSizeType(typ.Elem()) {
		for i := range val.Len() {
			// If each element is not variable size, we simply encode
			// sequentially and write
			// into the buffer at the last index we wrote at.
			index, err = s.Marshal(val.Index(i), typ.Elem(), buf, index)
			if err != nil {
				return 0, err
			}
		}
		return index, nil
	}
	fixedIndex := index
	//#nosec:G701 // int overflow should be caught earlier in the stack
	currentOffsetIndex := startOffset + uint64(val.Len())*BytesPerLengthOffset
	//nolint:wastedassign // the underlying passed in input buffer is read
	nextOffsetIndex := currentOffsetIndex
	// If the elements are variable size, we need to include offset indices
	// in the serialized output list.
	for i := range val.Len() {
		nextOffsetIndex, err = s.Marshal(
			val.Index(i),
			typ.Elem(),
			buf,
			currentOffsetIndex,
		)
		if err != nil {
			return 0, err
		}
		// Write the offset.
		offsetBuf := make([]byte, BytesPerLengthOffset)
		//#nosec:G701 // int overflow should be caught earlier in the stack
		binary.LittleEndian.PutUint32(
			offsetBuf,
			uint32(currentOffsetIndex-startOffset),
		)
		copy(buf[fixedIndex:fixedIndex+BytesPerLengthOffset], offsetBuf)

		// We increase the offset indices accordingly.
		currentOffsetIndex = nextOffsetIndex
		fixedIndex += BytesPerLengthOffset
	}
	index = currentOffsetIndex
	return index, nil
}

// TODO
// func (s *Serializer) MarshalStruct(
// 	val reflect.Value,
// 	typ reflect.Type,
// 	buf []byte,
// 	startOffset uint64,
// ) (uint64, error) {
// 	if typ.Kind() == reflect.Ptr {
// 		if val.IsNil() {
// 			newVal := reflect.New(typ.Elem()).Elem()
// 			return s.Marshal(newVal, newVal.Type(), buf, startOffset)
// 		}
// 		return s.Marshal(val.Elem(), typ.Elem(), buf, startOffset)
// 	}
// 	fixedIndex := startOffset
// 	fixedLength := uint64(0)
// 	// For every field, we add up the total length of the items depending if
// 	// they
// 	// are variable or fixed-size fields.
// 	for i := range typ.NumField() {
// 		// We skip protobuf related metadata fields.
// 		if strings.Contains(typ.Field(i).Name, "XXX_") {
// 			continue
// 		}
// 		fType, err := determineFieldType(typ.Field(i))
// 		if err != nil {
// 			return 0, err
// 		}
// 		if isVariableSizeType(fType) {
// 			fixedLength += BytesPerLengthOffset
// 		} else {
// 			if val.Type().Kind() == reflect.Ptr && val.IsNil() {
// 				elem := reflect.New(val.Type().Elem()).Elem()
// 				fixedLength += determineFixedSize(elem, fType)
// 			} else {
// 				fixedLength += determineFixedSize(val.Field(i), fType)
// 			}
// 		}
// 	}
// 	//nolint:wastedassign // the underlying passed in input buffer is read
// 	currentOffsetIndex := startOffset + fixedLength
// 	//nolint:wastedassign // the underlying passed in input buffer is read
// 	nextOffsetIndex := currentOffsetIndex
// 	for i := range typ.NumField() {
// 		// We skip protobuf related metadata fields.
// 		if strings.Contains(typ.Field(i).Name, "XXX_") {
// 			continue
// 		}
// 		fType, err := determineFieldType(typ.Field(i))
// 		if err != nil {
// 			return 0, err
// 		}

// 		if !isVariableSizeType(fType) {
// 			fixedIndex, err = s.Marshal(val.Field(i), fType, buf, fixedIndex)
// 			if err != nil {
// 				return 0, err
// 			}
// 		} else {
// 			nextOffsetIndex, err = s.Marshal(
// 				val.Field(i), fType, buf, currentOffsetIndex)
// 			if err != nil {
// 				return 0, err
// 			}
// 			// Write the offset.
// 			offsetBuf := make([]byte, BytesPerLengthOffset)
// 			//#nosec:G701 // int overflow should be caught earlier in the stack
// 			binary.LittleEndian.PutUint32(offsetBuf,
// uint32(currentOffsetIndex-startOffset))
// 			copy(buf[fixedIndex:fixedIndex+BytesPerLengthOffset], offsetBuf)

// 			// We increase the offset indices accordingly.
// 			currentOffsetIndex = nextOffsetIndex
// 			fixedIndex += BytesPerLengthOffset
// 		}
// 	}
// 	return currentOffsetIndex, nil
// }
