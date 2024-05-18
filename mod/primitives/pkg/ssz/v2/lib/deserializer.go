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
)

/***
Deserialization
ideas:
accept an interface, e.g. Type BeaconBlockDeneb via deneb.BeaconState{}
accept marshalled buffer
recursively understand first the interface object
	recursively look at each field val
		for each fixed item run determineSize() and make a map
		for each variable sized item, add to map with size 0
		(if any of the fixed parts eval'd to nil and we have a variable offset at that loc? panic?)

for each item in our fixed/simple types to size map
	slice the buffer till item size, pass the rest on
	set interface val to unmarshalled data

we now have remaining variable parts
for each variable field in map
	call readOffset and div by BYTES_PER_LENGTH_OFFSET to get size (at the start?)
	read as fixed size item with the now known size

*/

// Future home of our deserializer
func (s *Serializer) UnmarshalSSZ(data []byte, val interface{}) error {
	// Check if the provided value is a pointer
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}

	// Create a new reflect.Value from the provided value
	valReflect := reflect.ValueOf(val).Elem()

	// Unmarshal the data into the provided struct using reflection
	if err := s.unmarshalSSZRecursive(data, valReflect); err != nil {
		return err
	}

	return nil
}

// func RouteUintUnmarshal(val reflect.Value, buf []byte) reflect.Value {
// 	kind := val.Kind()
// 	switch kind {
// 	case reflect.Uint8:
// 		return reflect.ValueOf(ssz.UnmarshalU8[uint8](val.Interface().([]byte)))
// 	case reflect.Uint16:
// 		return reflect.ValueOf(ssz.UnmarshalU16[uint16](val.Interface().([]byte)))
// 	case reflect.Uint32:
// 		return reflect.ValueOf(ssz.UnmarshalU32[uint32](val.Interface().([]byte)))
// 	case reflect.Uint64:
// 		// handle native
// 		// if data, ok := val.Interface().([]byte); ok {
// 		// 	u64Val := ssz.UnmarshalU64(data)
// 		// 	return ssz.MarshalU64(u64Val)
// 		// }
// 		return reflect.ValueOf(ssz.UnmarshalU64[uint64](val.Interface().([]byte)))

// 	// TODO(Chibera): Handle numbers over 64bit?
// 	// case reflect.Uint128:
// 	// 	return UnmarshalU128(val.Interface().([]byte))
// 	// case reflect.Uint256:
// 	// 	return UnmarshalU256(val.Interface().([]byte))
// 	default:
// 		return reflect.ValueOf(make([]byte, 0))
// 	}
// }

func (s *Serializer) unmarshalSSZRecursive(data []byte, val reflect.Value) error {
	size := DetermineSize(val)
	k := val.Kind()
	buf := data[:size]
	elem := reflect.New(val.Type().Elem())
	typ := val.Type()

	if IsUintLike(k) {
		val.Set(RouteUintUnmarshal(k, buf))
	}
	switch val.Kind() {
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			// return s.UnmarshalByteArray(data, val, typ)
		}
		if isBasicType(typ.Elem().Kind()) {
			return s.UnmarshalNDimensionalArray(data, val)
		}
		if isVariableSizeType(typ) || isVariableSizeType(typ.Elem()) {
			return s.UnmarshalComposite(data, val)
		}
		if IsNDimensionalSliceLike(typ) {
			return s.UnmarshalNDimensionalArray(data, val)
		}
		val.Set(elem)
		fallthrough
	case reflect.Ptr, reflect.Struct:
		if k == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		}
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if err := s.unmarshalSSZRecursive(data, field); err != nil {
				return err
			}
		}
	case reflect.Slice:
		// Unmarshal the length of the slice
		length := binary.LittleEndian.Uint64(data[:8])
		data = data[8:length]

		// Create a new slice with the specified length
		slice := reflect.MakeSlice(val.Type(), int(length), int(length))

		// Unmarshal each element of the slice
		for i := 0; i < int(length); i++ {
			if err := s.unmarshalSSZRecursive(data, slice.Index(i)); err != nil {
				return err
			}
		}

		// Set the slice value to the provided struct field
		val.Set(slice)
	case reflect.Uint64:
		// Unmarshal the uint64 value
		val.SetUint(binary.LittleEndian.Uint64(data[:8]))
	// case reflect.Array:
	// 	// Unmarshal each element of the array
	// 	for i := 0; i < val.Len(); i++ {
	// 		if err := s.unmarshalSSZRecursive(data, val.Index(i)); err != nil {
	// 			return err
	// 		}
	// 	}
	default:
		return errors.New("unsupported type")
	}

	return nil
}

