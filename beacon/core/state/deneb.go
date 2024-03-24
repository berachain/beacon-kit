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

package state

import (
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultBeaconStateDeneb returns a default BeaconStateDeneb.
func DefaultBeaconStateDeneb() *BeaconStateDeneb {
	return &BeaconStateDeneb{
		Eth1GenesisHash: common.HexToHash(
			"0xa63c365d92faa4de2a64a80ed4759c3e9dfa939065c10af08d2d8d017a29f5f4",
		),
		RandaoMix: make([]byte, types.MixLength),
	}
}

// TODO setup this properly.
//
//go:generate go run github.com/fjl/gencodec -type BeaconStateDeneb -field-override beaconStateDenebJSONMarshaling -out deneb.json.go
type BeaconStateDeneb struct {
	Eth1GenesisHash primitives.ExecutionHash `json:"eth1GenesisHash" ssz-size:"32"`
	RandaoMix       []byte                   `json:"randaoMix"       ssz-size:"32"`
}

// beaconStateDenebJSONMarshaling is a type used to marshal/unmarshal
// BeaconStateDeneb.
type beaconStateDenebJSONMarshaling struct {
	RandaoMix hexutil.Bytes
}
