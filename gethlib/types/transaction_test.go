// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package types

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestTransactionAccessors(t *testing.T) {
	t.Parallel()

	to := common.HexToAddress("0x1000000000000000000000000000000000000001")
	poLTo := common.HexToAddress("0x2000000000000000000000000000000000000002")
	accessList := AccessList{{
		Address:     common.HexToAddress("0x3000000000000000000000000000000000000003"),
		StorageKeys: []common.Hash{common.HexToHash("0x01")},
	}}
	blobHashes := []common.Hash{
		common.HexToHash("0x10"),
		common.HexToHash("0x20"),
	}
	sidecar := &BlobTxSidecar{Version: coretypes.BlobSidecarVersion1}

	tests := []struct {
		name              string
		tx                *Transaction
		wantAccessList    AccessList
		wantBlobGas       uint64
		wantBlobGasFeeCap *big.Int
		wantBlobHashes    []common.Hash
		wantBlobTxSidecar *BlobTxSidecar
		wantData          []byte
		wantGas           uint64
		wantGasFeeCap     *big.Int
		wantGasPrice      *big.Int
		wantGasTipCap     *big.Int
		wantNonce         uint64
		wantTo            *common.Address
		wantValue         *big.Int
	}{
		{
			name: "legacy",
			tx: &Transaction{inner: &LegacyTx{
				Nonce:    1,
				GasPrice: big.NewInt(2),
				Gas:      3,
				To:       &to,
				Value:    big.NewInt(4),
				Data:     []byte{0x05},
			}},
			wantData:      []byte{0x05},
			wantGas:       3,
			wantGasFeeCap: big.NewInt(2),
			wantGasPrice:  big.NewInt(2),
			wantGasTipCap: big.NewInt(2),
			wantNonce:     1,
			wantTo:        &to,
			wantValue:     big.NewInt(4),
		},
		{
			name: "access list",
			tx: &Transaction{inner: &AccessListTx{
				Nonce:      2,
				GasPrice:   big.NewInt(3),
				Gas:        4,
				To:         &to,
				Value:      big.NewInt(5),
				Data:       []byte{0x06},
				AccessList: accessList,
			}},
			wantAccessList: accessList,
			wantData:       []byte{0x06},
			wantGas:        4,
			wantGasFeeCap:  big.NewInt(3),
			wantGasPrice:   big.NewInt(3),
			wantGasTipCap:  big.NewInt(3),
			wantNonce:      2,
			wantTo:         &to,
			wantValue:      big.NewInt(5),
		},
		{
			name: "dynamic fee",
			tx: &Transaction{inner: &DynamicFeeTx{
				Nonce:      3,
				GasTipCap:  big.NewInt(4),
				GasFeeCap:  big.NewInt(5),
				Gas:        6,
				To:         &to,
				Value:      big.NewInt(7),
				Data:       []byte{0x08},
				AccessList: accessList,
			}},
			wantAccessList: accessList,
			wantData:       []byte{0x08},
			wantGas:        6,
			wantGasFeeCap:  big.NewInt(5),
			wantGasPrice:   big.NewInt(5),
			wantGasTipCap:  big.NewInt(4),
			wantNonce:      3,
			wantTo:         &to,
			wantValue:      big.NewInt(7),
		},
		{
			name: "blob",
			tx: &Transaction{inner: &BlobTx{
				Nonce:      4,
				GasTipCap:  uint256.NewInt(5),
				GasFeeCap:  uint256.NewInt(6),
				Gas:        7,
				To:         to,
				Value:      uint256.NewInt(8),
				Data:       []byte{0x09},
				AccessList: accessList,
				BlobFeeCap: uint256.NewInt(10),
				BlobHashes: blobHashes,
				Sidecar:    sidecar,
			}},
			wantAccessList:    accessList,
			wantBlobGas:       params.BlobTxBlobGasPerBlob * uint64(len(blobHashes)),
			wantBlobGasFeeCap: big.NewInt(10),
			wantBlobHashes:    blobHashes,
			wantBlobTxSidecar: sidecar,
			wantData:          []byte{0x09},
			wantGas:           7,
			wantGasFeeCap:     big.NewInt(6),
			wantGasPrice:      big.NewInt(6),
			wantGasTipCap:     big.NewInt(5),
			wantNonce:         4,
			wantTo:            &to,
			wantValue:         big.NewInt(8),
		},
		{
			name: "set code",
			tx: &Transaction{inner: &SetCodeTx{
				Nonce:      5,
				GasTipCap:  uint256.NewInt(6),
				GasFeeCap:  uint256.NewInt(7),
				Gas:        8,
				To:         to,
				Value:      uint256.NewInt(9),
				Data:       []byte{0x0a},
				AccessList: accessList,
			}},
			wantAccessList: accessList,
			wantData:       []byte{0x0a},
			wantGas:        8,
			wantGasFeeCap:  big.NewInt(7),
			wantGasPrice:   big.NewInt(7),
			wantGasTipCap:  big.NewInt(6),
			wantNonce:      5,
			wantTo:         &to,
			wantValue:      big.NewInt(9),
		},
		{
			name: "pol",
			tx: &Transaction{inner: &PoLTx{
				Nonce:    6,
				GasLimit: 9,
				GasPrice: big.NewInt(10),
				To:       poLTo,
				Data:     []byte{0x0b},
			}},
			wantData:      []byte{0x0b},
			wantGas:       9,
			wantGasFeeCap: big.NewInt(10),
			wantGasPrice:  big.NewInt(10),
			wantGasTipCap: common.Big0,
			wantNonce:     6,
			wantTo:        &poLTo,
			wantValue:     common.Big0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.wantAccessList, tt.tx.AccessList())
			require.Equal(t, tt.wantBlobGas, tt.tx.BlobGas())
			requireBigIntEqual(t, tt.wantBlobGasFeeCap, tt.tx.BlobGasFeeCap())
			require.Equal(t, tt.wantBlobHashes, tt.tx.BlobHashes())
			require.Same(t, tt.wantBlobTxSidecar, tt.tx.BlobTxSidecar())
			require.Equal(t, tt.wantData, tt.tx.Data())
			require.Equal(t, tt.wantGas, tt.tx.Gas())
			requireBigIntEqual(t, tt.wantGasFeeCap, tt.tx.GasFeeCap())
			requireBigIntEqual(t, tt.wantGasPrice, tt.tx.GasPrice())
			requireBigIntEqual(t, tt.wantGasTipCap, tt.tx.GasTipCap())
			require.Equal(t, tt.wantNonce, tt.tx.Nonce())
			requireAddressPtrEqual(t, tt.wantTo, tt.tx.To())
			require.True(t, tt.tx.Time().IsZero())
			requireBigIntEqual(t, tt.wantValue, tt.tx.Value())
		})
	}
}

