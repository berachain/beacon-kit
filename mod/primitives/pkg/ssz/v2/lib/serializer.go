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

//nolint:wastedassign // experimental
package ssz

import (
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
	Marshal(
		val reflect.Value,
		typ reflect.Type,
		input []byte,
		startOffset uint64,
	) (uint64, error)
	MarshalSSZ(c interface{}) ([]byte, error)
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

	if IsUintLike(k) {
		return RouteUint(val, typ), nil
	}

	switch k {
	case reflect.Bool:
		return ssz.MarshalBool(c.(bool)), nil
	case reflect.Slice, reflect.Array:
		// 1 dimensional array of uint8s or bytearray []byte.
		if typ.Elem().Kind() == reflect.Uint8 {
			return s.MarshalToDefaultBuffer(val, typ, s.MarshalByteArray)
		}
		// We follow fastssz generated code samples in
		// bellatrix.ssz.go for these.
		if isBasicType(typ.Elem().Kind()) {
			return s.MarshalNDimensionalArray(val)
		}
		if isVariableSizeType(typ) || isVariableSizeType(typ.Elem()) {
			// composite arr.
			return s.MarshalToDefaultBuffer(val, typ, s.MarshalComposite)
		}
		if IsNDimensionalSliceLike(typ) {
			return s.MarshalNDimensionalArray(val)
		}
		fallthrough
	case reflect.Ptr, reflect.Struct:
		// Composite structs appear initially as pointers so we Look inside
		if k == reflect.Struct || typ.Elem().Kind() == reflect.Struct {
			return s.MarshalToDefaultBuffer(val, typ, s.MarshalStruct)
		}
		fallthrough
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
	//nolint:wastedassign,staticcheck // the underlying passed in input buffer
	// is read
	input = append(input[startOffset:], marshalled...)
	return offset, err
}

func (s *Serializer) MarshalToDefaultBuffer(
	val reflect.Value,
	typ reflect.Type,
	cb func(reflect.Value, reflect.Type, *[]byte, uint64) (uint64, error),
) ([]byte, error) {
	aLen := 0
	//nolint:ineffassign,staticcheck // wtf lol.
	err := errors.New("MarshalToDefaultBuffer Failure")
	switch {
	case IsStruct(typ, val):
		aLen, err = CalculateBufferSizeForStruct(val)
		if err != nil {
			return nil, err
		}
	case val.Kind() == reflect.Array || val.Kind() == reflect.Slice:
		aLen = GetNestedArrayLength(val)
	default:
		aLen = val.Len()
	}
	buf := make([]byte, aLen)
	_, err = cb(val, typ, &buf, 0)
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
	var errCheck []error
	var processMember = func(c interface{}) {
		// Serialize single element
		bytes, err := s.MarshalSSZ(c)
		if err != nil {
			errCheck = append(errCheck, err)
		}
		buffer = append(buffer, bytes...)
	}

	// Start the recursive serialization
	if err := SerializeRecursive(val, processMember); err != nil {
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
	buf *[]byte,
	startOffset uint64,
) (uint64, error) {
	bufLocal := *buf
	if val.Kind() == reflect.Array {
		for i := range val.Len() {
			//#nosec:G701 // int overflow should be caught earlier in the stack.
			bufLocal[int(startOffset)+i] = uint8(val.Index(i).Uint())
		}
		*buf = bufLocal
		//#nosec:G701 // int overflow should be caught earlier in the stack.
		return startOffset + uint64(val.Len()), nil
	}
	if val.IsNil() {
		item := make([]byte, typ.Len())
		copy(bufLocal[startOffset:], item)
		*buf = bufLocal
		//#nosec:G701 // int overflow should be caught earlier in the stack.
		return startOffset + uint64(typ.Len()), nil
	}
	copy(bufLocal[startOffset:], val.Bytes())
	*buf = bufLocal
	//#nosec:G701 // int overflow should be caught earlier in the stack.
	return startOffset + uint64(val.Len()), nil
}

func (s *Serializer) MarshalFixedSizeParts(
	val reflect.Value,
	fixedParts [][]byte,
	fixedLengths []int,
) ([][]byte, []int, error) {
	serialized, err := s.MarshalSSZ(val.Interface())
	if err != nil {
		return fixedParts, fixedLengths, err
	}
	fixedParts = append(fixedParts, serialized)
	partSize := BytesPerLengthOffset
	if len(serialized) > 0 {
		partSize = len(serialized)
	}
	fixedLengths = append(fixedLengths, partSize)
	return fixedParts, fixedLengths, nil
}

func (s *Serializer) MarshalVariableSizeParts(
	val reflect.Value,
	variableParts [][]byte,
	variableLengths []int,
) ([][]byte, []int, error) {
	serialized, err := s.MarshalSSZ(val.Interface())
	if err != nil {
		return variableParts, variableLengths, err
	}
	variableParts = append(variableParts, serialized)
	variableLengths = append(variableLengths, len(serialized))
	return variableParts, variableLengths, nil
}

func (s *Serializer) MarshalStruct(
	val reflect.Value,
	typ reflect.Type,
	buf *[]byte,
	startOffset uint64,
) (uint64, error) {
	if !IsStruct(typ, val) {
		return 0, errors.New("input is not a struct")
	}

	var fixedParts [][]byte
	var variableParts [][]byte
	var fixedLengths []int
	var variableLengths []int
	var errCheck []error

	var processStructField = func(
		typ reflect.Type,
		val reflect.Value,
		field reflect.StructField,
		err error,
	) {
		if err != nil {
			errCheck = append(errCheck, err)
			return
		}

		var serializationErr error
		// If the field has a ssz-size tag set, we treat it as a fixed size
		// field
		if hasUndefinedSizeTag(field) && isVariableSizeType(typ) {
			variableParts,
				variableLengths,
				serializationErr = s.MarshalVariableSizeParts(
				val,
				variableParts,
				variableLengths,
			)
			// We create holes in fixedParts using nil
			// which is where we slot in offsets in interleaveOffsets.
			fixedParts = append(fixedParts, nil)
			fixedLengths = append(fixedLengths, BytesPerLengthOffset)
		} else {
			fixedParts,
				fixedLengths,
				serializationErr = s.MarshalFixedSizeParts(
				val,
				fixedParts,
				fixedLengths,
			)
			// We populate variable parts with an empty item based on the
			// spec
			variableParts = append(variableParts, make([]byte, 0))
			variableLengths = append(variableLengths, 0)
		}
		if serializationErr != nil {
			errCheck = append(errCheck, serializationErr)
			return
		}
	}

	IterStructFields(
		val,
		processStructField,
	)
	if len(errCheck) > 0 {
		return 0, errCheck[0]
	}

	// Check lengths and
	// Interleave offsets of variable-size parts with fixed-size parts.
	res, err := InterleaveOffsets(
		fixedParts,
		fixedLengths,
		variableParts,
		variableLengths,
	)
	if err != nil {
		return 0, err
	}
	SafeCopyBuffer(res, buf, startOffset)
	return uint64(len(res)), nil
}

func (s *Serializer) MarshalComposite(
	val reflect.Value,
	_ reflect.Type,
	buf *[]byte,
	startOffset uint64,
) (uint64, error) {
	var fixedParts [][]byte
	var variableParts [][]byte
	var fixedLengths []int
	var variableLengths []int
	var errCheck []error

	var processMember = func(
		c interface{},
	) {
		memberTyp := reflect.TypeOf(c)
		memberVal := reflect.ValueOf(c)
		var serializationErr error
		// If the field has a ssz-size tag set, we treat it as a fixed size
		// field
		if isVariableSizeType(memberTyp) {
			variableParts,
				variableLengths,
				serializationErr = s.MarshalVariableSizeParts(
				memberVal,
				variableParts,
				variableLengths,
			)
			// spec-deviation: we differ from the ssz.dev composite spec but
			// align
			// with fastssz output by not writing a nil value into fixedParts.
			fixedLengths = append(fixedLengths, BytesPerLengthOffset)
		} else {
			fixedParts,
				fixedLengths,
				serializationErr = s.MarshalFixedSizeParts(
				memberVal,
				fixedParts,
				fixedLengths,
			)
			// We populate variable parts with an empty item based on the
			// spec
			variableParts = append(variableParts, make([]byte, 0))
			variableLengths = append(variableLengths, 0)
		}
		if serializationErr != nil {
			errCheck = append(errCheck, serializationErr)
			return
		}
	}

	// Start the recursive serialization
	for i := range val.Len() {
		processMember(val.Index(i).Interface())
	}
	if len(errCheck) > 0 {
		return 0, errCheck[0]
	}

	// Check lengths and
	// Interleave offsets of variable-size parts with fixed-size parts.
	res, err := InterleaveOffsets(
		fixedParts,
		fixedLengths,
		variableParts,
		variableLengths,
	)
	if err != nil {
		return 0, err
	}
	SafeCopyBuffer(res, buf, startOffset)
	return uint64(len(res)), nil
}
