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
	"io"
	"testing"

	types "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/mocks"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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

	mocksSigner := &mocks.Blssigner{}
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

	buf := make([]byte, original.SizeSSZ())
	data, err := original.MarshalSSZTo(buf)
	require.NoError(t, err)

	var unmarshalled types.DepositMessage
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestDepositMessage_UnmarshalSSZ_ErrSize(t *testing.T) {
	buf := make([]byte, 10) // size less than 88

	var unmarshalledDepositMessage types.DepositMessage
	err := unmarshalledDepositMessage.UnmarshalSSZ(buf)

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
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
