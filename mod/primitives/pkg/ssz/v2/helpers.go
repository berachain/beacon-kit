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
	"strconv"
	"strings"

	"github.com/berachain/beacon-kit/mod/errors"
	ssz "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// DetermineSize returns the required byte size of a buffer for
// using SSZ to marshal an object.
func DetermineSize(val reflect.Value) uint64 {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return DetermineSize(reflect.New(val.Type().Elem()).Elem())
		}
		return DetermineSize(val.Elem())
	}
	if isVariableSizeType(val.Type()) {
		return determineVariableSize(val, val.Type())
	}
	return determineFixedSize(val, val.Type())
}

func isBasicType(kind reflect.Kind) bool {
	return kind == reflect.Bool ||
		kind == reflect.Int32 ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

func isBasicTypeArray(typ reflect.Type, kind reflect.Kind) bool {
	return kind == reflect.Array && isBasicType(typ.Elem().Kind())
}

func isRootsArray(_ reflect.Value, typ reflect.Type) bool {
	elemTyp := typ.Elem()
	elemKind := elemTyp.Kind()
	arrCheck := elemKind == reflect.Array
	uintCheck := elemTyp.Elem().Kind() == reflect.Uint8
	isByteArray := arrCheck && uintCheck
	return isByteArray && elemTyp.Len() == 32
}

func isVariableSizeType(typ reflect.Type) bool {
	kind := typ.Kind()
	switch {
	case isBasicType(kind):
		return false
	case isBasicTypeArray(typ, kind):
		return false
	case kind == reflect.Slice:
		return true
	case kind == reflect.String:
		return true
	case kind == reflect.Array:
		return isVariableSizeType(typ.Elem())
	case kind == reflect.Struct:
		n := typ.NumField()
		for i := range n {
			if strings.Contains(typ.Field(i).Name, "XXX_") {
				continue
			}
			f := typ.Field(i)
			fType, err := determineFieldType(f)
			if err != nil {
				return false
			}
			if isVariableSizeType(fType) {
				return true
			}
		}
		return false
	case kind == reflect.Ptr:
		return isVariableSizeType(typ.Elem())
	}
	return false
}

func determineFixedSize(val reflect.Value, typ reflect.Type) uint64 {
	kind := typ.Kind()
	switch {
	case kind == reflect.Bool:

		return 1
	case kind == reflect.Uint8:

		return 1
	case kind == reflect.Uint16:
		//nolint:mnd // static mapped types
		return 2
	case kind == reflect.Uint32 || kind == reflect.Int32:
		//nolint:mnd // static mapped types
		return 4
	case kind == reflect.Uint64:
		//nolint:mnd // static mapped types
		return 8
	case kind == reflect.Array && typ.Elem().Kind() == reflect.Uint8:
		//#nosec:G701 // will not realistically cause a problem.
		return uint64(typ.Len())
	case kind == reflect.Slice && typ.Elem().Kind() == reflect.Uint8:
		//#nosec:G701 // will not realistically cause a problem.
		return uint64(val.Len())
	case kind == reflect.Array || kind == reflect.Slice:
		var num uint64
		n := val.Len()
		for i := range n {
			num += determineFixedSize(val.Index(i), typ.Elem())
		}
		return num
	case kind == reflect.Struct:
		totalSize := uint64(0)
		n := typ.NumField()
		for i := range n {
			if strings.Contains(typ.Field(i).Name, "XXX_") {
				continue
			}
			f := typ.Field(i)
			fType, err := determineFieldType(f)
			if err != nil {
				return 0
			}
			totalSize += determineFixedSize(val.Field(i), fType)
		}
		return totalSize
	case kind == reflect.Ptr:
		if val.IsNil() {
			newElem := reflect.New(typ.Elem()).Elem()
			return determineVariableSize(newElem, newElem.Type())
		}
		return determineFixedSize(val.Elem(), typ.Elem())
	default:
		return 0
	}
}

func determineVariableSize(val reflect.Value, typ reflect.Type) uint64 {
	kind := typ.Kind()
	switch {
	case kind == reflect.Slice && typ.Elem().Kind() == reflect.Uint8:
		//#nosec:G701 // will not realistically cause a problem.
		return uint64(val.Len())
	case kind == reflect.String:
		//#nosec:G701 // will not realistically cause a problem.
		return uint64(val.Len())
	case kind == reflect.Slice || kind == reflect.Array:
		return determineSizeSliceOrArray(val, typ)
	case kind == reflect.Struct:
		return determineSizeStruct(val, typ)
	case kind == reflect.Ptr:
		if val.IsNil() {
			newElem := reflect.New(typ.Elem()).Elem()
			return determineVariableSize(newElem, newElem.Type())
		}
		return determineVariableSize(val.Elem(), val.Elem().Type())
	default:
		return 0
	}
}

func determineSizeStruct(val reflect.Value, typ reflect.Type) uint64 {
	totalSize := uint64(0)
	for i := range typ.NumField() {
		if strings.Contains(typ.Field(i).Name, "XXX_") {
			continue
		}
		f := typ.Field(i)
		fType, err := determineFieldType(f)
		if err != nil {
			return 0
		}
		if isVariableSizeType(fType) {
			varSize := determineVariableSize(val.Field(i), fType)
			totalSize += varSize + ssz.BytesPerLengthOffset
		} else {
			varSize := determineFixedSize(val.Field(i), fType)
			totalSize += varSize
		}
	}
	return totalSize
}

func determineSizeSliceOrArray(val reflect.Value, typ reflect.Type) uint64 {
	totalSize := uint64(0)
	n := val.Len()
	for i := range n {
		varSize := DetermineSize(val.Index(i))
		if isVariableSizeType(typ.Elem()) {
			totalSize += varSize + ssz.BytesPerLengthOffset
		} else {
			totalSize += varSize
		}
	}
	return totalSize
}

func determineFieldType(field reflect.StructField) (reflect.Type, error) {
	fieldSizeTags, exists, err := parseSSZFieldTags(field)
	if err != nil {
		return reflect.TypeOf(reflect.Invalid), errors.Wrap(
			err,
			"could not parse ssz struct field tags")
	}
	if exists {
		// If the field does indeed specify ssz struct tags
		// we infer the field's type.
		return inferFieldTypeFromSizeTags(field, fieldSizeTags), nil
	}
	return field.Type, nil
}

func getSSZFieldTags(field reflect.StructField) []string {
	tag, exists := field.Tag.Lookup("ssz-size")
	if !exists {
		return make([]string, 0)
	}
	items := strings.Split(tag, ",")
	if items == nil {
		return make([]string, 0)
	}
	return items
}

func hasUndefinedSizeTag(field reflect.StructField) bool {
	items := getSSZFieldTags(field)
	for i := range items {
		// If a field is unbounded, we mark it as dynamic length otherwise we
		// treat it as fixed length
		if items[i] == ssz.UnboundedSSZFieldSizeMarker {
			return true
		}
	}
	//#nosec:G703 // idc about the strconv err.
	// its impossible if the previous step didnt detect a string and return
	// already.
	sizes, found, err := parseSSZFieldTags(field)
	if !found || err != nil {
		return true
	}
	sumUintArr := sumArr[[]uint64]
	return sumUintArr(sizes) <= 0
}

func parseSSZFieldTags(field reflect.StructField) ([]uint64, bool, error) {
	items := getSSZFieldTags(field)
	sizes := make([]uint64, len(items))
	var err error
	for i := range len(items) {
		// If a field is unbounded, we mark it with a size of 0.
		if items[i] == ssz.UnboundedSSZFieldSizeMarker {
			sizes[i] = 0
			continue
		}
		sizes[i], err = strconv.ParseUint(items[i], 10, 64)
		if err != nil {
			return make([]uint64, 0), false, err
		}
	}
	return sizes, true, nil
}

func inferFieldTypeFromSizeTags(
	field reflect.StructField, sizes []uint64) reflect.Type {
	innerElement := field.Type.Elem()
	for i := 1; i < len(sizes); i++ {
		innerElement = innerElement.Elem()
	}
	currentType := innerElement
	for i := len(sizes) - 1; i >= 0; i-- {
		if sizes[i] == 0 {
			currentType = reflect.SliceOf(currentType)
		} else {
			//#nosec:G701 // will not realistically cause a problem.
			currentType = reflect.ArrayOf(int(sizes[i]), currentType)
		}
	}
	return currentType
}

// Recursive function to calculate the length of an N-dimensional array.
func GetNestedArrayLength(val reflect.Value) int {
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		//#nosec:G701 // int overflow should be caught earlier in the stack.
		return int(DetermineSize(val))
	}
	length := val.Len()
	if length == 0 {
		return 0 // Early return for empty arrays/slices.
	}

	// Recursively calculate the length of the first element if it is an
	// array/slice.
	elementLength := GetNestedArrayLength(val.Index(0))
	return length * elementLength
}

