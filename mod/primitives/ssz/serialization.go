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
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

const (
	BytesPerChunk        int = 32 // Number of bytes per chunk.
	BytesPerLengthOffset int = 4  // Number of bytes per serialized length offset.
	BitsPerByte          int = 8  // Number of bits per byte.
)

type Serializable interface {
	MarshalSSZ() ([]byte, error)
	IsVariableSize() bool
	SizeSSZ() int
}

func SerializeBasic[
	U64T U64[U64T],
	U256L U256LT,
	B Basic[RootT, SpecT],
	SpecT any,
	RootT ~[32]byte,
](value B) ([]byte, error) {
	switch el := reflect.ValueOf(value).Interface().(type) {
	case bool:
		var buffer [1]byte
		if el {
			buffer[0] = 1
		}
		return buffer[:], nil
	case uint8:
		var buffer [1]byte
		buffer[0] = el
		return buffer[:], nil
	case uint16:
		var buffer [2]byte
		binary.LittleEndian.PutUint16(buffer[:], el)
		return buffer[:], nil
	case uint32:
		var buffer [4]byte
		binary.LittleEndian.PutUint32(buffer[:], el)
		return buffer[:], nil
	case U64T:
		var buffer [8]byte
		//#nosec:G701 // This is a safe operation.
		binary.LittleEndian.PutUint64(buffer[:], uint64(el))
		return buffer[:], nil
	case U256L:
		return el[:], nil
	}
	return nil, errors.New("unsupported type")
}

func SerializeContainer[
	U64T U64[U64T], RootT ~[32]byte,
	SpecT any, C Container[RootT, SpecT],
](value C) ([]byte, error) {
	rValue := reflect.ValueOf(value)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	numFields := rValue.NumField()
	buf := make([]byte, value.SizeSSZ())

	// initial offset = Sum of all basic fields + 8 * number of variable fields.
	offset := 0 // todo magioc number
	cursor := 0
	variableFields := []int{}

	for i := range numFields {
		field, ok := rValue.Field(i).Interface().(Serializable)
		if !ok {
			return nil, errors.New("field is not hashable")
		}
		fmt.Println(field.SizeSSZ(), field.IsVariableSize())
		fieldSize := field.SizeSSZ()
		if field.IsVariableSize() {
			// Copy offset to the buffer
			// Move cursor forward by 4 bytes.
			// Move offset forward by fieldSize bytes.
			copy(buf[cursor:cursor+4], serializeUint32(uint32(offset)))
			cursor += 4
			offset += fieldSize
			variableFields = append(variableFields, i)
		} else {
			// Copy bytes from the field to the buffer
			// Move cursor forward by size bytes.
			bz, err := field.MarshalSSZ()
			if err != nil {
				return nil, err
			}
			fmt.Println(bz, cursor, fieldSize)
			copy(buf[cursor:cursor+fieldSize], bz)
			cursor += fieldSize
		}
	}

	if cursor != value.SizeSSZ()-offset {
		return nil, errors.New("invalid size")
	}

	// Serialize variable fields
	for x := range variableFields {
		field, ok := rValue.Field(x).Interface().(Serializable)
		if !ok {
			return nil, errors.New("field is not hashable")
		}
		bz, err := field.MarshalSSZ()
		if err != nil {
			return nil, err
		}

		// Copy bytes from the field to the buffer
		// Move cursor forward by size bytes.
		copy(bz, buf[cursor:cursor+len(bz)])
		cursor += len(bz)
	}

	if cursor != value.SizeSSZ() {
		return nil, errors.New("invalid size")
	}
	return buf, nil
}

func serializeUint32(x uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, x)
	return buf
}

func sum(slice []int) int {
	total := 0
	for _, v := range slice {
		total += v
	}
	return total
}
