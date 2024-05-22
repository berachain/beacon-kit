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
package deneb_test

import (
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	sszv2 "github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/v2/lib"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // test debug default err msg
var defaultErrMsg = "local output & fastssz output doesnt match"

// Test using local deneb genesis beaconstate.
func TestParityDenebLocal(t *testing.T) {
	// Demonstrate the block is valid by proving
	// the block can be serialized
	// and deserialized back to the same object using fastssz
	block := deneb.BeaconState{}
	genesis, getBlockErr := deneb.DefaultBeaconState()
	require.NoError(t, getBlockErr)
	data, fastSSZMarshalErr := genesis.MarshalSSZ()
	require.NoError(t, fastSSZMarshalErr)
	if data == nil {
		panic("Data is nil")
	}

	if err := block.UnmarshalSSZ(data); err != nil {
		panic(err)
	}

	if block.SizeSSZ() == 0 {
		panic("Block is nil")
	}

	destBlockBz, err := block.MarshalSSZ()
	if err != nil {
		panic(`Step 1: Deserialize-Serialize 
		-- could not serialize back the 
		deserialized input block`)
	}

	if !reflect.DeepEqual(data, destBlockBz) {
		panic(`Step 2: Deserialize-Serialize 
		-- input != serialize(deserialize(input))`)
	}

	// Use our native serializer to do the same
	s := sszv2.NewSerializer()
	o2, err3 := s.MarshalSSZ(genesis)
	require.NoError(t, err3)
	require.Equal(t, len(o2), len(data), defaultErrMsg)

	// TODO: not a full match yet
	// require.Equal(t, o2, data, "local output and fastssz output doesnt
	// match")
}
