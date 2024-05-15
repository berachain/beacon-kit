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

//go:build !ckzg

package ckzg_test

import (
	"encoding/json"
	"os"
	"testing"

	ckzg "github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var verifier *ckzg.Verifier

var (
	validBlob       = &eip4844.Blob{}
	validProof      = eip4844.KZGProof{}
	validCommitment = eip4844.KZGCommitment{}
)

func TestMain(m *testing.M) {
	// Load the trusted setup before any tests are run
	fs := afero.NewOsFs()
	file, err := afero.ReadFile(fs, "./files/kzg-trusted-setup.json")
	if err != nil {
		panic(err)
	}

	var ts gokzg4844.JSONTrustedSetup
	err = json.Unmarshal(file, &ts)
	if err != nil {
		panic(err)
	}

	verifier, err = ckzg.NewVerifier(&ts)
	if err != nil {
		panic(err)
	}

	// Run the tests
	os.Exit(m.Run())
}

func setup(t *testing.T, filePath string) {
	data, err := os.ReadFile(filePath)
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

	errBlob := validBlob.UnmarshalJSON([]byte(`"` + test.Input.Blob + `"`))
	require.NoError(t, errBlob)

	errCommitment := validCommitment.UnmarshalJSON([]byte(
		`"` + test.Input.Commitment + `"`))
	require.NoError(t, errCommitment)

	errProof := validProof.UnmarshalJSON([]byte(`"` + test.Input.Proof + `"`))
	require.NoError(t, errProof)
}

func TestVerifyBlobKZGProof(t *testing.T) {
	setup(t, "./files/test_data.json")
	testCases := []struct {
		name        string
		blob        *eip4844.Blob
		proof       eip4844.KZGProof
		commitment  eip4844.KZGCommitment
		expectError bool
	}{
		{
			name:        "Valid Proof",
			blob:        validBlob,
			proof:       validProof,
			commitment:  validCommitment,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifier.VerifyBlobProof(tc.blob, tc.proof, tc.commitment)
			if tc.expectError {
				require.Error(t, err, "cgo is not enabled")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
