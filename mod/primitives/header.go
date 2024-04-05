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

package primitives

import (
	"encoding/json"
)

// BeaconBlockHeader is the header of a beacon block.
//
//go:generate go run github.com/fjl/gencodec -type BeaconBlockHeader -out header.json.go
type BeaconBlockHeader struct {
	// Slot is the slot number of the block.
	Slot Slot `json:"slot"`
	// ProposerIndex is the index of the proposer of the block.
	ProposerIndex ValidatorIndex `json:"proposerIndex"`
	// ParentRoot is the root of the parent block.
	ParentRoot Root `json:"parentRoot"    ssz-size:"32"`
	// StateRoot is the root of the beacon state after executing
	// the block. Will be 0x00...00 prior to execution.
	StateRoot Root `json:"stateRoot"     ssz-size:"32"`
	// 	// BodyRoot is the root of the block body.
	BodyRoot Root `json:"bodyRoot"      ssz-size:"32"`
}

// String returns a string representation of the beacon block header.
func (h *BeaconBlockHeader) String() string {
	//#nosec:G703 // ignore potential marshalling failure.
	output, _ := json.Marshal(h)
	return string(output)
}
