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

package blockchain

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	beacontypes "github.com/itsdevbear/bolaris/beacon/core/types"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/config/version"
	"github.com/itsdevbear/bolaris/primitives"
)

// sendFCU sends a forkchoice update to the execution client.
func (s *Service) sendFCU(
	ctx context.Context,
	headEth1Hash common.Hash,
) error {
	_, err := s.es.NotifyForkchoiceUpdate(
		ctx, &execution.FCUConfig{
			HeadEth1Hash: headEth1Hash,
		})
	return err
}

// sendFCUWithAttributes sends a forkchoice update to the
// execution client with payload attributes. It does
// so via the local builder service.
func (s *Service) sendFCUWithAttributes(
	ctx context.Context,
	headEth1Hash common.Hash,
	slot primitives.Slot,
	parentBlockRoot [32]byte,
) error {
	_, err := s.lb.BuildLocalPayload(
		ctx,
		headEth1Hash,
		slot+1,
		//#nosec:G701 // won't realistically overflow.
		uint64(time.Now().Unix()),
		parentBlockRoot,
	)
	return err
}

// notifyNewPayload notifies the execution client of a new payload.
func (s *Service) notifyNewPayload(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBuoy,
) (bool, error) {
	if err := beacontypes.BeaconBuoyIsNil(blk); err != nil {
		return false, err
	}

	executionData, err := blk.ExecutionPayload()
	if err != nil {
		return false, err
	}

	if executionData == nil || executionData.IsEmpty() {
		return false, errors.New("no payload in beacon block")
	}

	var versionedHashes []common.Hash
	if blk.Version() >= version.Deneb {
		// TODO: EIP-4844
		versionedHashes = make([]common.Hash, 0)
		// var versionedHashes []common.Hash
		// versionedHashes, err =
		// kzgCommitmentsToVersionedHashes(blk.Block().Body())
		// if err != nil {
		// 	return false, errors.Wrap(err, "could not get versioned hashes to
		// feed the engine")
		// }
	}

	return s.es.NotifyNewPayload(
		ctx,
		blk.GetSlot(),
		executionData,
		versionedHashes,
		common.Hash(blk.GetParentRoot()),
	)
}
