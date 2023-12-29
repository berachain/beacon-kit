// SPDX-License-Identifier: Apache-2.0
//
// Copyright (c) 2023 Berachain Foundation
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

package callback

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// validateArg uses reflection to verify the implementation arg matches the ABI arg.
func validateArg(implMethodVar reflect.Value, abiMethodVar reflect.Value) error {
	implMethodVarType := implMethodVar.Type()
	abiMethodVarType := abiMethodVar.Type()

	switch implMethodVarType.Kind() { //nolint:exhaustive // todo verify its okay.
	case reflect.TypeOf(common.Hash{}).Kind():
	case reflect.String:
		return validateString(implMethodVarType, abiMethodVarType)
	case reflect.Array, reflect.Slice:
		return validateArrayOrSlice(implMethodVarType, abiMethodVarType)
	case abiMethodVarType.Kind():
		return validateSameKind(implMethodVarType, abiMethodVarType)
	case reflect.Ptr:
		return validatePointer(implMethodVarType, abiMethodVarType)
	case reflect.Interface:
		// If it's `any` (reflect.Interface), we leave it to the implementer to make sure that it is
		// used/converted correctly.
	default:
		return fmt.Errorf("type mismatch: %v != %v", implMethodVarType, abiMethodVarType)
	}

	return nil
}

// validateString checks if the implementation type is a string, the ABI type must be a string.
func validateString(implMethodVarType reflect.Type, abiMethodVarType reflect.Type) error {
	if abiMethodVarType.Kind() != reflect.String {
		return fmt.Errorf(
			"type mismatch: %v != %v", implMethodVarType, abiMethodVarType,
		)
	}
	return nil
}

// validateArrayOrSlice checks if the array is not a slice/array of structs, return an error.
// If it is a slice/array of structs, check if the struct fields match.
func validateArrayOrSlice(implMethodVarType reflect.Type, abiMethodVarType reflect.Type) error {
	if implMethodVarType.Elem() != abiMethodVarType.Elem() {
		if implMethodVarType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf(
				"type mismatch: %v != %v", implMethodVarType, abiMethodVarType,
			)
		}

		if err := validateStruct(implMethodVarType.Elem(), abiMethodVarType.Elem()); err != nil {
			return err
		}
	}
	return nil
}

// validateSameKind checks if it's a struct, check all the fields to match.
// If the types (primitives) match, we're good.
func validateSameKind(implMethodVarType reflect.Type, abiMethodVarType reflect.Type) error {
	if implMethodVarType.Kind() == reflect.Struct {
		if err := validateStruct(implMethodVarType, abiMethodVarType); err != nil {
			return err
		}
	}
	return nil
}

// validatePointer checks if the corresponding ABI type must be a struct.
// Any implementation type that is a pointer must point to a struct.
// Check if the struct fields match.
func validatePointer(implMethodVarType reflect.Type, abiMethodVarType reflect.Type) error {
	if abiMethodVarType.Kind() != reflect.Struct {
		return fmt.Errorf(
			"type mismatch: %v != %v", implMethodVarType, abiMethodVarType,
		)
	}

	if implMethodVarType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf(
			"type mismatch: %v != %v", implMethodVarType, abiMethodVarType,
		)
	}

	return validateStruct(implMethodVarType.Elem(), abiMethodVarType)
}

// validateStruct checks to make sure that the implementation struct's fields match the ABI
// struct's fields.
func validateStruct(implMethodVarType reflect.Type, abiMethodVarType reflect.Type) error {
	if implMethodVarType.Kind() != reflect.Struct || abiMethodVarType.Kind() != reflect.Struct {
		return errors.New("validateStruct: not a struct")
	}

	if implMethodVarType.NumField() != abiMethodVarType.NumField() {
		return fmt.Errorf(
			"struct %v has %v fields, but struct %v has %v fields",
			implMethodVarType.Name(),
			implMethodVarType.NumField(),
			abiMethodVarType.Name(),
			abiMethodVarType.NumField(),
		)
	}

	// match every individual field
	for j := 0; j < implMethodVarType.NumField(); j++ {
		if err := validateArg(
			reflect.New(implMethodVarType.Field(j).Type).Elem(),
			reflect.New(abiMethodVarType.Field(j).Type).Elem(),
		); err != nil {
			return err
		}
	}
	return nil
}
