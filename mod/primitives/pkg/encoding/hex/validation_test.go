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

func TestValidateBasicHex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "Valid hex string",
			input:       "0x48656c6c6f",
			expectedErr: nil,
		},
		{
			name:        "Empty string",
			input:       "",
			expectedErr: hex.ErrEmptyString,
		},
		{
			name:        "No 0x prefix",
			input:       "48656c6c6f",
			expectedErr: hex.ErrMissingPrefix,
		},
		{
			name:        "Valid single hex character",
			input:       "0x0",
			expectedErr: nil,
		},
		{
			name:        "Empty hex string",
			input:       "0x",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hex.ValidateBasicHex(tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
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
			err := hex.ValidateQuotedString(test.input)
			if test.expected != nil {
				require.Equal(t, test.expected, err, "Test case %s", test.name)
			} else {
				require.NoError(t, err, "Test case %s", test.name)
			}
		})
	}
}
