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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAndSignDepositMessage(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{0x00, 0x00, 0x00, 0x04},
		GenesisValidatorsRoot: common.Root{0x00, 0x00, 0x00, 0x00},
	}

	domainType := common.DomainType{
		0x01, 0x00, 0x00, 0x00,
	}

	mocksSigner := &mocks.BLSSigner{}
	mocksSigner.On("PublicKey").Return(crypto.BLSPubkey{})
	mocksSigner.On("Sign", mock.Anything).Return(crypto.BLSSignature{}, nil)

	credentials := types.WithdrawalCredentials{}
	amount := math.Gwei(32)

	depositMessage, signature, err := types.CreateAndSignDepositMessage(
		forkData, domainType, mocksSigner, credentials, amount,
	)

	require.NoError(t, err)
	require.NotNil(t, depositMessage)
	require.NotNil(t, signature)
}

func TestNewDepositMessage(t *testing.T) {
	pubKey := crypto.BLSPubkey{}
	credentials := types.WithdrawalCredentials{}
	amount := math.Gwei(32)
	depositMessage := types.DepositMessage{}
	newDepositMessage := depositMessage.New(pubKey, credentials, amount)

	require.NotNil(t, newDepositMessage)
}

func TestDepositMessage_MarshalUnmarshalSSZ(t *testing.T) {
	original := &types.DepositMessage{
		Pubkey:      crypto.BLSPubkey{},
		Credentials: types.WithdrawalCredentials{},
		Amount:      math.Gwei(1000),
	}

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)
	var unmarshalled types.DepositMessage
	err = unmarshalled.UnmarshalSSZ(data)

	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestDepositMessage_MarshalSSZTo(t *testing.T) {
	original := &types.DepositMessage{
		Pubkey:      crypto.BLSPubkey{},
		Credentials: types.WithdrawalCredentials{},
		Amount:      math.Gwei(1000),
	}

	buf := make([]byte, 0, original.SizeSSZ())
	data, err := original.MarshalSSZTo(buf)
	require.NoError(t, err)

	var unmarshalled types.DepositMessage
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestDepositMessage_GetTree(t *testing.T) {
	original := &types.DepositMessage{
		Pubkey:      crypto.BLSPubkey{},
		Credentials: types.WithdrawalCredentials{},
		Amount:      math.Gwei(1000),
	}

	tree, err := original.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}
func TestDepositMessage_UnmarshalSSZ_ErrSize(t *testing.T) {
	buf := make([]byte, 10) // size less than 88

	var unmarshalledDepositMessage types.DepositMessage
	err := unmarshalledDepositMessage.UnmarshalSSZ(buf)

	require.ErrorIs(t, err, ssz.ErrSize)
}

func TestDepositMessage_VerifyCreateValidator_Error(t *testing.T) {
	original := &types.DepositMessage{
		Pubkey:      crypto.BLSPubkey{},
		Credentials: types.WithdrawalCredentials{},
		Amount:      math.Gwei(1000),
	}

	forkData := &types.ForkData{
		CurrentVersion:        common.Version{0, 0, 0, 0},
		GenesisValidatorsRoot: common.Root{},
	}

	signature := crypto.BLSSignature{}

	// Define a signature verification function that always returns an error
	signatureVerificationFn := func(
		_ crypto.BLSPubkey, _ []byte, _ crypto.BLSSignature,
	) error {
		return errors.New("signature verification failed")
	}
	domainType := common.DomainType{
		0x01, 0x00, 0x00, 0x00,
	}

	err := original.VerifyCreateValidator(
		forkData, signature, domainType, signatureVerificationFn,
	)

	require.ErrorIs(t, err, types.ErrDepositMessage)
}
