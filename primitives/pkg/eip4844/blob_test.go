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

package eip4844_test

import (
	"encoding/hex"
	"testing"

	"github.com/berachain/beacon-kit/primitives/pkg/eip4844"
	"github.com/stretchr/testify/require"
)

func TestBlob_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name: "valid hex input",
			input: []byte(
				`"0x` + hex.EncodeToString(make([]byte, 131072)) + `"`,
			),
			expected: []byte{
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
				0x0,
			},
			wantErr: false,
		},
		{
			name:    "invalid hex input",
			input:   []byte(`"invalidhex"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b eip4844.Blob
			err := b.UnmarshalJSON(tt.input)
			if tt.wantErr {
				require.Error(t, err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, tt.expected, b[:len(tt.expected)],
					"Test case: %s", tt.name)
			}
		})
	}
}

func TestBlob_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    eip4844.Blob
		expected string
	}{
		{
			name: "valid blob",
			input: func() eip4844.Blob {
				var b eip4844.Blob
				copy(
					b[:],
					[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
				)
				return b
			}(),
			expected: func() string {
				var b eip4844.Blob
				copy(
					b[:],
					[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
				)
				return "0x" + hex.EncodeToString(b[:])
			}(),
		},
		{
			name: "all zero bytes",
			input: func() eip4844.Blob {
				var b eip4844.Blob
				copy(b[:], make([]byte, len(b)))
				return b
			}(),
			expected: func() string {
				var b eip4844.Blob
				copy(b[:], make([]byte, len(b)))
				return "0x" + hex.EncodeToString(b[:])
			}(),
		},
		{
			name: "all max bytes",
			input: func() eip4844.Blob {
				var b eip4844.Blob
				for i := range b {
					b[i] = 0xFF
				}
				return b
			}(),
			expected: func() string {
				var b eip4844.Blob
				for i := range b {
					b[i] = 0xFF
				}
				return "0x" + hex.EncodeToString(b[:])
			}(),
		},
		{
			name: "mixed values",
			input: func() eip4844.Blob {
				var b eip4844.Blob
				copy(
					b[:],
					[]byte{0x00, 0xFF, 0xAA, 0x55, 0x11, 0x22, 0x33, 0x44,
						0x88, 0x99, 0x77, 0x66, 0xEE, 0xDD, 0xCC, 0xBB},
				)
				return b
			}(),
			expected: func() string {
				var b eip4844.Blob
				copy(
					b[:],
					[]byte{0x00, 0xFF, 0xAA, 0x55, 0x11, 0x22, 0x33, 0x44,
						0x88, 0x99, 0x77, 0x66, 0xEE, 0xDD, 0xCC, 0xBB},
				)
				return "0x" + hex.EncodeToString(b[:])
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.input.MarshalText()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(
				t,
				tt.expected,
				string(output),
				"Test case: %s",
				tt.name,
			)
		})
	}
}
