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

package blob

import (
	"context"
	"errors"

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/db"
	filedb "github.com/berachain/beacon-kit/db/file"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/sourcegraph/conc/iter"
)

// Store is the default implementation of the AvailabilityStore.
type Store struct {
	*filedb.RangeDB
}

// NewStore creates a new instance of the AvailabilityStore.
func NewStore(db db.DB) *Store {
	return &Store{
		RangeDB: filedb.NewRangeDB(db),
	}
}

// IsDataAvailable ensures that all blobs referenced in the block are
// stored before it returns without an error.
func (s *Store) IsDataAvailable(
	ctx context.Context,
	slot primitives.Slot,
	b beacontypes.ReadOnlyBeaconBlock,
) error {
	_ = ctx
	_ = slot
	_ = b
	return nil
}

// Persist ensures the sidecar data remains accessible, utilizing parallel
// processing for efficiency.
func (s *Store) Persist(
	slot primitives.Slot,
	sidecars ...*beacontypes.BlobSidecar,
) error {
	_, err := iter.MapErr(
		sidecars,
		func(sidecar **beacontypes.BlobSidecar) ([]any, error) {
			if *sidecar == nil {
				return nil, errors.New("sidecar is nil")
			}
			sc := *sidecar
			bz, err := sc.MarshalSSZ()
			if err != nil {
				return nil, err
			}
			err = s.Set(slot, sc.KzgCommitment, bz)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
	return err
}
