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

	"github.com/itsdevbear/bolaris/types/consensus"
)

// ABCIRequest is the interface for an ABCI request.
type ABCIRequest interface {
	GetHeight() int64
	GetTime() time.Time
	GetTxs() [][]byte
}

// ReadOnlyBeaconKitBlockFromABCIRequest assembles a
// new read-only beacon block by extracting a marshalled
// block out of an ABCI request.
func ReadOnlyBeaconKitBlockFromABCIRequest(
	req ABCIRequest,
	bzIndex uint,
	forkVersion int,
) (consensus.ReadOnlyBeaconKitBlock, error) {
	txs := req.GetTxs()

	// Ensure there are transactions in the request and
	// that the request is valid.
	if lenTxs := uint(len(txs)); lenTxs == 0 {
		return nil, ErrNoBeaconBlockInRequest
	} else if bzIndex >= lenTxs {
		return nil, ErrBzIndexOutOfBounds
	}

	// Extract the beacon block from the ABCI request.
	return consensus.BeaconKitBlockFromSSZ(txs[bzIndex], forkVersion)
}
