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

package deposit_test // Change to your package name

import (
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/mod/node-builder/commands/deposit"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// errToString converts error to string, treating nil as empty.
func errToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// TestConvertFunctions tests various string conversion functions.
func TestConvertFunctions(t *testing.T) {
	var tests = []struct {
		name           string
		convertFunc    func(string) (interface{}, error)
		input          string
		expectedOutput interface{}
		expectedErr    string
		description    string // Describes the test purpose
	}{
		{
			name: "Valid Pubkey",
			convertFunc: func(s string) (interface{}, error) {
				return deposit.ConvertPubkey(s)
			},
			input: "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccdd",
			expectedOutput: primitives.BLSPubkey{
				0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11,
				0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
				0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11,
				0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
			},
			expectedErr: "",
			description: "Converts a valid hex string to a BLSPubkey.",
		},
		// {
		// 	name: "Invalid Pubkey Length",
		// 	convertFunc: func(s string) (interface{}, error) {
		// 		return deposit.ConvertPubkey(s)
		// 	},
		// 	input:          "aabbccdd",
		// 	expectedOutput: primitives.BLSPubkey{},
		// 	expectedErr:    deposit.ErrInvalidPubKeyLength.Error(),
		// 	description:    "Handles pubkey strings that are too short.",
		// },
		// {
		// 	name: "Valid WithdrawalCredentials",
		// 	convertFunc: func(s string) (interface{}, error) {
		// 		return deposit.ConvertWithdrawalCredentials(s)
		// 	},
		// 	input:
		// "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
		// 	expectedOutput: types.WithdrawalCredentials{},
		// 	expectedErr:    "",
		// 	description:    "Correctly converts a valid hex string to
		// WithdrawalCredentials.",
		// },
		// {
		// 	name: "Invalid WithdrawalCredentials Length",
		// 	convertFunc: func(s string) (interface{}, error) {
		// 		return deposit.ConvertWithdrawalCredentials(s)
		// 	},
		// 	input:          "aabbccdd",
		// 	expectedOutput: types.WithdrawalCredentials{},
		// 	expectedErr:    "ErrInvalidWithdrawalCredentialsLength",
		// 	description:    "Checks for WithdrawalCredentials string length
		// errors.",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.convertFunc(tt.input)
			if errToString(err) != tt.expectedErr {
				t.Errorf("%s: %s - expected error %s, got %s",
					tt.name, tt.description, tt.expectedErr, errToString(err))
			}
			if !reflect.DeepEqual(output, tt.expectedOutput) &&
				tt.expectedErr == "" {
				t.Errorf("%s: %s - expected output %v, got %v",
					tt.name, tt.description, tt.expectedOutput, output)
			}
		})
	}
}
