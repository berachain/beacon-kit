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
	"bytes"
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxData is the underlying data of a transaction.
type TxData interface {
	txType() byte // returns the type ID

	encode(*bytes.Buffer) error
	decode([]byte) error
}

// Transaction is similar to coretypes.Transaction, with added support
// for the BRIP-0004 PoL transaction type.
type Transaction struct {
	inner TxData // Consensus contents of a transaction

	// cache
	hash atomic.Pointer[common.Hash]
}

// BlobHashes returns the hashes of the blob commitments for blob transactions, nil otherwise.
func (tx *Transaction) BlobHashes() []common.Hash {
	if blobtx, ok := tx.inner.(*BlobTx); ok {
		return blobtx.BlobHashes
	}
	return nil
}

// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return *hash
	}

	var h common.Hash
	if tx.Type() == coretypes.LegacyTxType {
		h = rlpHash(tx.inner)
	} else {
		h = prefixedRlpHash(tx.Type(), tx.inner)
	}
	tx.hash.Store(&h)
	return h
}

// Type returns the transaction type.
func (tx *Transaction) Type() uint8 {
	return tx.inner.txType()
}

// RawSignatureValues returns the transaction signature values.
// PoL transactions intentionally return nil values because they are unsigned.
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	switch itx := tx.inner.(type) {
	case *LegacyTx:
		return itx.V, itx.R, itx.S
	case *AccessListTx:
		return itx.V, itx.R, itx.S
	case *DynamicFeeTx:
		return itx.V, itx.R, itx.S
	case *BlobTx:
		if itx.V == nil || itx.R == nil || itx.S == nil {
			return nil, nil, nil
		}
		return itx.V.ToBig(), itx.R.ToBig(), itx.S.ToBig()
	case *SetCodeTx:
		if itx.V == nil || itx.R == nil || itx.S == nil {
			return nil, nil, nil
		}
		return itx.V.ToBig(), itx.R.ToBig(), itx.S.ToBig()
	default:
		return nil, nil, nil
	}
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	if tx.Type() == coretypes.LegacyTxType {
		return rlp.EncodeToBytes(tx.inner)
	}
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP-2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return err
		}
		tx.setDecoded(&data)
		return nil
	}
	// It's an EIP-2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner)
	return nil
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return tx.inner.encode(w)
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction) decodeTyped(b []byte) (TxData, error) {
	if len(b) <= 1 {
		return nil, errors.New("typed transaction too short")
	}
	var inner TxData
	switch b[0] {
	case coretypes.AccessListTxType:
		inner = new(AccessListTx)
	case coretypes.DynamicFeeTxType:
		inner = new(DynamicFeeTx)
	case coretypes.BlobTxType:
		inner = new(BlobTx)
	case coretypes.SetCodeTxType:
		inner = new(SetCodeTx)
	case PoLTxType:
		inner = new(PoLTx)
	default:
		return nil, coretypes.ErrTxTypeNotSupported
	}
	err := inner.decode(b[1:])
	return inner, err
}

// setDecoded sets the inner transaction and clears hash cache.
func (tx *Transaction) setDecoded(inner TxData) {
	tx.inner = inner
	tx.hash.Store(nil)
}

// Transactions implements DerivableList for transactions.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txs that were either
// constructed by decoding or via public API in this package.
func (s Transactions) EncodeIndex(i int, w *bytes.Buffer) {
	tx := s[i]
	if tx.Type() == coretypes.LegacyTxType {
		rlp.Encode(w, tx.inner) //#nosec:G104 copied from go-ethereum
	} else {
		tx.encodeTyped(w) //#nosec:G104 copied from go-ethereum
	}
}
