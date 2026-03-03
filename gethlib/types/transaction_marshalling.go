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
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
)

var (
	errInvalidYParity   = errors.New("'yParity' field must be 0 or 1")
	errVYParityMismatch = errors.New("'v' and 'yParity' fields do not match")
	errVYParityMissing  = errors.New("missing 'yParity' or 'v' field in transaction")
)

const (
	recoveryIDByteLen      = 8
	replayProtectionBitLen = 64
	legacyVValue27         = 27
	legacyVValue28         = 28
	replayProtectionBase   = 35
	chainIDDivisor         = 2
)

// txJSON is the JSON representation of transactions.
type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID              *hexutil.Big                     `json:"chainId,omitempty"`
	From                 *common.Address                  `json:"from,omitempty"`
	Nonce                *hexutil.Uint64                  `json:"nonce"`
	To                   *common.Address                  `json:"to"`
	Gas                  *hexutil.Uint64                  `json:"gas"`
	GasPrice             *hexutil.Big                     `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big                     `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big                     `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *hexutil.Big                     `json:"maxFeePerBlobGas,omitempty"`
	Value                *hexutil.Big                     `json:"value"`
	Input                *hexutil.Bytes                   `json:"input"`
	AccessList           *coretypes.AccessList            `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash                    `json:"blobVersionedHashes,omitempty"`
	AuthorizationList    []coretypes.SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                    *hexutil.Big                     `json:"v"`
	R                    *hexutil.Big                     `json:"r"`
	S                    *hexutil.Big                     `json:"s"`
	YParity              *hexutil.Uint64                  `json:"yParity,omitempty"`

	// Blob transaction sidecar encoding:
	Blobs       []kzg4844.Blob       `json:"blobs,omitempty"`
	Commitments []kzg4844.Commitment `json:"commitments,omitempty"`
	Proofs      []kzg4844.Proof      `json:"proofs,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

// yParityValue returns the YParity value from JSON. For backwards-compatibility reasons,
// this can be given in the 'v' field or the 'yParity' field. If both exist, they must match.
func (tx *txJSON) yParityValue() (*big.Int, error) {
	if tx.YParity != nil {
		val := uint64(*tx.YParity)
		if val != 0 && val != 1 {
			return nil, errInvalidYParity
		}
		bigval := new(big.Int).SetUint64(val)
		if tx.V != nil && tx.V.ToInt().Cmp(bigval) != 0 {
			return nil, errVYParityMismatch
		}
		return bigval, nil
	}
	if tx.V != nil {
		return tx.V.ToInt(), nil
	}
	return nil, errVYParityMissing
}

// MarshalJSON marshals as JSON with a hash.
//
//nolint:funlen // Mirrors geth transaction JSON shape for wire compatibility.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	var enc txJSON
	hash := tx.Hash()

	// These are set for all tx types.
	enc.Hash = hash
	enc.Type = hexutil.Uint64(tx.Type())

	// Other fields are set conditionally depending on tx type.
	switch itx := tx.inner.(type) {
	case *LegacyTx:
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = itx.To
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		if itx.V != nil && isProtectedV(itx.V) {
			enc.ChainID = (*hexutil.Big)(deriveChainID(itx.V))
		}

	case *AccessListTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = itx.To
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		if itx.V != nil {
			yparity := itx.V.Uint64()
			enc.YParity = (*hexutil.Uint64)(&yparity)
		}

	case *DynamicFeeTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = itx.To
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		if itx.V != nil {
			yparity := itx.V.Uint64()
			enc.YParity = (*hexutil.Uint64)(&yparity)
		}

	case *BlobTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID.ToBig())
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap.ToBig())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap.ToBig())
		enc.MaxFeePerBlobGas = (*hexutil.Big)(itx.BlobFeeCap.ToBig())
		enc.Value = (*hexutil.Big)(itx.Value.ToBig())
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.BlobVersionedHashes = itx.BlobHashes
		enc.To = &itx.To
		enc.V = (*hexutil.Big)(itx.V.ToBig())
		enc.R = (*hexutil.Big)(itx.R.ToBig())
		enc.S = (*hexutil.Big)(itx.S.ToBig())
		yparity := itx.V.Uint64()
		enc.YParity = (*hexutil.Uint64)(&yparity)
		if sidecar := itx.Sidecar; sidecar != nil {
			enc.Blobs = itx.Sidecar.Blobs
			enc.Commitments = itx.Sidecar.Commitments
			enc.Proofs = itx.Sidecar.Proofs
		}

	case *SetCodeTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID.ToBig())
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = &itx.To
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap.ToBig())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap.ToBig())
		enc.Value = (*hexutil.Big)(itx.Value.ToBig())
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.AuthorizationList = itx.AuthList
		enc.V = (*hexutil.Big)(itx.V.ToBig())
		enc.R = (*hexutil.Big)(itx.R.ToBig())
		enc.S = (*hexutil.Big)(itx.S.ToBig())
		yparity := itx.V.Uint64()
		enc.YParity = (*hexutil.Uint64)(&yparity)

	case *PoLTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.From = &itx.From
		enc.To = &itx.To
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		gas := hexutil.Uint64(itx.GasLimit)
		enc.Gas = &gas
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
	}
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
//
//nolint:gocognit,funlen,gocyclo,cyclop,maintidx // Mirrors geth transaction JSON shape for wire compatibility.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Decode / verify fields according to transaction type.
	var inner TxData
	txType := uint64(dec.Type)
	switch txType {
	case uint64(coretypes.LegacyTxType):
		var itx LegacyTx
		inner = &itx
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input

		// signature R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// signature S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// signature V
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			sigErr := sanityCheckSignature(itx.V, itx.R, itx.S, true)
			if sigErr != nil {
				return sigErr
			}
		}

	case uint64(coretypes.AccessListTxType):
		var itx AccessListTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}

		// signature R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// signature S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// signature V
		vParity, parityErr := dec.yParityValue()
		if parityErr != nil {
			return parityErr
		}
		itx.V = vParity
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			sigErr := sanityCheckSignature(itx.V, itx.R, itx.S, false)
			if sigErr != nil {
				return sigErr
			}
		}

	case uint64(coretypes.DynamicFeeTxType):
		var itx DynamicFeeTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}

		// signature R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// signature S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// signature V
		vParity, parityErr := dec.yParityValue()
		if parityErr != nil {
			return parityErr
		}
		itx.V = vParity
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			sigErr := sanityCheckSignature(itx.V, itx.R, itx.S, false)
			if sigErr != nil {
				return sigErr
			}
		}

	case uint64(coretypes.BlobTxType):
		var itx BlobTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		var overflow bool
		itx.ChainID, overflow = uint256.FromBig(dec.ChainID.ToInt())
		if overflow {
			return errors.New("'chainId' value overflows uint256")
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To == nil {
			return errors.New("missing required field 'to' in transaction")
		}
		itx.To = *dec.To
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = uint256.MustFromBig((*big.Int)(dec.MaxPriorityFeePerGas))
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerGas))
		if dec.MaxFeePerBlobGas == nil {
			return errors.New("missing required field 'maxFeePerBlobGas' for txdata")
		}
		itx.BlobFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerBlobGas))
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = uint256.MustFromBig((*big.Int)(dec.Value))
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		if dec.BlobVersionedHashes == nil {
			return errors.New("missing required field 'blobVersionedHashes' in transaction")
		}
		itx.BlobHashes = dec.BlobVersionedHashes

		// signature R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R, overflow = uint256.FromBig((*big.Int)(dec.R))
		if overflow {
			return errors.New("'r' value overflows uint256")
		}
		// signature S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S, overflow = uint256.FromBig((*big.Int)(dec.S))
		if overflow {
			return errors.New("'s' value overflows uint256")
		}
		// signature V
		vParity, parityErr := dec.yParityValue()
		if parityErr != nil {
			return parityErr
		}
		itx.V, overflow = uint256.FromBig(vParity)
		if overflow {
			return errors.New("'v' value overflows uint256")
		}
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			sigErr := sanityCheckSignature(vParity, itx.R.ToBig(), itx.S.ToBig(), false)
			if sigErr != nil {
				return sigErr
			}
		}

	case uint64(coretypes.SetCodeTxType):
		var itx SetCodeTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		var overflow bool
		itx.ChainID, overflow = uint256.FromBig(dec.ChainID.ToInt())
		if overflow {
			return errors.New("'chainId' value overflows uint256")
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To == nil {
			return errors.New("missing required field 'to' in transaction")
		}
		itx.To = *dec.To
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = uint256.MustFromBig((*big.Int)(dec.MaxPriorityFeePerGas))
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerGas))
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = uint256.MustFromBig((*big.Int)(dec.Value))
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		if dec.AuthorizationList == nil {
			return errors.New("missing required field 'authorizationList' in transaction")
		}
		itx.AuthList = dec.AuthorizationList

		// signature R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R, overflow = uint256.FromBig((*big.Int)(dec.R))
		if overflow {
			return errors.New("'r' value overflows uint256")
		}
		// signature S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S, overflow = uint256.FromBig((*big.Int)(dec.S))
		if overflow {
			return errors.New("'s' value overflows uint256")
		}
		// signature V
		vParity, parityErr := dec.yParityValue()
		if parityErr != nil {
			return parityErr
		}
		itx.V, overflow = uint256.FromBig(vParity)
		if overflow {
			return errors.New("'v' value overflows uint256")
		}
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			sigErr := sanityCheckSignature(vParity, itx.R.ToBig(), itx.S.ToBig(), false)
			if sigErr != nil {
				return sigErr
			}
		}

	case uint64(PoLTxType):
		var itx PoLTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		itx.From = *dec.From
		if dec.To == nil {
			return errors.New("missing required field 'to' in transaction")
		}
		itx.To = *dec.To
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.GasLimit = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input

	default:
		return coretypes.ErrTxTypeNotSupported
	}

	// Set the inner transaction.
	tx.setDecoded(inner)
	return nil
}

func sanityCheckSignature(v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return coretypes.ErrUnexpectedProtection
	}

	var plainV uint64
	switch {
	case isProtectedV(v):
		chainID := deriveChainID(v).Uint64()
		plainV = v.Uint64() - replayProtectionBase - chainIDDivisor*chainID
	case maybeProtected:
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = v.Uint64() - legacyVValue27
	default:
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = v.Uint64()
	}

	var recoveryID byte
	switch plainV {
	case 0:
		recoveryID = 0
	case 1:
		recoveryID = 1
	default:
		return coretypes.ErrInvalidSig
	}
	if !crypto.ValidateSignatureValues(recoveryID, r, s, false) {
		return coretypes.ErrInvalidSig
	}

	return nil
}

func isProtectedV(v *big.Int) bool {
	if v.BitLen() <= recoveryIDByteLen {
		val := v.Uint64()
		return val != legacyVValue27 && val != legacyVValue28 && val != 1 && val != 0
	}
	// Anything not 27 or 28 is considered protected.
	return true
}

// deriveChainID derives the chain ID from a signature v value.
func deriveChainID(v *big.Int) *big.Int {
	if v.BitLen() <= replayProtectionBitLen {
		val := v.Uint64()
		if val == legacyVValue27 || val == legacyVValue28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((val - replayProtectionBase) / chainIDDivisor)
	}
	vCopy := new(big.Int).Sub(v, big.NewInt(replayProtectionBase))
	return vCopy.Rsh(vCopy, 1)
}
