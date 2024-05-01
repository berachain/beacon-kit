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
package ssz_test

import (
	"fmt"
	"os"
	"testing"

	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

const TestFileName = "fixtures/beacon_state_bellatrix.ssz" // https://goerli.beaconcha.in/slot/4744352
var debug = false

func debugPrint(debug bool, s ...any) {
	if debug {
		fmt.Println(s...)
	}
}
func TestParityUint64(t *testing.T) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(t, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(t, err)

	object := sszState.LatestBlockHeader
	slot := object.Slot

	s := sszv2.NewSerializer()
	o2, err := s.MarshalSSZ(sszState.LatestBlockHeader.Slot)
	debugPrint(debug, "Local Serializer output:", o2, err)

	res := make([]byte, 0)
	res = ssz.MarshalUint64(res, slot)
	debugPrint(debug, "FastSSZ Output:", res)
	require.Equal(t, o2, res, "local output and fastssz output doesnt match")
}

func BenchmarkNativeUint64(b *testing.B) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	s := sszv2.NewSerializer()
	for i := 0; i < b.N; i++ {
		// Native impl
		o2, err := s.MarshalSSZ(sszState.LatestBlockHeader.Slot)
		debugPrint(false, "Local Serializer output:", o2, err)
	}
}

func BenchmarkFastSSZUint64(b *testing.B) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		res := make([]byte, 0)
		res = ssz.MarshalUint64(res, sszState.LatestBlockHeader.Slot)
		debugPrint(false, "FastSSZ Output:", res)
	}
}

func TestParityByteArray(t *testing.T) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(t, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(t, err)

	s := sszv2.NewSerializer()
	exp, err := s.MarshalSSZ(sszState.LatestBlockHeader.ParentRoot)
	debugPrint(debug, "Local Serializer output:", exp, err)

	res := make([]byte, 0)
	res, err = sszState.LatestBlockHeader.MarshalSSZ()
	prInRes := res[16:48]

	debugPrint(debug, "FastSSZ Output:", prInRes)
	require.Equal(t, exp, prInRes, "local output and fastssz output doesnt match")
}

func BenchmarkNativeByteArray(b *testing.B) {
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	s := sszv2.NewSerializer()
	for i := 0; i < b.N; i++ {
		// Native impl
		exp, err := s.MarshalSSZ(sszState.LatestBlockHeader.ParentRoot)
		debugPrint(debug, "Local Serializer output:", exp, err)
	}
}

func BenchmarkFastSSZByteArray(b *testing.B) {
	debug = false
	data, err := os.ReadFile(TestFileName)
	require.NoError(b, err)

	sszState := sszv2.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(data)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		res := make([]byte, 0)
		res, err = sszState.LatestBlockHeader.MarshalSSZ()
		prInRes := res[16:48]
		debugPrint(debug, "FastSSZ Output:", prInRes)
	}
}

func TestParityByteArrayLarge(t *testing.T) {
	debug = true
	data, err := os.ReadFile(TestFileName)
	require.NoError(t, err)

	sszState := sszv2.BeaconStateBellatrix{}

	err = sszState.UnmarshalSSZ(data)
	require.NoError(t, err)

	s := sszv2.NewSerializer()
	exp, err := s.MarshalSSZ(sszState.StateRoots)
	debugPrint(true, "Local Serializer output:", exp, err)

	// res := make([]byte, 0)
	_, err = sszState.MarshalSSZ()
	prInRes := sszState.StateRoots
	// prInRes := res[16:48]

	// debugPrint(true, "FastSSZ Output:", prInRes, res)
	require.Equal(t, exp, prInRes, "local output and fastssz output doesnt match")
}

func BenchmarkNativeByteArrayLarge(b *testing.B) {
	// 8192 size vs 32
	// size := 8192
	// arr := make([]byte, size)
	// for i := range size {
	// 	arr[i] = byte(i)
	// }
	arr := sszv2.BeaconStateBellatrix{}.StateRoots

	s := sszv2.NewSerializer()

	for i := 0; i < b.N; i++ {
		// Native impl
		exp, err := s.MarshalSSZ(arr)
		debugPrint(debug, "Local Serializer output:", exp, err)
	}
}

func BenchmarkFastSSZByteArrayLarge(b *testing.B) {

	sszState := sszv2.BeaconStateBellatrix{}
	for i := 0; i < b.N; i++ {
		res := make([]byte, 0)
		res, _ = sszState.MarshalSSZ()
		prInRes := res[262320:524464]
		debugPrint(debug, "FastSSZ Output:", prInRes)
	}
}
