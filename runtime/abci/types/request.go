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

	cometabci "github.com/cometbft/cometbft/abci/types"
)

// WrappedABCIRequest provides methods to attach and detach BeaconKitBlock.
// It also provides a method to get the height of the block.
type WrappedABCIRequest interface {
	AttachBeaconKitBlock(tx []byte)
	DettachBeaconKitBlock() []byte
	GetHeight() int64
	GetTime() time.Time
	GetProposerAddress() []byte
}

var _ WrappedABCIRequest = (*WrappedRequestFinalizeBlock)(nil)

// WrappedRequestFinalizeBlock wraps the RequestFinalizeBlock struct from the cometabci package.
type WrappedRequestFinalizeBlock struct {
	*cometabci.RequestFinalizeBlock
}

// FromRequestFinalizeBlock creates a new WrappedRequestFinalizeBlock
// from a given RequestFinalizeBlock.
func FromRequestFinalizeBlock(
	req *cometabci.RequestFinalizeBlock,
) *WrappedRequestFinalizeBlock {
	return &WrappedRequestFinalizeBlock{RequestFinalizeBlock: req}
}

// AttachBeaconKitBlock attaches a BeaconKitBlock to the WrappedRequestFinalizeBlock.
func (w *WrappedRequestFinalizeBlock) AttachBeaconKitBlock(tx []byte) {
	w.RequestFinalizeBlock.Txs = append([][]byte{tx}, w.RequestFinalizeBlock.Txs...)
}

// DettachBeaconKitBlock detaches a BeaconKitBlock from the WrappedRequestFinalizeBlock.
func (w *WrappedRequestFinalizeBlock) DettachBeaconKitBlock() []byte {
	tx := w.RequestFinalizeBlock.Txs[0]
	w.RequestFinalizeBlock.Txs = w.RequestFinalizeBlock.Txs[1:]
	return tx
}

// Txs returns the transactions in the WrappedRequestFinalizeBlock.
func (w *WrappedRequestFinalizeBlock) Txs() [][]byte {
	return w.RequestFinalizeBlock.Txs
}

// WrappedResponsePrepareProposal wraps the ResponsePrepareProposal struct
// from the cometabci package.
type WrappedResponsePrepareProposal struct {
	req *cometabci.ResponsePrepareProposal
}

// FromResponsePrepareProposal creates a new WrappedResponsePrepareProposal
// from a given ResponsePrepareProposal.
func FromResponsePrepareProposal(
	req *cometabci.ResponsePrepareProposal,
) *WrappedResponsePrepareProposal {
	return &WrappedResponsePrepareProposal{req: req}
}

// AttachBeaconKitBlock attaches a BeaconKitBlock to the WrappedResponsePrepareProposal.
func (w *WrappedResponsePrepareProposal) AttachBeaconKitBlock(tx []byte) {
	w.req.Txs = append([][]byte{tx}, w.req.Txs...)
}

// DettachBeaconKitBlock detaches a BeaconKitBlock from the WrappedResponsePrepareProposal.
func (w *WrappedResponsePrepareProposal) DettachBeaconKitBlock() []byte {
	tx := w.req.Txs[0]
	w.req.Txs = w.req.Txs[1:]
	return tx
}

// WrappedRequestProcessProposal wraps the RequestProcessProposal struct
// from the cometabci package.
type WrappedRequestProcessProposal struct {
	req *cometabci.RequestProcessProposal
}

// Height returns the height of the block in the WrappedRequestProcessProposal.
func (w *WrappedRequestProcessProposal) Height() int64 {
	return w.req.Height
}

// FromRequestProcessProposal creates a new WrappedRequestProcessProposal
// from a given RequestProcessProposal.
func FromRequestProcessProposal(
	req *cometabci.RequestProcessProposal,
) *WrappedRequestProcessProposal {
	return &WrappedRequestProcessProposal{req: req}
}

// DettachBeaconKitBlock detaches a BeaconKitBlock from the WrappedRequestProcessProposal.
func (w *WrappedRequestProcessProposal) DettachBeaconKitBlock() []byte {
	tx := w.req.Txs[0]
	w.req.Txs = w.req.Txs[1:]
	return tx
}
