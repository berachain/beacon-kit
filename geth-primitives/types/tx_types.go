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
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// BRIP-0004 PoL transaction type.
const PoLTxType = 0x7E

// AccessListTx is the data of EIP-2930 access list transactions.
type AccessListTx struct {
	ChainID    *big.Int             // destination chain ID
	Nonce      uint64               // nonce of sender account
	GasPrice   *big.Int             // wei per gas
	Gas        uint64               // gas limit
	To         *common.Address      `rlp:"nil"` // nil means contract creation
	Value      *big.Int             // wei amount
	Data       []byte               // contract invocation input data
	AccessList coretypes.AccessList // EIP-2930 access list
	V, R, S    *big.Int             // signature values
}

func (tx *AccessListTx) txType() byte { return coretypes.AccessListTxType }

func (tx *AccessListTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *AccessListTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

// BlobTx represents an EIP-4844 transaction.
type BlobTx struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList coretypes.AccessList
	BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes []common.Hash

	// A blob transaction can optionally contain blobs. This field must be set when BlobTx
	// is used to create a transaction for signing.
	Sidecar *coretypes.BlobTxSidecar `rlp:"-"`

	// Signature values
	V *uint256.Int
	R *uint256.Int
	S *uint256.Int
}

func (tx *BlobTx) txType() byte { return coretypes.BlobTxType }

func (tx *BlobTx) encode(b *bytes.Buffer) error {
	switch {
	case tx.Sidecar == nil:
		return rlp.Encode(b, tx)

	case tx.Sidecar.Version == coretypes.BlobSidecarVersion0:
		return rlp.Encode(b, &blobTxWithBlobsV0{
			BlobTx:      tx,
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})

	case tx.Sidecar.Version == coretypes.BlobSidecarVersion1:
		return rlp.Encode(b, &blobTxWithBlobsV1{
			BlobTx:      tx,
			Version:     tx.Sidecar.Version,
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})

	default:
		return errors.New("unsupported sidecar version")
	}
}

func (tx *BlobTx) decode(input []byte) error {
	// Here we need to support two outer formats: the network protocol encoding of the tx
	// (with blobs) or the canonical encoding without blobs.
	//
	// The canonical encoding is just a list of fields:
	//
	//     [chainID, nonce, ...]
	//
	// The network encoding is a list where the first element is the tx in the canonical encoding,
	// and the remaining elements are the 'sidecar':
	//
	//     [[chainID, nonce, ...], ...]
	//
	// The two outer encodings can be distinguished by checking whether the first element
	// of the input list is itself a list. If it's the canonical encoding, the first
	// element is the chainID, which is a number.

	firstElem, _, err := rlp.SplitList(input)
	if err != nil {
		return err
	}
	firstElemKind, _, secondElem, err := rlp.Split(firstElem)
	if err != nil {
		return err
	}
	if firstElemKind != rlp.List {
		// Blob tx without blobs.
		return rlp.DecodeBytes(input, tx)
	}

	// Now we know it's the network encoding with the blob sidecar. Here we again need to
	// support multiple encodings: legacy sidecars (v0) with a blob proof, and versioned
	// sidecars.
	//
	// The legacy encoding is:
	//
	//     [tx, blobs, commitments, proofs]
	//
	// The versioned encoding is:
	//
	//     [tx, version, blobs, ...]
	//
	// We can tell the two apart by checking whether the second element is the version byte.
	// For legacy sidecar the second element is a list of blobs.

	secondElemKind, _, _, err := rlp.Split(secondElem)
	if err != nil {
		return err
	}
	var payload blobTxWithBlobs
	if secondElemKind == rlp.List {
		// No version byte: blob sidecar v0.
		payload = new(blobTxWithBlobsV0)
	} else {
		// It has a version byte. Decode as v1, version is checked by assign()
		payload = new(blobTxWithBlobsV1)
	}
	if err := rlp.DecodeBytes(input, payload); err != nil {
		return err
	}
	sc := new(coretypes.BlobTxSidecar)
	if err := payload.assign(sc); err != nil {
		return err
	}
	*tx = *payload.tx()
	tx.Sidecar = sc
	return nil
}

// blobTxWithBlobs represents blob tx with its corresponding sidecar.
// This is an interface because sidecars are versioned.
type blobTxWithBlobs interface {
	tx() *BlobTx
	assign(*coretypes.BlobTxSidecar) error
}

type blobTxWithBlobsV0 struct {
	BlobTx      *BlobTx
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

type blobTxWithBlobsV1 struct {
	BlobTx      *BlobTx
	Version     byte
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func (btx *blobTxWithBlobsV0) tx() *BlobTx {
	return btx.BlobTx
}

func (btx *blobTxWithBlobsV0) assign(sc *coretypes.BlobTxSidecar) error {
	sc.Version = coretypes.BlobSidecarVersion0
	sc.Blobs = btx.Blobs
	sc.Commitments = btx.Commitments
	sc.Proofs = btx.Proofs
	return nil
}

func (btx *blobTxWithBlobsV1) tx() *BlobTx {
	return btx.BlobTx
}

func (btx *blobTxWithBlobsV1) assign(sc *coretypes.BlobTxSidecar) error {
	if btx.Version != coretypes.BlobSidecarVersion1 {
		return fmt.Errorf("unsupported blob tx version %d", btx.Version)
	}
	sc.Version = coretypes.BlobSidecarVersion1
	sc.Blobs = btx.Blobs
	sc.Commitments = btx.Commitments
	sc.Proofs = btx.Proofs
	return nil
}

// DynamicFeeTx represents an EIP-1559 transaction.
type DynamicFeeTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList coretypes.AccessList

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int
}

func (tx *DynamicFeeTx) txType() byte { return coretypes.DynamicFeeTxType }

func (tx *DynamicFeeTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *DynamicFeeTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

// LegacyTx is the transaction data of the original Ethereum transactions.
type LegacyTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

func (tx *LegacyTx) txType() byte { return coretypes.LegacyTxType }

func (tx *LegacyTx) encode(*bytes.Buffer) error {
	panic("encode called on LegacyTx")
}

func (tx *LegacyTx) decode([]byte) error {
	panic("decode called on LegacyTx)")
}

// PoLTx represents an BRIP-0004 transaction. No gas is consumed for execution.
type PoLTx struct {
	ChainID  *big.Int
	From     common.Address // system address
	To       common.Address // address of the PoL Distributor contract
	Nonce    uint64         // block number distributing for
	GasLimit uint64         // artificial gas limit for the PoL tx, not consumed against the block gas limit
	GasPrice *big.Int       // gas price is set to the baseFee to make the tx valid for EIP-1559 rules
	Data     []byte         // encodes the pubkey distributing for
}

func (*PoLTx) txType() byte { return PoLTxType }

func (tx *PoLTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *PoLTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

// SetCodeTx implements the EIP-7702 transaction type which temporarily installs
// the code at the signer's address.
type SetCodeTx struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList coretypes.AccessList
	AuthList   []coretypes.SetCodeAuthorization

	// Signature values
	V *uint256.Int
	R *uint256.Int
	S *uint256.Int
}

func (tx *SetCodeTx) txType() byte { return coretypes.SetCodeTxType }

func (tx *SetCodeTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *SetCodeTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
