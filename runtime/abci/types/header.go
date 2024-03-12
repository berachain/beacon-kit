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

package types

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs CometBFTHeader,BlockID,PartSetHeader -output generated.ssz.go

// CometBFTHeader is the header of a block in the Comet BFT consensus algorithm.
// Doc: https://docs.cometbft.com/main/spec/core/data_structures#header
type CometBFTHeader struct {
	Version *Consensus
	// ChainID is an unstructured string with a max length of 50-bytes.
	ChainID []byte `ssz-size:"50"`
	Height  uint64
	// Unix time in seconds.
	Time uint64

	LastBlockID *BlockID

	LastCommitHash [32]byte `ssz-size:"32"`
	DataHash       [32]byte `ssz-size:"32"`

	ValidatorsHash     [32]byte `ssz-size:"32"`
	NextValidatorsHash [32]byte `ssz-size:"32"`
	ConsensusHash      [32]byte `ssz-size:"32"`
	AppHash            [32]byte `ssz-size:"32"`
	LastResultsHash    [32]byte `ssz-size:"32"`

	EvidenceHash [32]byte `ssz-size:"32"`
	ProposerAddr [20]byte `ssz-size:"20"`
}

type BlockID struct {
	Hash        [32]byte `ssz-size:"32"`
	PartsHeader *PartSetHeader
}

type PartSetHeader struct {
	Hash  [32]byte `ssz-size:"32"`
	Total uint32
}

type Consensus struct {
	Block uint64
	App   uint64
}
