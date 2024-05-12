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

// func TestVerifier_VerifyBlobProof(t *testing.T) {
//	t.Run("should verify proof", func(t *testing.T) {
//		fs := afero.NewOsFs()
//		file, err := afero.ReadFile(
//		fs, "../../../../../testing/files/kzg-trusted-setup.json")
//		if err != nil {
//			fmt.Println("err", err)
//		}
//		require.NoError(t, err)
//
//		// Get the contents from file
//		var ts gokzg4844.JSONTrustedSetup
//		err = json.Unmarshal(file, &ts)
//		if err != nil {
//			require.Error(t, err)
//			return
//		}
//
//		verifier, err := ckzg.NewVerifier(&ts)
//		require.NoError(t, err)
//		require.NotNil(t, verifier)
//
//		// Load the test data
//		file, err = afero.ReadFile(
//		fs, "../../../../../testing/files/kzg-proof.json")
//		if err != nil {
//			fmt.Println("err", err)
//		}
//		require.NoError(t, err)
//
//		// Get the contents from file
//
//		var proofData ckzg4844.Blob
//		err = json.Unmarshal(file, &proofData)
//		if err != nil {
//			require.Error(t, err)
//			return
//		}
//
//		// Verify the proof
//		err = verifier.VerifyBlobProof(
//			(*ckzg4844.Blob)(&proofData.Blob),
//			eip4844.KZGProof((ckzg4844.Bytes48)(proofData.Proof)),
//			eip4844.KZGCommitment((ckzg4844.Bytes48)(proofData.Commitment)),
//		)
//		require.NoError(t, err)
//	})
//}
