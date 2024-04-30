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
	"fmt"
	"reflect"

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
	//nolint:lll
	return &SerializeError{Err: fmt.Errorf("the encoded length is %d which meets or exceeds the maximum length %d", size, MaximumLength)}
}

// NewSerializeErrorInvalidInstance creates
// a new SerializeError for invalid instances.
func NewSerializeErrorInvalidInstance(err error) *SerializeError {
	return &SerializeError{Err: fmt.Errorf("invalid instance: %w", err)}
}

// NewSerializeErrorInvalidType creates a new SerializeError for invalid types.
func NewSerializeErrorInvalidType(err error) *SerializeError {
	return &SerializeError{Err: fmt.Errorf("invalid type: %w", err)}
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
	// TODO(Chibera): Handle numbers over 64bit
	// case reflect.Uint128:
	// 	return MarshalU128(val.Interface().(uint128))
	// case reflect.Uint256:
	// 	return MarshalU256(val.Interface().(uint256))
	default:
		return make([]byte, 0)
	}
}

func IsUintLike(typ reflect.Type) bool {
	kind := typ.Kind()
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

func (s *Serializer) MarshalSSZ(c interface{}) ([]byte, error) {
	typ := reflect.TypeOf(c)
	k := typ.Kind()
	isUintLike := IsUintLike(typ)

	if isUintLike {
		return RouteUint(reflect.ValueOf(c), reflect.TypeOf(c)), nil
	}
	switch k {
	case reflect.Bool:
		return ssz.MarshalBool(reflect.ValueOf(c).Interface().(bool)), nil
	// TODO(Chibera): handle composite types. same algo for all 3
	// case KindList:
	// 	return true
	// case KindVector:
	// 	return IsVariableSize(t.(VectorType).ElemType())
	// case KindContainer:
	// 	for _, ft := range t.(ContainerType).FieldTypes() {
	// 		if IsVariableSize(ft) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	default:
		return make([]byte, 0), nil
	}
}