func TestTransactionAccessorsReturnCopies(t *testing.T) {
	t.Parallel()

	to := common.HexToAddress("0x1000000000000000000000000000000000000001")
	tx := &Transaction{inner: &DynamicFeeTx{
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(2),
		To:        &to,
		Value:     big.NewInt(3),
	}}

	tx.GasFeeCap().SetUint64(100)
	tx.GasPrice().SetUint64(101)
	tx.GasTipCap().SetUint64(102)
	tx.Value().SetUint64(103)
	*tx.To() = common.HexToAddress("0x2000000000000000000000000000000000000002")

	requireBigIntEqual(t, big.NewInt(2), tx.GasFeeCap())
	requireBigIntEqual(t, big.NewInt(2), tx.GasPrice())
	requireBigIntEqual(t, big.NewInt(1), tx.GasTipCap())
	requireBigIntEqual(t, big.NewInt(3), tx.Value())
	require.Equal(t, to, *tx.To())
}

func TestTransactionTimeSetOnDecode(t *testing.T) {
	t.Parallel()

	tx := new(Transaction)
	before := time.Now()
	tx.setDecoded(&LegacyTx{GasPrice: big.NewInt(1), Value: big.NewInt(2)})
	after := time.Now()

	require.False(t, tx.Time().Before(before))
	require.False(t, tx.Time().After(after))
}

func TestTransactionToReturnsNilForContractCreation(t *testing.T) {
	t.Parallel()

	tx := &Transaction{inner: &LegacyTx{GasPrice: big.NewInt(1), Value: big.NewInt(2)}}
	require.Nil(t, tx.To())
}

func requireAddressPtrEqual(t *testing.T, want *common.Address, got *common.Address) {
	t.Helper()
	if want == nil {
		require.Nil(t, got)
		return
	}
	require.NotNil(t, got)
	require.Equal(t, *want, *got)
}

func requireBigIntEqual(t *testing.T, want *big.Int, got *big.Int) {
	t.Helper()
	if want == nil {
		require.Nil(t, got)
		return
	}
	require.NotNil(t, got)
	require.Zero(t, want.Cmp(got))
}
