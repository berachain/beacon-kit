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

package hex_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero", 0, "0x0"},
		{"MaxByte", 255, "0xff"},
		{"MaxWord", 65535, "0xffff"},
		{"MaxDWord", 4294967295, "0xffffffff"},
		{"MaxQWord", 18446744073709551615, "0xffffffffffffffff"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := hex.MarshalText(test.input)
			require.NoError(t, err, "Test case %s", test.name)
			require.Equal(
				t,
				test.expected,
				string(result),
				"Test case %s",
				test.name,
			)
		})
	}
}

func TestValidateUnmarshalInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected error
	}{
		{"ValidQuotedString", []byte(`"0x0"`), nil},
		{"ValidQuotedStringFF", []byte(`"0xff"`), nil},
		{"NonQuotedString", []byte(`0xff`), hex.ErrNonQuotedString},
		{"InvalidQuotedString", []byte(`"z`), hex.ErrNonQuotedString},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := hex.ValidateUnmarshalInput(test.input)
			if test.expected != nil {
				require.Equal(t, test.expected, err, "Test case %s", test.name)
			} else {
				require.NoError(t, err, "Test case %s", test.name)
			}
		})
	}
}

func TestUnmarshalUint64Text(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected uint64
		err      error
	}{
		{"Zero", []byte("0x0"), 0, nil},
		{"MaxByte", []byte("0xff"), 255, nil},
		{"MaxWord", []byte("0xffff"), 65535, nil},
		{"MaxDWord", []byte("0xffffffff"), 4294967295, nil},
		{"MaxQWord", []byte("0xffffffffffffffff"), 18446744073709551615, nil},
		{"OutOfRange", []byte("0x10000000000000000"), 0, hex.ErrUint64Range},
		{"InvalidString", []byte("0xzz"), 0, hex.ErrInvalidString},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := hex.UnmarshalUint64Text(test.input)
			if test.err != nil {
				require.Equal(t, test.err, err, "Test case %s", test.name)
			} else {
				require.Equal(t, test.expected, result, "Test case %s", test.name)
				require.NoError(t, err)
			}
		})
	}
}
