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
	"testing"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewVerifier(t *testing.T) {
	t.Run("should create verifier", func(t *testing.T) {
		fs := afero.NewOsFs()
		file, err := afero.ReadFile(
			fs, "../../../../../testing/files/kzg-trusted-setup.json")

		require.NoError(t, err)

		// Get the contents from file
		var ts gokzg4844.JSONTrustedSetup
		err = json.Unmarshal(file, &ts)
		if err != nil {
			require.Error(t, err)
			return
		}

		verifier, err := ckzg.NewVerifier(&ts)
		require.NoError(t, err)
		require.NotNil(t, verifier)
	})
}
