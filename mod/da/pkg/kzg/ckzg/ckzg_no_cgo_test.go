//go:build !ckzg

package ckzg_test

import (
	"encoding/json"
	ckzg "github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"os"

	"testing"
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
		}
	}
	var test Test

	err = json.Unmarshal(data, &test)
	require.NoError(t, err)

	errBlob := validBlob.UnmarshalJSON([]byte(`"` + test.Input.Blob + `"`))
	require.NoError(t, errBlob)

	err = validCommitment.UnmarshalJSON([]byte(`"` + test.Input.Commitment + `"`))
	require.NoError(t, errBlob)

	err = validProof.UnmarshalJSON([]byte(`"` + test.Input.Proof + `"`))
	require.NoError(t, errBlob)

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
