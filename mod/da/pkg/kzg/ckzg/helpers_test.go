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

package ckzg_test

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // this is a test.
var globalVerifier *ckzg.Verifier

//nolint:gochecknoglobals // this is a test.
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
