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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/stretchr/testify/require"
)

func TestNewCredentialsFromExecutionAddress(t *testing.T) {
	address := primitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	expectedCredentials := types.DepositCredentials{}
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
		primitives.ExecutionAddress(expectedCredentials[12:]),
		"Expected address to be set correctly",
	)
	credentials := types.NewCredentialsFromExecutionAddress(address)
	require.Equal(
		t,
		expectedCredentials,
		credentials,
		"Generated credentials do not match expected",
	)
}

func TestToExecutionAddress(t *testing.T) {
	expectedAddress := primitives.ExecutionAddress{0xde, 0xad, 0xbe, 0xef}
	credentials := types.DepositCredentials{}
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
	require.NoError(t, err, "Conversion to execution address should not error")
	require.Equal(
		t,
		expectedAddress,
		address,
		"Converted address does not match expected",
	)
}

func TestToExecutionAddress_InvalidPrefix(t *testing.T) {
	credentials := types.DepositCredentials{}
	for i := range credentials {
		credentials[i] = 0x00 // Invalid prefix
	}

	_, err := credentials.ToExecutionAddress()

	require.Error(t, err, "Expected an error due to invalid prefix")
}
