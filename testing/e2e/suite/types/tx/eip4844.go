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

package tx

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// New4844Tx creates a new 4844 tx.
func New4844Tx(
	nonce uint64, to *common.Address, gasLimit uint64,
	chainID, tip, feeCap, value *big.Int, code []byte,
	blobFeeCap *big.Int, blobData []byte, al types.AccessList,
) *types.Transaction {
	blobs, commits, aggProof, versionedHashes, err := EncodeBlobs(blobData)
	if err != nil {
		panic(err)
	}
	if to == nil {
		to = &common.Address{}
	}
	return types.NewTx(&types.BlobTx{
		ChainID:    uint256.MustFromBig(chainID),
		Nonce:      nonce,
		GasTipCap:  uint256.MustFromBig(tip),
		GasFeeCap:  uint256.MustFromBig(feeCap),
		Gas:        gasLimit,
		To:         *to,
		Value:      uint256.MustFromBig(value),
		Data:       code,
		AccessList: al,
		BlobFeeCap: uint256.MustFromBig(blobFeeCap),
		BlobHashes: versionedHashes,
		Sidecar: &types.BlobTxSidecar{
			Blobs:       blobs,
			Commitments: commits,
			Proofs:      aggProof,
		},
	})
}

// EncodeBlobs encodes blobs.
func EncodeBlobs(
	data []byte,
) (
	[]kzg4844.Blob, []kzg4844.Commitment,
	[]kzg4844.Proof, []common.Hash, error,
) {
	blobs := encodeBlobs(data)
	commits := make([]kzg4844.Commitment, 0, len(blobs))
	proofs := make([]kzg4844.Proof, 0, len(blobs))
	versionedHashes := make([]common.Hash, 0, len(blobs))
	for i := range blobs {
		commit, err := kzg4844.BlobToCommitment(&blobs[i])
		if err != nil {
			return nil, nil, nil, nil, err
		}
		commits = append(commits, commit)

		proof, err := kzg4844.ComputeBlobProof(&blobs[i], commit)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		proofs = append(proofs, proof)

		versionedHashes = append(
			versionedHashes,
			eip4844.KZGCommitment(commit).ToVersionedHash(),
		)
	}
	return blobs, commits, proofs, versionedHashes, nil
}

// encodeBlobs encodes data into blobs.
//
//nolint:mnd // its chill.
func encodeBlobs(data []byte) []kzg4844.Blob {
	blobs := []kzg4844.Blob{{}}
	blobIndex := 0
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			blobs = append(blobs, kzg4844.Blob{})
			blobIndex++
			fieldIndex = 0
		}
		maxSize := i + 31
		if maxSize > len(data) {
			maxSize = len(data)
		}
		copy(blobs[blobIndex][fieldIndex*32+1:], data[i:maxSize])
	}
	return blobs
}