// func (s *Serializer) UnmarshalByteArray(
// 	val reflect.Value,
// 	_ reflect.Type,
// 	input []byte,
// 	startOffset uint64,
// ) (uint64, error) {
// 	offset := startOffset + uint64(len(input))
// 	val.SetBytes(input[startOffset:offset])
// 	return offset, nil
// }

// UnmarshalSSZ takes a byte buffer, reflects on the provided interface object or struct, and unmarshals the buffer into it.
func (s *Serializer) UnmarshalSSZ2(data []byte, c interface{}) error {
	val := reflect.ValueOf(c)
	typ := reflect.TypeOf(c)

	if IsUintLike(typ.Kind()) {
		return errors.New("unmarshaling into uint-like types is not supported")
	}

	switch typ.Kind() {
	case reflect.Bool:
		return errors.New("unmarshaling into bool types is not supported")
	case reflect.Slice, reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			// return s.UnmarshalByteArray(data, val, typ)
		}
		if isBasicType(typ.Elem().Kind()) {
			return s.UnmarshalNDimensionalArray(data, val)
		}
		if isVariableSizeType(typ) || isVariableSizeType(typ.Elem()) {
			return s.UnmarshalComposite(data, val)
		}
		if IsNDimensionalSliceLike(typ) {
			return s.UnmarshalNDimensionalArray(data, val)
		}
		fallthrough
	case reflect.Ptr, reflect.Struct:
		if typ.Kind() == reflect.Struct || typ.Elem().Kind() == reflect.Struct {
			return s.UnmarshalStruct(data, val, typ)
		}
		fallthrough
	default:
		return errors.Newf("type %v is not deserializable", val.Type())
	}
}

func (s *Serializer) UnmarshalNDimensionalArray(data []byte, val reflect.Value) error {
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return errors.New("input is not an array or slice")
	}

	dimensionality := GetArrayDimensionality(val)
	if dimensionality == 0 {
		return errors.New("zero-dimensional array provided")
	}

	// Deserialize each element in the array
	var offset uint64
	var errCheck []error
	var processMember = func(c interface{}) {
		err := s.UnmarshalSSZ(data[offset:], c)
		if err != nil {
			errCheck = append(errCheck, err)
		}
		offset += uint64(len(data))
	}

	// Start the recursive deserialization
	if err := SerializeRecursive(val, processMember); err != nil {
		return err
	}

	if len(errCheck) > 0 {
		return errCheck[0]
	}
	return nil
}

func (s *Serializer) UnmarshalStruct(data []byte, val reflect.Value, typ reflect.Type) error {
	if !IsStruct(typ, val) {
		return errors.New("input is not a struct")
	}

	var errCheck []error

	var processStructField = func(typ reflect.Type, val reflect.Value, field reflect.StructField, err error) {
		if err != nil {
			errCheck = append(errCheck, err)
			return
		}

		if hasUndefinedSizeTag(field) && isVariableSizeType(typ) {
			// Handle variable size fields
			// Not implemented in this example
		} else {
			// Handle fixed size fields
			// Not implemented in this example
		}
	}

	IterStructFields(val, processStructField)

	if len(errCheck) > 0 {
		return errCheck[0]
	}

	return nil
}

func (s *Serializer) UnmarshalComposite(data []byte, val reflect.Value) error {
	var errCheck []error

	var processMember = func(c interface{}) {
		memberTyp := reflect.TypeOf(c)
		// memberVal := reflect.ValueOf(c)

		if isVariableSizeType(memberTyp) {
			// Handle variable size fields
			// Not implemented in this example
		} else {
			// Handle fixed size fields
			// Not implemented in this example
		}
	}

	// Start the recursive deserialization
	for i := range val.Len() {
		processMember(val.Index(i).Interface())
	}

	if len(errCheck) > 0 {
		return errCheck[0]
	}

	return nil
}
