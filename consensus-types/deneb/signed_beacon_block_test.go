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

package deneb_test

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/consensus-types/deneb"
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

func generateFakeSignedBeaconBlock(t *testing.T) *deneb.SignedBeaconBlock {
	t.Helper()

	blk := generateValidBeaconBlock(t)
	signature := crypto.BLSSignature{}
	return &deneb.SignedBeaconBlock{
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

func generateSigningRoot(blk *deneb.BeaconBlock) (common.Root, error) {
	cs, err := spec.DevnetChainSpec()
	if err != nil {
		return common.Root{}, err
	}
	domain := (&deneb.ForkData{}).ComputeDomain(cs.DomainTypeProposer())
	signingRoot := deneb.ComputeSigningRoot(blk, domain)
	return signingRoot, nil
}

func generateRealSignedBeaconBlock(t *testing.T, blsSigner crypto.BLSSigner) (*deneb.SignedBeaconBlock, error) {
	t.Helper()

	blk := generateValidBeaconBlock(t)

	signingRoot, err := generateSigningRoot(blk)
	if err != nil {
		return nil, err
	}
	signature, err := blsSigner.Sign(signingRoot[:])
	if err != nil {
		return nil, err
	}
	return &deneb.SignedBeaconBlock{
		Message:   blk,
		Signature: signature,
	}, nil
}

// TestNewSignedBeaconBlockFromSSZ tests the roundtrip SSZ encoding for Deneb.
func TestNewSignedBeaconBlockFromSSZ(t *testing.T) {
	t.Parallel()
	originalBlock := generateFakeSignedBeaconBlock(t)
	blockBytes, err := originalBlock.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, blockBytes)

	newBlock, err := deneb.NewSignedBeaconBlockFromSSZ(
		blockBytes, originalBlock.Message.Version(),
	)
	require.NoError(t, err)
	require.NotNil(t, newBlock)
	require.Equal(t, originalBlock, newBlock)
}

func TestNewSignedBeaconBlockFromSSZForkVersionNotSupported(t *testing.T) {
	t.Parallel()
	_, err := deneb.NewSignedBeaconBlockFromSSZ([]byte{}, version.Altair())
	require.ErrorIs(t, err, deneb.ErrForkVersionNotSupported)
}

func TestSignedBeaconBlock_HashTreeRoot(t *testing.T) {
	t.Parallel()
	sBlk := generateFakeSignedBeaconBlock(t)
	sBlk.HashTreeRoot()
}

// TestSignedBeaconBlock_SignBeaconBlock ensures the validity of the block
// signatures.
func TestSignedBeaconBlock_SignBeaconBlock(t *testing.T) {
	t.Parallel()
	// Generate a new bls key signer
	filePV, err := privval.GenFilePV(
		"signed_beacon_block_test_filepv_key",
		"signed_beacon_block_test_filepv_state",
		generatePrivKey,
	)
	require.NoError(t, err)
	blsSigner := signer.BLSSigner{PrivValidator: filePV}

	// Generate real signed beacon block
	signedBlk, err := generateRealSignedBeaconBlock(t, blsSigner)
	require.NoError(t, err)
	require.NotNil(t, signedBlk)

	// Use SignBeaconBlock to sign the same BeaconBlock
	cs, err := spec.DevnetChainSpec()
	require.NoError(t, err)
	newSignedBlk, err := deneb.NewSignedBeaconBlock(
		signedBlk.GetMessage(),
		&deneb.ForkData{},
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
	t.Parallel()
	sBlk := generateFakeSignedBeaconBlock(t)
	size := ssz.Size(sBlk)
	require.Positive(t, size)
}

func TestSignedBeaconBlock_EmptySerialization(t *testing.T) {
	t.Parallel()
	orig := &deneb.SignedBeaconBlock{}
	data, err := orig.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled deneb.SignedBeaconBlock
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
