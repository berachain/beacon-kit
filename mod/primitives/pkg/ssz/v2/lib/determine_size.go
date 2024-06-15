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
	switch kind {
	case reflect.Bool,
		reflect.Int32,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	default:
		return false
	}
}

func isBasicTypeSliceOrArr(t reflect.Type, k reflect.Kind) bool {
	return isBasicTypeSlice(t, k) || isBasicTypeArray(t, k)
}

func isBasicTypeSlice(typ reflect.Type, kind reflect.Kind) bool {
	return isBasicList(reflect.Slice, typ, kind)
}

func isBasicTypeArray(typ reflect.Type, kind reflect.Kind) bool {
	return isBasicList(reflect.Array, typ, kind)
}

func isBasicList(
	listType reflect.Kind,
	typ reflect.Type,
	kind reflect.Kind) bool {
	return kind == listType && isBasicType(typ.Elem().Kind())
}

func isVariableSizeType(typ reflect.Type) bool {
	kind := typ.Kind()
	switch {
	case isBasicType(kind),
		isBasicTypeArray(typ, kind):
		return false
	case kind == reflect.Slice,
		kind == reflect.String:
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

func isFixedSizeType(typ reflect.Type) bool {
	kind := typ.Kind()
	switch kind {
	case reflect.Bool,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Int32,
		reflect.Uint64:
		return true
	default:
		return false
	}
}

func determineFixedSize(val reflect.Value, typ reflect.Type) uint64 {
	kind := typ.Kind()
	switch kind {
	case reflect.Bool,
		reflect.Uint8:
		return 1
	case reflect.Uint16:
		//nolint:mnd // static mapped types
		return 2
	case reflect.Uint32, reflect.Int32:
		//nolint:mnd // static mapped types
		return 4
	case reflect.Uint64:
		//nolint:mnd // static mapped types
		return 8
	case reflect.Array, reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			//#nosec:G701 // will not realistically cause a problem.
			return uint64(val.Len())
		}
		var num uint64
		n := val.Len()
		for i := range n {
			num += determineFixedSize(val.Index(i), typ.Elem())
		}
		return num
	case reflect.Struct:
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
	case reflect.Ptr:
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
	switch kind {
	case reflect.String:
		//#nosec:G701 // will not realistically cause a problem.
		return uint64(val.Len())
	case reflect.Slice, reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			//#nosec:G701 // will not realistically cause a problem.
			return uint64(val.Len())
		}
		return determineSizeSliceOrArray(val, typ)
	case reflect.Struct:
		return determineSizeStruct(val, typ)
	case reflect.Ptr:
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
	if exists && len(fieldSizeTags) > 0 && !hasUndefinedSizeTag(field) {
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
		// If a field is unbounded, we skip it
		if items[i] != ssz.UnboundedSSZFieldSizeMarker {
			sizes[i], err = strconv.ParseUint(items[i], 10, 64)
			if err != nil {
				return make([]uint64, 0), false, err
			}
		}
	}
	return sizes, len(sizes) > 0, nil
}

func inferFieldTypeFromSizeTags(
	field reflect.StructField, sizes []uint64) reflect.Type {
	if isFixedSizeType(field.Type) {
		return field.Type
	}
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

// GetArrayDimensionality is function to determine the dimensionality
// of an N-dimensional array.
func GetArrayDimensionality(val reflect.Value) int {
	dimensionality := 0
	typ := reflect.TypeOf(val.Interface())
	kind := typ.Kind()

	if val.Len() == 0 && isBasicTypeSliceOrArr(typ, kind) {
		// 1 dimensional empty arr. e.g. Balances   []uint64
		return 1
	}

	for kind == reflect.Array || kind == reflect.Slice {
		dimensionality++
		if val.Len() == 0 {
			// Empty array, get the element type and update the kind
			typ = typ.Elem()
			kind = typ.Kind()
		} else {
			// Move to the next nested array
			val = val.Index(0)
			typ = val.Type()
			kind = typ.Kind()
		}
	}

	return dimensionality
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
