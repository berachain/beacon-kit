// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package ssz

import (
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

func IsCompositeType(t reflect.Type) bool {
	// array is fixed length and analogous to vector
	// slice is variable and analogous to list
	// Vectors, containers, lists, unions are considered composite types
	// Since we pre-handle Arrays and slices we return false for now
	// We only trigger on containers
	return t.Kind() == reflect.Struct
}

func IsNDimensionalArrayLike(typ reflect.Type) bool {
	ct := reflect.Array
	// A N dimensional array has a top level type of array and elem type also of
	// arr.
	return typ.Kind() == ct && typ.Elem().Kind() == ct
}

func IsNDimensionalSliceLike(typ reflect.Type) bool {
	ct := reflect.Slice
	// A N dimensional array has a top level type of Slice and elem type also of
	// Slice.
	return typ.Kind() == ct && typ.Elem().Kind() == ct
}

func RouteUint(val reflect.Value, typ reflect.Type) []byte {
	kind := typ.Kind()
	switch kind {
	case reflect.Uint8:
		return ssz.MarshalU8(val.Interface().(uint8))
	case reflect.Uint16:
		return ssz.MarshalU16(val.Interface().(uint16))
	case reflect.Uint32:
		return ssz.MarshalU32(val.Interface().(uint32))
	case reflect.Uint64:
		// handle native
		if data, ok := val.Interface().(math.U64); ok {
			serialized, serializationErr := data.MarshalSSZ()
			if serializationErr != nil {
				panic(serializationErr)
			}
			return serialized
		}
		return ssz.MarshalU64(val.Interface().(uint64))
	// TODO(Chibera): Handle numbers over 64bit?
	// case reflect.Uint128:
	// 	return MarshalU128(val.Interface().(uint128))
	// case reflect.Uint256:
	// 	return MarshalU256(val.Interface().(uint256))
	default:
		return make([]byte, 0)
	}
}

func IsUintLike(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	default:
		return false
	}
}

// Helper to iterate over fields in a struct.
func IterStructFields(
	val reflect.Value,
	cb func(
		typ reflect.Type,
		val reflect.Value,
		field reflect.StructField,
		err error,
	),
) {
	typ := reflect.TypeOf(val.Interface())
	vf := make([]reflect.StructField, 0)
	if !IsStruct(typ, val) {
		// Kick back incoming types in case of err for debugging upstream in the
		// caller fn.
		cb(
			typ,
			val,
			vf[0],
			errors.Newf("wrong data type provided to IterStructFields"),
		)
		return
	}

	// Deref the pointer
	subtyp := typ
	numFields := 0
	subval := val
	if typ.Kind() == reflect.Ptr {
		subtyp = reflect.TypeOf(val.Interface()).Elem()
		subval = val.Elem()
		numFields = subval.NumField()
	}
	if typ.Kind() == reflect.Struct {
		numFields = val.NumField()
	}

	vf = reflect.VisibleFields(subtyp)
	// Double check field count for rare nested cases
	iterLen := len(vf)
	if numFields < len(vf) {
		iterLen = numFields
	}

	for i := range iterLen {
		sf := vf[i]
		// Note: You can get the name this way for deserialization
		// name := sf.Name
		sft := sf.Type
		sfv := subval.Field(i)
		cb(sft, sfv, sf, nil)
	}
}

// Recursive function to traverse and serialize elements in slice or arr.
func SerializeRecursive(currentVal reflect.Value, cb func(interface{})) error {
	if currentVal.Kind() == reflect.Array ||
		currentVal.Kind() == reflect.Slice {
		for i := range currentVal.Len() {
			if err := SerializeRecursive(currentVal.Index(i), cb); err != nil {
				return err
			}
		}
	} else {
		// Serialize single element
		cb(currentVal.Interface())
	}
	return nil
}

func InterleaveOffsets(
	fixedParts [][]byte,
	fixedLengths []int,
	variableParts [][]byte,
	variableLengths []int,
) ([]byte, error) {
	sumIntArr := sumArr[[]int]
	// Check lengths
	totalLength := sumIntArr(fixedLengths) + sumIntArr(variableLengths)
	if totalLength >= 1<<(BytesPerLengthOffset*BitsPerByte) {
		return nil, errors.New("total length exceeds allowable limit")
	}

	if len(variableLengths) != len(variableParts) ||
		len(variableParts) < len(fixedParts) {
		return nil, errors.New(
			"variableParts & variableLengths must be same length",
		)
	}

	// Compute the offsets for the variable-size parts by summing up the fixed
	// lengthsand the variable lengths up to each index. The offsets are then
	// serialized as uint32 values and stored in the variableOffsets slice.
	offsetSum := sumIntArr(fixedLengths)
	variableOffsets := make([][]byte, len(variableParts))

	for i := range len(variableParts) {
		// #nosec:G701 // converting an int of max is 4294967295 to uint64 max
		// of 2147483647.
		// Wont realisticially overflow.
		variableOffsets[i] = ssz.MarshalU32(uint32(offsetSum))

		// Increment offsetSum by the length of the current variable part.
		offsetSum += variableLengths[i]
	}

	fixedPartsWithOffsets := make([][]byte, len(fixedParts))
	for i, part := range fixedParts {
		if part == nil {
			fixedPartsWithOffsets[i] = variableOffsets[i]
		} else {
			fixedPartsWithOffsets[i] = part
		}
	}

	// Flatten the nested arr to a 1d []byte
	allParts := make([][]byte, 0)
	allParts = append(allParts, fixedPartsWithOffsets...)
	allParts = append(allParts, variableParts...)
	res := make([]byte, 0)
	for i := range allParts {
		res = append(res, allParts[i]...)
	}

	return res, nil
}

func sumArr[S ~[]E, E ~int | ~uint | ~float64 | ~uint64](s S) E {
	var total E
	for _, v := range s {
		total += v
	}
	return total
}

func IsStruct(typ reflect.Type, val reflect.Value) bool {
	return typ.Kind() == reflect.Struct ||
		(typ.Kind() == reflect.Ptr &&
			val.Elem().Kind() == reflect.Struct)
}

func SafeCopyBuffer(res []byte, buf *[]byte, startOffset uint64) {
	bufLocal := *buf
	if len(res) != len(bufLocal) {
		//#nosec:G701 // will not realistically cause a problem.
		buf2 := make([]byte, len(res)+int(startOffset))
		copy(buf2, bufLocal[:startOffset])
		copy(buf2[startOffset:], res)
		*buf = buf2
		return
	}
	copy(bufLocal[startOffset:], res)
	*buf = bufLocal
}
