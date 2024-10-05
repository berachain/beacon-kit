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

package hex

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"github.com/berachain/beacon-kit/mod/errors"
)

// String represents a hex string with 0x prefix.
// Invariants: IsEmpty(s) > 0, has0xPrefix(s) == true.
type String string

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It validates the input text as a hex string and
// assigns it to the String type.
// Returns an error if the input is not a valid hex string.
func (s *String) UnmarshalText(text []byte) error {
	if _, err := IsValidHex(text); err != nil {
		return errors.Wrapf(ErrInvalidString, "%s", text)
	}
	*s = String(text)
	return nil
}

// IsValidHex performs basic validations that every hex string
// must pass (there may be extra ones depending on the type encoded)
// It returns the suffix (dropping 0x prefix) in the hope to appease nilaway.
func IsValidHex[T ~[]byte | ~string](s T) (T, error) {
	if len(s) == 0 {
		return *new(T), ErrEmptyString
	}
	if len(s) < prefixLen {
		return *new(T), ErrMissingPrefix
	}
	if strings.ToLower(string(s[:prefixLen])) != prefix {
		return *new(T), ErrMissingPrefix
	}
	return s[prefixLen:], nil
}

// NewStringStrict creates a hex string with 0x prefix. It errors if any of the
// string invariants are violated.
func NewStringStrict[T []byte | string](s T) (String, error) {
	if _, err := IsValidHex(s); err != nil {
		return "", err
	}
	return String(s), nil
}

// NewString creates a hex string with 0x prefix. It modifies the input to
// ensure that the string invariants are satisfied.
func NewString[T []byte | string](s T) String {
	str := string(s)
	switch _, err := IsValidHex(str); {
	case err == nil:
		break // already well formatted
	case errors.Is(err, ErrEmptyString):
		str = prefix + "0"
	default:
		str = prefix + str
	}
	return String(str)
}

// FromBytes creates a hex string with 0x prefix.
func FromBytes[B ~[]byte](input B) String {
	b := EncodeBytes(input)
	return NewString(b)
}

// FromUint64 encodes i as a hex string with 0x prefix.
func FromUint64[U ~uint64](i U) String {
	enc := make([]byte, prefixLen, initialCapacity)
	copy(enc, prefix)
	//#nosec:G701 // i is a uint64, so it can't overflow.
	return String(strconv.AppendUint(enc, uint64(i), hexBase))
}

// FromBigInt encodes bigint as a hex string with 0x prefix.
// Precondition: bigint is non-negative.
func FromBigInt(bigint *big.Int) String {
	if sign := bigint.Sign(); sign == 0 {
		return NewString(prefix + "0")
	} else if sign > 0 {
		return NewString(prefix + bigint.Text(hexBase))
	}
	// this return should never reach if precondition is met
	return NewString(prefix + bigint.Text(hexBase)[1:])
}

func FromJSONString[B ~[]byte](b B) String {
	return NewString(bytes.Trim(b, "\""))
}

// ToBytes decodes a hex string with 0x prefix.
func (s String) ToBytes() ([]byte, error) {
	return hex.DecodeString(string(s[prefixLen:]))
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
	if len(raw) > nibblesPer256Bits {
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
