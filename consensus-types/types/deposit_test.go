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

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	ssz "github.com/ferranbt/fastssz"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// generateValidDeposit generates a valid deposit for testing purposes.
func generateValidDeposit() *types.Deposit {
	var pubKey crypto.BLSPubkey
	var signature crypto.BLSSignature
	var credentials types.WithdrawalCredentials
	amount := math.Gwei(32)
	index := uint64(1)

	return &types.Deposit{
		Pubkey:      pubKey,
		Credentials: credentials,
		Amount:      amount,
		Signature:   signature,
		Index:       index,
	}
}

func TestDeposit_New(t *testing.T) {
	deposit := generateValidDeposit()
	newDeposit := deposit.New(
		deposit.Pubkey,
		deposit.Credentials,
		deposit.Amount,
		deposit.Signature,
		deposit.Index,
	)
	require.Equal(t, deposit, newDeposit)
}

func TestDeposit_MarshalUnmarshalSSZ(t *testing.T) {
	originalDeposit := generateValidDeposit()

	// Marshal the original deposit to SSZ
	sszDeposit, err := originalDeposit.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, sszDeposit)

	var unmarshalledDeposit types.Deposit
	err = unmarshalledDeposit.UnmarshalSSZ(sszDeposit)
	require.NoError(t, err)

	require.Equal(t, originalDeposit, &unmarshalledDeposit)
}

func TestDeposit_MarshalSSZTo(t *testing.T) {
	deposit := generateValidDeposit()
	buf := make([]byte, karalabessz.Size(deposit))
	target, err := deposit.MarshalSSZTo(buf)
	require.NoError(t, err)
	require.NotNil(t, target)
}

func TestDeposit_HashTreeRoot(t *testing.T) {
	deposit := generateValidDeposit()
	require.NotPanics(t, func() {
		_ = deposit.HashTreeRoot()
	})
}

func TestDeposit_SizeSSZ(t *testing.T) {
	deposit := generateValidDeposit()

	require.Equal(t, uint32(192), karalabessz.Size(deposit))
}

func TestDeposit_HashTreeRootWith(t *testing.T) {
	deposit := generateValidDeposit()
	require.NotNil(t, deposit)
	hasher := ssz.NewHasher()
	require.NotNil(t, hasher)
	err := deposit.HashTreeRootWith(hasher)
	require.NoError(t, err)
}

func TestDeposit_GetTree(t *testing.T) {
	deposit := generateValidDeposit()
	_, err := deposit.GetTree()
	require.NoError(t, err)
}

func TestDeposit_UnmarshalSSZ_ErrSize(t *testing.T) {
	// Create a byte slice of incorrect size
	buf := make([]byte, 10) // size less than 192

	var unmarshalledDeposit types.Deposit
	err := unmarshalledDeposit.UnmarshalSSZ(buf)

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestDeposit_VerifySignature(t *testing.T) {
	deposit := generateValidDeposit()

	forkData := &types.ForkData{
		CurrentVersion:        common.Version{0x00, 0x00, 0x00, 0x04},
		GenesisValidatorsRoot: common.Root{0x00, 0x00, 0x00, 0x00},
	}

	signatureVerificationFn := func(
		_ crypto.BLSPubkey, _ []byte, _ crypto.BLSSignature,
	) error {
		return nil
	}

	errVerify := deposit.VerifySignature(forkData, common.DomainType{
		0x01, 0x00, 0x00, 0x00,
	}, signatureVerificationFn)
	require.NoError(t, errVerify)
}

func TestDeposit_Getters(t *testing.T) {
	deposit := generateValidDeposit()

	require.Equal(t, deposit.Pubkey, deposit.GetPubkey())
	require.Equal(t, deposit.Credentials, deposit.GetWithdrawalCredentials())
	require.Equal(t, deposit.Amount, deposit.GetAmount())
	require.Equal(t, deposit.Signature, deposit.GetSignature())
	require.Equal(t, math.U64(deposit.Index), deposit.GetIndex())
}
