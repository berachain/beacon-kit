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

package engine

import (
	eth "github.com/itsdevbear/bolaris/engine/ethclient"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// processPayloadStatusResult processes the payload status result and
// returns the latest valid hash or an error.
func processPayloadStatusResult(
	result *enginev1.PayloadStatus,
) ([]byte, error) {
	switch result.GetStatus() {
	case enginev1.PayloadStatus_INVALID_BLOCK_HASH:
		return nil, eth.ErrInvalidBlockHashPayloadStatus
	case enginev1.PayloadStatus_ACCEPTED, enginev1.PayloadStatus_SYNCING:
		return nil, eth.ErrAcceptedSyncingPayloadStatus
	case enginev1.PayloadStatus_INVALID:
		return result.GetLatestValidHash(), eth.ErrInvalidPayloadStatus
	case enginev1.PayloadStatus_VALID:
		return result.GetLatestValidHash(), nil
	case enginev1.PayloadStatus_UNKNOWN:
		fallthrough
	default:
		return nil, eth.ErrUnknownPayloadStatus
	}
}
