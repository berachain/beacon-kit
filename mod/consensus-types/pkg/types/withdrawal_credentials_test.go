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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/stretchr/testify/require"
)

func TestNewCredentialsFromExecutionAddress(t *testing.T) {
	address := gethprimitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	expectedCredentials := types.WithdrawalCredentials{}
	expectedCredentials[0] = 0x01 // EthSecp256k1CredentialPrefix
	copy(expectedCredentials[12:], address[:])
	for i := 1; i < 12; i++ {
		expectedCredentials[i] = 0x00
	}
	require.Len(
		t,
		expectedCredentials,
		32,
		"Expected credentials to be 32 bytes long",
	)
	require.Equal(
		t,
		byte(0x01),
		expectedCredentials[0],
		"Expected prefix to be 0x01",
	)
	require.Equal(
		t,
		address,
		gethprimitives.ExecutionAddress(expectedCredentials[12:]),
		"Expected address to be set correctly",
	)
	credentials := types.
		NewCredentialsFromExecutionAddress(address)
	require.Equal(
		t,
		expectedCredentials,
		credentials,
		"Generated credentials do not match expected",
	)
}

func TestToExecutionAddress(t *testing.T) {
	expectedAddress := gethprimitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	credentials := types.WithdrawalCredentials{}
	for i := range credentials {
		// First byte should be 0x01
		switch {
		case i == 0:
			credentials[i] = 0x01 // EthSecp256k1CredentialPrefix
		case i > 0 && i < 12:
			credentials[i] = 0x00 // then we have 11 bytes of padding
		default:
			credentials[i] = expectedAddress[i-12] // then the address
		}
	}

	address, err := credentials.ToExecutionAddress()
	require.NoError(t, err,
		"Conversion to execution address should not error")
	require.Equal(
		t,
		expectedAddress,
		address,
		"Converted address does not match expected",
	)
}

func TestToExecutionAddress_InvalidPrefix(t *testing.T) {
	credentials := types.WithdrawalCredentials{}
	for i := range credentials {
		credentials[i] = 0x00 // Invalid prefix
	}

	_, err := credentials.ToExecutionAddress()

	require.Error(t, err, "Expected an error due to invalid prefix")
}

func TestWithdrawalCredentials_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.WithdrawalCredentials
		wantErr  bool
	}{
		{
			name: "valid JSON",
			input: `"0x0100000000000000000000000000000` +
				`000000000000000000000000000000000"`,
			expected: types.WithdrawalCredentials{0x01},
			wantErr:  false,
		},
		{
			name:    "invalid JSON",
			input:   `"invalid"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wc types.WithdrawalCredentials

			err := wc.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err, "Test case %s", tt.name)
				require.Equal(t, tt.expected, wc, "Test case %s", tt.name)
			}
		})
	}
}

func TestWithdrawalCredentials_String(t *testing.T) {
	tests := []struct {
		name     string
		wc       types.WithdrawalCredentials
		expected string
	}{
		{
			name: "valid string",
			wc:   types.WithdrawalCredentials{0x01},
			expected: "0x010000000000000000000000000000" +
				"0000000000000000000000000000000000",
		},
		{
			name: "valid string with full address",
			wc: types.WithdrawalCredentials{0x01, 0xde, 0xad, 0xbe, 0xef, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00},
			expected: "0x01deadbeef00000000000000000000000" +
				"0000000000000000000000000000000",
		},
		{
			name: "empty credentials",
			wc:   types.WithdrawalCredentials{},
			expected: "0x0000000000000000000000000000000000" +
				"000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.wc.String(),
				"Test case %s", tt.name)
		})
	}
}

func TestWithdrawalCredentials_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		wc       types.WithdrawalCredentials
		expected []byte
		wantErr  bool
	}{
		{
			name: "valid marshal",
			wc:   types.WithdrawalCredentials{0x01},
			expected: []byte("0x010000000000000000000000000000000000" +
				"0000000000000000000000000000"),
			wantErr: false,
		},
		{
			name: "valid marshal with different prefix",
			wc:   types.WithdrawalCredentials{0x02},
			expected: []byte("0x020000000000000000000000000000000000" +
				"0000000000000000000000000000"),
			wantErr: false,
		},
		{
			name: "valid marshal with full address",
			wc: types.WithdrawalCredentials{0x01, 0xde, 0xad, 0xbe, 0xef, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00},
			expected: []byte("0x01deadbeef00000000000000000000000000000" +
				"0000000000000000000000000"),
			wantErr: false,
		},
		{
			name: "empty credentials",
			wc:   types.WithdrawalCredentials{},
			expected: []byte("0x0000000000000000000000000000000000000000" +
				"000000000000000000000000"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.wc.MarshalText()
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err, "Test case %s", tt.name)
				require.Equal(t, tt.expected, result,
					"Test case %s", tt.name)
			}
		})
	}
}

func TestWithdrawalCredentials_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected types.WithdrawalCredentials
		wantErr  bool
	}{
		{
			name: "valid unmarshal",
			input: []byte("0x010000000000000000000000000000000" +
				"0000000000000000000000000000000"),
			expected: types.WithdrawalCredentials{0x01},
			wantErr:  false,
		},
		{
			name:    "invalid unmarshal",
			input:   []byte("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wc types.WithdrawalCredentials
			err := wc.UnmarshalText(tt.input)
			if tt.wantErr {
				require.Error(t, err, "Test case %s", tt.name)
			} else {
				require.NoError(t, err, "Test case %s", tt.name)
				require.Equal(t, tt.expected, wc,
					"Test case %s", tt.name)
			}
		})
	}
}
