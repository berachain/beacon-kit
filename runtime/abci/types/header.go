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

import (
	"time"

	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cometbft "github.com/cometbft/cometbft/types"
)

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

	EvidenceHash    [32]byte `ssz-size:"32"`
	ProposerAddress [20]byte `ssz-size:"20"`
}

type BlockID struct {
	Hash          [32]byte `ssz-size:"32"`
	PartSetHeader *PartSetHeader
}

type PartSetHeader struct {
	Hash  [32]byte `ssz-size:"32"`
	Total uint32
}

type Consensus struct {
	Block uint64
	App   uint64
}

func (c *Consensus) ToCometBFT() cmtversion.Consensus {
	return cmtversion.Consensus{
		Block: c.Block,
		App:   c.App,
	}
}

func (c *Consensus) FromCometBFT(consensus cmtversion.Consensus) {
	c.Block = consensus.Block
	c.App = consensus.App
}

func (b *BlockID) ToCometBFT() cometbft.BlockID {
	if b.PartSetHeader == nil {
		b.PartSetHeader = &PartSetHeader{}
	}
	return cometbft.BlockID{
		Hash:          cmtbytes.HexBytes(b.Hash[:]),
		PartSetHeader: b.PartSetHeader.ToCometBFT(),
	}
}

func (b *BlockID) FromCometBFT(blockID cometbft.BlockID) {
	if b.PartSetHeader == nil {
		b.PartSetHeader = &PartSetHeader{}
	}
	b.Hash = byteslib.ToBytes32(blockID.Hash)
	b.PartSetHeader.FromCometBFT(blockID.PartSetHeader)
}

func (p *PartSetHeader) ToCometBFT() cometbft.PartSetHeader {
	return cometbft.PartSetHeader{
		Hash:  cmtbytes.HexBytes(p.Hash[:]),
		Total: p.Total,
	}
}

func (p *PartSetHeader) FromCometBFT(partSetHeader cometbft.PartSetHeader) {
	p.Hash = byteslib.ToBytes32(partSetHeader.Hash)
	p.Total = partSetHeader.Total
}

func (h *CometBFTHeader) ToCometBFT() cometbft.Header {
	if h.Version == nil {
		h.Version = &Consensus{}
	}
	if h.LastBlockID == nil {
		h.LastBlockID = &BlockID{}
	}
	return cometbft.Header{
		Version:            h.Version.ToCometBFT(),
		ChainID:            string(h.ChainID),
		Height:             int64(h.Height),
		Time:               time.Unix(int64(h.Time), 0).UTC(),
		LastBlockID:        h.LastBlockID.ToCometBFT(),
		LastCommitHash:     cmtbytes.HexBytes(h.LastCommitHash[:]),
		DataHash:           cmtbytes.HexBytes(h.DataHash[:]),
		ValidatorsHash:     cmtbytes.HexBytes(h.ValidatorsHash[:]),
		NextValidatorsHash: cmtbytes.HexBytes(h.NextValidatorsHash[:]),
		ConsensusHash:      cmtbytes.HexBytes(h.ConsensusHash[:]),
		AppHash:            cmtbytes.HexBytes(h.AppHash[:]),
		LastResultsHash:    cmtbytes.HexBytes(h.LastResultsHash[:]),
		EvidenceHash:       cmtbytes.HexBytes(h.EvidenceHash[:]),
		ProposerAddress:    h.ProposerAddress[:],
	}
}

func (h *CometBFTHeader) FromCometBFT(header cometbft.Header) {
	if h.Version == nil {
		h.Version = &Consensus{}
	}
	if h.LastBlockID == nil {
		h.LastBlockID = &BlockID{}
	}
	h.Version.FromCometBFT(header.Version)
	h.ChainID = []byte(header.ChainID)
	h.Height = uint64(header.Height)
	h.Time = uint64(header.Time.Unix())
	h.LastBlockID.FromCometBFT(header.LastBlockID)
	h.LastCommitHash = byteslib.ToBytes32(header.LastCommitHash)
	h.DataHash = byteslib.ToBytes32(header.DataHash)
	h.ValidatorsHash = byteslib.ToBytes32(header.ValidatorsHash)
	h.NextValidatorsHash = byteslib.ToBytes32(header.NextValidatorsHash)
	h.ConsensusHash = byteslib.ToBytes32(header.ConsensusHash)
	h.AppHash = byteslib.ToBytes32(header.AppHash)
	h.LastResultsHash = byteslib.ToBytes32(header.LastResultsHash)
	h.EvidenceHash = byteslib.ToBytes32(header.EvidenceHash)
	h.ProposerAddress = byteslib.ToBytes20(header.ProposerAddress)
}
