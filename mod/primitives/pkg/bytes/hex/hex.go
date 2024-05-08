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
	"encoding/hex"
	"math/big"
	"strconv"
)

// String represents a hex string with 0x prefix.
type String string

// FromBytes creates a hex string with 0x prefix.
func FromBytes(b []byte) String {
	enc := make([]byte, len(b)*2+prefixLen)
	copy(enc, prefix)
	hex.Encode(enc[2:], b)
	return String(enc)
}

// FromUint64 encodes i as a hex string with 0x prefix.
func FromUint64(i uint64) String {
	enc := make([]byte, prefixLen, initialCapacity)
	copy(enc, prefix)
	return String(strconv.AppendUint(enc, i, hexBase))
}

// FromBig encodes bigint as a hex string with 0x prefix.
func FromBig(bigint *big.Int) String {
	if sign := bigint.Sign(); sign == 0 {
		return String("0x0")
	} else if sign > 0 {
		return String("0x" + bigint.Text(hexBase))
	}
	return String("-0x" + bigint.Text(hexBase)[1:])
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
	if s.IsEmpty() {
		return nil, ErrEmptyString
	} else if s.Has0xPrefix() {
		return nil, ErrMissingPrefix
	}
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
	raw := string(s)
	err := validateNumber(raw)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(raw, 16, 64)
}

// MustToUint64 decodes a hex string with 0x prefix.
// It panics for invalid input.
func (s String) MustToUint64() uint64 {
	i, err := s.ToUint64()
	if err != nil {
		panic(err)
	}
	return i
}

// ToBigInt decodes a hex string with 0x prefix.
func (s String) ToBigInt() (*big.Int, error) {
	raw := string(s)
	err := validateNumber(raw)
	if err != nil {
		return nil, err
	}
	if len(raw) > bytesIn256Bits {
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
				return nil, ErrSyntax
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

// Unwrap returns the string value.
func (s String) Unwrap() string {
	return string(s)
}
