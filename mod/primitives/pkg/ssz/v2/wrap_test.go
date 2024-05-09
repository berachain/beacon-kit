// SPDX-License-Identifier: MIT
//
// # Copyright (c) 2024 Berachain Foundation
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
//
//nolint:errcheck // do not check for err returns
package ssz_test

import (
	"os"
	"testing"

	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2"
	"github.com/stretchr/testify/require"
)

type Test struct {
	A uint64
}

func getTestStruct() (*sszv2.Checkpoint, error) {
	// https://goerli.beaconcha.in/slot/4744352
	// Test fixture from fastssz.
	const tfn = "fixtures/beacon_state_bellatrix.ssz"
	// A checkpt is the simplest field.
	data, err := os.ReadFile(tfn)
	if err != nil {
		return nil, err
	}
	sszState := sszv2.BeaconStateBellatrix{}
	sszState.UnmarshalSSZ(data)
	if sszState.CurrentJustifiedCheckpoint == nil {
		return nil, err
	}
	return sszState.CurrentJustifiedCheckpoint, nil
}

// Test cases for SSZWrapper.
func TestSSZWrapper_SizeSSZ(t *testing.T) {
	testStruct, err := getTestStruct()
	require.NoError(t, err)

	wrapper := sszv2.Wrap(testStruct)
	size := wrapper.SizeSSZ()
	require.Equal(t, 40, size, "incorrect size")
}

func TestSSZWrapper_MarshalSSZ(t *testing.T) {
	testStruct, err := getTestStruct()
	require.NoError(t, err)

	wrapper := sszv2.Wrap(testStruct)
	data, err := wrapper.MarshalSSZ()
	if err != nil {
		t.Errorf("Failed to marshal SSZ: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Expected non-empty data after marshaling")
	}
}

func TestSSZWrapper_HashTreeRoot(t *testing.T) {
	testStruct, err := getTestStruct()
	require.NoError(t, err)

	wrapper := sszv2.Wrap(testStruct)
	_, err2 := wrapper.HashTreeRoot()
	if err2 != nil {
		t.Errorf("Failed to compute hash tree root: %v", err2)
	}
}
