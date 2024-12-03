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

package ckzg_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/da/kzg/ckzg"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/json"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var globalVerifier *ckzg.Verifier

var baseDir = "../../../../../testing/files/"

func TestMain(m *testing.M) {
	var err error
	globalVerifier, err = setupVerifier()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func setupVerifier() (*ckzg.Verifier, error) {
	fs := afero.NewOsFs()
	fileName := "kzg-trusted-setup.json"
	fullPath := filepath.Join(baseDir, fileName)
	file, err := afero.ReadFile(fs, fullPath)

	if err != nil {
		return nil, err
	}

	var ts gokzg4844.JSONTrustedSetup
	if errUnmarshal := json.Unmarshal(file, &ts); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	verifier, errVerifier := ckzg.NewVerifier(&ts)
	if errVerifier != nil {
		return nil, errVerifier
	}
	return verifier, nil
}

func setupTestData(t *testing.T, fileName string) (
	*eip4844.Blob, eip4844.KZGProof, eip4844.KZGCommitment) {
	t.Helper()

	filePath := filepath.Join(baseDir, fileName)
	data, err := afero.ReadFile(afero.NewOsFs(), filePath)
	require.NoError(t, err)

	type Test struct {
		Input struct {
			Blob       string `json:"blob"`
			Commitment string `json:"commitment"`
			Proof      string `json:"proof"`
		} `json:"input"`
	}
	var test Test

	err = json.Unmarshal(data, &test)
	require.NoError(t, err)

	var blob eip4844.Blob
	errBlob := blob.UnmarshalJSON([]byte(`"` + test.Input.Blob + `"`))
	require.NoError(t, errBlob)

	var commitment eip4844.KZGCommitment
	errCommitment := commitment.UnmarshalJSON(
		[]byte(`"` + test.Input.Commitment + `"`),
	)
	require.NoError(t, errCommitment)

	var proof eip4844.KZGProof
	errProof := proof.UnmarshalJSON([]byte(`"` + test.Input.Proof + `"`))
	require.NoError(t, errProof)

	return &blob, proof, commitment
}
