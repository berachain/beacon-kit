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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/version"
	cmtcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/privval"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func generateFakeSignedBeaconBlock() *types.SignedBeaconBlock {
	blk := generateValidBeaconBlock()
	signature := crypto.BLSSignature{}
	return &types.SignedBeaconBlock{
		Message:   blk,
		Signature: signature,
	}
}

func generatePrivKey() (cmtcrypto.PrivKey, error) {
	privKey, err := bls12381.GenPrivKey()
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func generateSigningRoot(blk *types.BeaconBlock) (common.Root, error) {
	cs, err := chain.NewSpec(spec.BaseSpec())
	if err != nil {
		return common.Root{}, err
	}
	forkData := types.ForkData{}
	domain := forkData.ComputeDomain(cs.DomainTypeProposer())
	signingRoot := types.ComputeSigningRoot(blk, domain)
	return signingRoot, nil
}

func generateRealSignedBeaconBlock(blsSigner crypto.BLSSigner) (*types.SignedBeaconBlock, error) {
	blk := generateValidBeaconBlock()

	signingRoot, err := generateSigningRoot(blk)
	if err != nil {
		return nil, err
	}
	signature, err := blsSigner.Sign(signingRoot[:])
	if err != nil {
		return nil, err
	}
	return &types.SignedBeaconBlock{
		Message:   blk,
		Signature: signature,
	}, nil
}

// TestNewSignedBeaconBlockFromSSZ tests the roundtrip SSZ encoding for Deneb.
func TestNewSignedBeaconBlockFromSSZ(t *testing.T) {
	originalBlock := generateFakeSignedBeaconBlock()
	blockBytes, err := originalBlock.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, blockBytes)

	newBlock, err := types.NewSignedBeaconBlockFromSSZ(
		blockBytes, version.Deneb,
	)
	require.NoError(t, err)
	require.NotNil(t, newBlock)
	require.Equal(t, originalBlock, newBlock)
}

func TestNewSignedBeaconBlockFromSSZForkVersionNotSupported(t *testing.T) {
	_, err := types.NewSignedBeaconBlockFromSSZ([]byte{}, 1)
	require.ErrorIs(t, err, types.ErrForkVersionNotSupported)
}

func TestSignedBeaconBlock_HashTreeRoot(_ *testing.T) {
	sBlk := generateFakeSignedBeaconBlock()
	sBlk.HashTreeRoot()
}

// TestSignedBeaconBlock_SignBeaconBlock ensures the validity of the block
// signatures.
func TestSignedBeaconBlock_SignBeaconBlock(t *testing.T) {
	// Generate a new bls key signer
	filePV, err := privval.GenFilePV(
		"signed_beacon_block_test_filepv_key",
		"signed_beacon_block_test_filepv_state",
		generatePrivKey,
	)
	require.NoError(t, err)
	blsSigner := signer.BLSSigner{PrivValidator: filePV}

	// Generate real signed beacon block
	signedBlk, err := generateRealSignedBeaconBlock(blsSigner)
	require.NoError(t, err)
	require.NotNil(t, signedBlk)

	// Use SignBeaconBlock to sign the same BeaconBlock
	cs, err := chain.NewSpec(spec.BaseSpec())
	require.NoError(t, err)
	newSignedBlk, err := types.NewSignedBeaconBlock(
		signedBlk.GetMessage(),
		&types.ForkData{},
		cs,
		blsSigner,
	)
	require.NoError(t, err)

	// Check that the signature from SignBeaconBlock matches
	sig1 := signedBlk.GetSignature()
	sig2 := newSignedBlk.GetSignature()
	require.Equal(t, sig1, sig2)

	// Verify the signature is good
	signingRoot, err := generateSigningRoot(newSignedBlk.GetMessage())
	require.NoError(t, err)
	err = blsSigner.VerifySignature(blsSigner.PublicKey(), signingRoot[:], newSignedBlk.GetSignature())
	require.NoError(t, err)
}

func TestSignedBeaconBlock_SizeSSZ(t *testing.T) {
	sBlk := generateFakeSignedBeaconBlock()
	size := ssz.Size(sBlk)
	require.Positive(t, size)
}

func TestSignedBeaconBlock_EmptySerialization(t *testing.T) {
	orig := &types.SignedBeaconBlock{}
	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.SignedBeaconBlock
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.NotNil(t, unmarshalled.GetMessage())
	require.NotNil(t, unmarshalled.GetSignature())

	buf := make([]byte, ssz.Size(orig))
	err = ssz.EncodeToBytes(buf, orig)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}
