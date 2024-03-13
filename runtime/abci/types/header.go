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
	"unsafe"

	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cometbft "github.com/cometbft/cometbft/types"
)

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs CometBFTHeader,BlockID,PartSetHeader -output generated.ssz.go

// CometBFTHeader is an ssz'able struct,
// equivalent to the Header struct in CometBFT.
// Doc: https://docs.cometbft.com/main/spec/core/data_structures#header
type CometBFTHeader struct {
	Version *Consensus
	// ChainID is an unstructured string with a max length of 50-bytes.
	ChainID []byte `ssz-max:"50"`
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

// BlockID is an ssz'able struct,
// equivalent to the BlockID struct in CometBFT.
type BlockID struct {
	Hash          [32]byte `ssz-size:"32"`
	PartSetHeader *PartSetHeader
}

// PartSetHeader is an ssz'able struct,
// equivalent to the PartSetHeader struct in CometBFT.
type PartSetHeader struct {
	Hash  [32]byte `ssz-size:"32"`
	Total uint32
}

// Consensus is an ssz'able struct,
// equivalent to the Consensus struct in CometBFT.
type Consensus struct {
	Block uint64
	App   uint64
}

// ToCometBFT converts CometBFTHeader to its CometBFT equivalent.
func (h *CometBFTHeader) ToCometBFT() cometbft.Header {
	return cometbft.Header{
		Version: *(*cmtversion.Consensus)(unsafe.Pointer(h.Version)),
		ChainID: string(h.ChainID),
		//#nosec:G701 // int64 is sufficient as block time is greater than a second.
		Height: int64(h.Height),
		//#nosec:G701 // int64 is sufficient for billions of years.
		Time:               time.Unix(int64(h.Time), 0).UTC(),
		LastBlockID:        *(*cometbft.BlockID)(unsafe.Pointer(h.LastBlockID)),
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

// FromCometBFT converts a CometBFT Header to its equivalent.
func (h *CometBFTHeader) FromCometBFT(header cometbft.Header) {
	h.Version = (*Consensus)(unsafe.Pointer(&header.Version))
	h.ChainID = []byte(header.ChainID)
	//#nosec:G701 // A positive int64 can never overflow a uint64.
	h.Height = uint64(header.Height)
	//#nosec:G701 // A positive int64 can never overflow a uint64.
	h.Time = uint64(header.Time.Unix())
	h.LastBlockID = (*BlockID)(unsafe.Pointer(&header.LastBlockID))
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
