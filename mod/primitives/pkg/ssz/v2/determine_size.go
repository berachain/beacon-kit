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

	"cosmossdk.io/errors"
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

func parseSSZFieldTags(field reflect.StructField) ([]uint64, bool, error) {
	tag, exists := field.Tag.Lookup("ssz-size")
	if !exists {
		return make([]uint64, 0), false, nil
	}
	items := strings.Split(tag, ",")
	if items == nil {
		return make([]uint64, 0), false, nil
	}
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
