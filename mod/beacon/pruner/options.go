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

package pruner

import (
	"errors"
	"time"

	"github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/runtime"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/prunedb"
)

// WithBaseService returns an Option that sets the BaseService for the Service.
func WithBaseService(base service.BaseService) service.Option[Service] {
	return func(s *Service) error {
		s.BaseService = base
		return nil
	}
}

// WithWindowSize returns an Option that sets the window size for the Service.
func WithWindowSize(size uint64) service.Option[Service] {
	return func(s *Service) error {
		s.windowSize = size
		return nil
	}
}

// WithStorageBackend returns an Option that sets the storage backend
// for the Service.
// TODO: wtf is going on this is so hood
func WithStorageBackend(bsb runtime.BeaconStorageBackend[
	consensus.ReadOnlyBeaconBlockBody, *datypes.BlobSidecars,
]) service.Option[Service] {
	return func(s *Service) error {
		avs := bsb.AvailabilityStore(nil)
		daStore, ok := avs.(*store.Store[consensus.ReadOnlyBeaconBlockBody])
		if !ok {
			return errors.New("invalid AvailabilityStore type")
		}
		rangeDB, ok := daStore.IndexDB.(*filedb.RangeDB)
		if !ok {
			return errors.New("invalid IndexDB type")
		}

		s.db = prunedb.New(rangeDB)
		return nil
	}
}

// WithInterval returns an Option that sets the pruning interval for the
// Service.
func WithInterval(interval time.Duration) service.Option[Service] {
	return func(s *Service) error {
		s.ticker = time.NewTicker(interval)
		return nil
	}
}
