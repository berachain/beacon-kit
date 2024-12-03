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

package kzg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/da/kzg"
	"github.com/berachain/beacon-kit/da/kzg/ckzg"
	"github.com/berachain/beacon-kit/da/kzg/gokzg"
	"github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var baseDir = "../../../../testing/files/"

func TestNewBlobProofVerifier_KzgImpl(t *testing.T) {
	ts, err := loadTrustedSetupFromFile()
	require.NoError(t, err)

	verifier, err := kzg.NewBlobProofVerifier(gokzg.Implementation, ts)
	require.NoError(t, err)
	require.NotNil(t, verifier)
	require.Equal(t, gokzg.Implementation, verifier.GetImplementation())
}

func TestNewBlobProofVerifier_CkzgImpl(t *testing.T) {
	ts, err := loadTrustedSetupFromFile()
	require.NoError(t, err)

	verifier, err := kzg.NewBlobProofVerifier(ckzg.Implementation, ts)
	require.NoError(t, err)
	require.NotNil(t, verifier)
	require.Equal(t, ckzg.Implementation, verifier.GetImplementation())
}

func TestNewBlobProofVerifier_InvalidImpl(t *testing.T) {
	ts, err := loadTrustedSetupFromFile()
	require.NoError(t, err)

	invalidImpl := "invalid-implementation"
	_, err = kzg.NewBlobProofVerifier(invalidImpl, ts)
	require.ErrorIs(t, err, kzg.ErrUnsupportedKzgImplementation)
}

// loadTrustedSetupFromFile is the helper function.
func loadTrustedSetupFromFile() (*gokzg4844.JSONTrustedSetup, error) {
	fileName := "kzg-trusted-setup.json"
	fullPath := filepath.Join(baseDir, fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var ts gokzg4844.JSONTrustedSetup
	err = json.Unmarshal(data, &ts)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

func TestArgsFromSidecars(t *testing.T) {
	fs := afero.NewOsFs()
	fullPath := filepath.Join(baseDir, "test_data.json")
	file, err := afero.ReadFile(fs, fullPath)
	require.NoError(t, err)

	type Data struct {
		Input struct {
			Blob       string `json:"blob"`
			Commitment string `json:"commitment"`
			Proof      string `json:"proof"`
		} `json:"input"`
	}
	var data Data

	err = json.Unmarshal(file, &data)
	require.NoError(t, err)

	scs := &types.BlobSidecars{
		Sidecars: []*types.BlobSidecar{
			{
				Blob:          eip4844.Blob{data.Input.Blob[0]},
				KzgProof:      eip4844.KZGProof{data.Input.Proof[0]},
				KzgCommitment: eip4844.KZGCommitment{data.Input.Commitment[0]},
			},
		},
	}

	args := kzg.ArgsFromSidecars[
		*types.BlobSidecar,
		*types.BlobSidecars,
	](scs)

	require.Len(t, args.Blobs, 1)
	require.Len(t, args.Proofs, 1)
	require.Len(t, args.Commitments, 1)
}
