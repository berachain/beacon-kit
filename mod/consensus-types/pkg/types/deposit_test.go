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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
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
	buf := make([]byte, deposit.SizeSSZ())
	target, err := deposit.MarshalSSZTo(buf)
	require.NoError(t, err)
	require.NotNil(t, target)
}

func TestDeposit_HashTreeRoot(t *testing.T) {
	deposit := generateValidDeposit()

	_, err := deposit.HashTreeRoot()
	require.NoError(t, err)
}

func TestDeposit_SizeSSZ(t *testing.T) {
	deposit := generateValidDeposit()

	require.Equal(t, 192, deposit.SizeSSZ())
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

	require.ErrorIs(t, err, ssz.ErrSize)
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
	require.Equal(t, deposit.Index, deposit.GetIndex())
}
