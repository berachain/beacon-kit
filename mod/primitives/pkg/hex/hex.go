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

package hex

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strconv"
)

// String represents a hex string with 0x prefix.
// Invariants: IsEmpty(s) > 0, has0xPrefix(s) == true.
type String string

// NewString creates a hex string with 0x prefix. It modifies the input to
// ensure that the string invariants are satisfied.
func NewString[T []byte | string](s T) String {
	str := string(s)
	str = ensureStringInvariants(str)
	return String(str)
}

// NewStringStrict creates a hex string with 0x prefix. It errors if any of the
// string invariants are violated.
func NewStringStrict[T []byte | string](s T) (String, error) {
	str := string(s)
	if len(str) == 0 {
		return "", ErrEmptyString
	} else if !has0xPrefix(str) {
		return "", ErrMissingPrefix
	}
	return String(str), nil
}

// FromBytes creates a hex string with 0x prefix.
func FromBytes[B ~[]byte](b B) String {
	enc := make([]byte, len(b)*2+prefixLen)
	copy(enc, prefix)
	hex.Encode(enc[2:], b)
	return NewString(enc)
}

// FromUint64 encodes i as a hex string with 0x prefix.
func FromUint64[U ~uint64](i U) String {
	enc := make([]byte, prefixLen, initialCapacity)
	copy(enc, prefix)
	//#nosec:G701 // i is a uint64, so it can't overflow.
	return String(strconv.AppendUint(enc, uint64(i), hexBase))
}

// FromBigInt encodes bigint as a hex string with 0x prefix.
func FromBigInt(bigint *big.Int) String {
	if sign := bigint.Sign(); sign == 0 {
		return NewString("0x0")
	} else if sign > 0 {
		return NewString("0x" + bigint.Text(hexBase))
	}
	return NewString("-0x" + bigint.Text(hexBase)[1:])
}

func FromJSONString[B ~[]byte](b B) String {
	return NewString(bytes.Trim(b, "\""))
}

// Has0xPrefix returns true if s has a 0x prefix.
func (s String) Has0xPrefix() bool {
	return has0xPrefix[string](string(s))
}

// IsEmpty returns true if s is empty.
func (s String) IsEmpty() bool {
	return len(s) == 0
}

// ToBytes decodes a hex string with 0x prefix.
func (s String) ToBytes() ([]byte, error) {
	return hex.DecodeString(string(s[2:]))
}

// MustToBytes decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToBytes() []byte {
	b, err := s.ToBytes()
	if err != nil {
		panic(err)
	}
	return b
}

// ToUint64 decodes a hex string with 0x prefix.
func (s String) ToUint64() (uint64, error) {
	raw, err := formatAndValidateNumber(s.Unwrap())
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(raw, 16, 64)
}

// MustToUInt64 decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToUInt64() uint64 {
	i, err := s.ToUint64()
	if err != nil {
		panic(err)
	}
	return i
}

// ToBigInt decodes a hex string with 0x prefix.
func (s String) ToBigInt() (*big.Int, error) {
	raw, err := formatAndValidateNumber(s.Unwrap())
	if err != nil {
		return nil, err
	}
	if len(raw) > bytesPer256Bits {
		return nil, ErrBig256Range
	}
	bigWordNibbles, err := getBigWordNibbles()
	if err != nil {
		return nil, err
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return nil, ErrInvalidString
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	dec := new(big.Int).SetBits(words)
	return dec, nil
}

// MustToBigInt decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToBigInt() *big.Int {
	bi, err := s.ToBigInt()
	if err != nil {
		panic(err)
	}
	return bi
}

func (s String) AddQuotes() String {
	return "\"" + s + "\""
}

// Unwrap returns the string value.
func (s String) Unwrap() string {
	return string(s)
}