// Function to determine the dimensionality of an N-dimensional array.
func GetArrayDimensionality(val reflect.Value) int {
	dimensionality := 0
	for val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		dimensionality++
		val = val.Index(0) // Move to the next nested array.
	}
	// for byte arrs
	return dimensionality
}

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
	isUintLike := false

	switch kind {
	case reflect.Uint8:
		isUintLike = true
	case reflect.Uint16:
		isUintLike = true
	case reflect.Uint32:
		isUintLike = true
	case reflect.Uint64:
		isUintLike = true
	default:
		return isUintLike
	}

	return isUintLike
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
	if typ.Kind() == reflect.Ptr {
		subtyp := reflect.TypeOf(val.Interface()).Elem()
		vf = reflect.VisibleFields(subtyp)
	}

	for i := range len(vf) {
		sf := vf[i]
		// Note: You can get the name this way for deserialization
		// name := sf.Name
		sft := sf.Type
		sfv := val.Elem().Field(i)
		cb(sft, sfv, sf, nil)
	}
}

// CalculateBufferSizeForStruct calculates the required buffer size for
// marshalling a struct using SSZ.
func CalculateBufferSizeForStruct(val reflect.Value) (int, error) {
	size := 0
	var errCheck []error
	IterStructFields(
		val,
		func(_ reflect.Type,
			val reflect.Value,
			_ reflect.StructField,
			err error,
		) {
			// #nosec:G701 // if we cant fit the size in an int. we cant fit the
			// value anywhere
			size += int(DetermineSize(val))
			if err != nil {
				// Track errors respective to what size was calculated when it
				// occurred.
				errCheck[size] = err
			}
		},
	)
	if len(errCheck) != 0 {
		return 0, errCheck[0]
	}
	return size, nil
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

	// Interleave offsets of variable-size parts with fixed-size parts.
	// variable_offsets = [serialize(uint32(sum(fixed_lengths +
	// variable_lengths[:i]))) for i in range(len(value))].
	offsetSum := sumIntArr(fixedLengths)
	variableOffsets := make([][]byte, len(variableParts))
	for i := range len(variableParts) {
		offsetSum += variableLengths[i]
		// #nosec:G701 // converting an int of max is 4294967295 to uint64 max
		// of 2147483647.
		// Wont realisticially overflow.
		variableOffsets[i] = ssz.MarshalU32(uint32(offsetSum))
	}

	for i, part := range fixedParts {
		if part == nil {
			fixedParts[i] = variableOffsets[i]
		}
	}

	// Flatten the nested arr to a 1d []byte
	allParts := make([][]byte, 0)
	allParts = append(allParts, variableParts...)
	allParts = append(allParts, fixedParts...)
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
	return typ.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct
}
