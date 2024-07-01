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
// conditions
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

package types

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// EthSecp256k1CredentialPrefix is the prefix for an Ethereum secp256k1.
const EthSecp256k1CredentialPrefix = byte(iota + 1)

// WithdrawalCredentials is a staking credential that is used to identify a
// validator.
type WithdrawalCredentials common.Bytes32

// NewCredentialsFromExecutionAddress creates a new WithdrawalCredentials from
// an.
func NewCredentialsFromExecutionAddress(
	address gethprimitives.ExecutionAddress,
) WithdrawalCredentials {
	credentials := WithdrawalCredentials{}
	credentials[0] = 0x01
	copy(credentials[12:], address[:])
	return credentials
}

// ToExecutionAddress converts the WithdrawalCredentials to an ExecutionAddress.
func (wc WithdrawalCredentials) ToExecutionAddress() (
	gethprimitives.ExecutionAddress,
	error,
) {
	if wc[0] != EthSecp256k1CredentialPrefix {
		return gethprimitives.ZeroAddress, ErrInvalidWithdrawalCredentials
	}
	return gethprimitives.ExecutionAddress(wc[12:]), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Bytes32.
// TODO: Figure out how to not have to do this.
func (wc *WithdrawalCredentials) UnmarshalJSON(input []byte) error {
	return (*common.Bytes32)(wc).UnmarshalJSON(input)
}

// String returns the hex string representation of Bytes32.
// TODO: Figure out how to not have to do this.
func (wc WithdrawalCredentials) String() string {
	return common.Bytes32(wc).String()
}

// MarshalText implements the encoding.TextMarshaler interface for Bytes32.
// TODO: Figure out how to not have to do this.
func (wc WithdrawalCredentials) MarshalText() ([]byte, error) {
	return common.Bytes32(wc).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Bytes32.
// TODO: Figure out how to not have to do this.
func (wc *WithdrawalCredentials) UnmarshalText(text []byte) error {
	return (*common.Bytes32)(wc).UnmarshalText(text)
}
