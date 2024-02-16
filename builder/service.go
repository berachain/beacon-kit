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

package builder

import (
	"context"
	"errors"
	"time"

	types "github.com/itsdevbear/bolaris/builder/types"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/sourcegraph/conc/iter"
)

const DefaultBuilderTimeout = 2 * time.Second

type Service struct {
	service.BaseService
	builders *types.BuilderRegistry
}

// NewService.
func NewService(
	base service.BaseService,
	opts ...Option,
) *Service {
	s := &Service{
		BaseService: base,
		builders:    types.NewBuilderRegistry(),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			panic(err)
		}
	}
	return s
}

func (s *Service) Start(context.Context) {}
func (s *Service) Status() error         { return nil }

// RequestBestBlock requests the best availible block from the builder.
func (s *Service) RequestBestBlock(
	ctx context.Context,
	slot primitives.Slot,
	// version int,
	// TODO: determine if we want this field, or should it be up to the builder
	// to determine the parent.
	/*eth1Parent common.Hash, */
) (interfaces.ReadOnlyBeaconKitBlock, error) {
	// Process all builders, both local and remote
	localBuilders := s.builders.LocalBuilders()
	// TODO: handle remote builders.
	// remoteBuilders := s.builders.RemoteBuilders()

	type concResponse struct {
		builderName string
		response    *types.RequestBestBlockResponse
		err         error
	}

	processBuilders := func(
		builders []*types.BuilderEntry,
	) ([]interfaces.ReadOnlyBeaconKitBlock, error) {
		resps := iter.Map[*types.BuilderEntry, *concResponse](
			builders,
			func(b **types.BuilderEntry) *concResponse {
				builder := *b
				dctx, cancel := context.WithTimeout(ctx, DefaultBuilderTimeout)
				defer cancel()
				resp, err := builder.BuilderServiceClient.RequestBestBlock(
					dctx, &types.RequestBestBlockRequest{
						Slot: slot,
					})
				if err != nil {
					s.Logger().Warn("Failed to request block from builder", "builder", builder.Name, "error", err)
				}
				return &concResponse{
					builderName: builder.Name,
					response:    resp,
					err:         err,
				}
			},
		)

		// Filter out any blocks with bad errors.
		var out []interfaces.ReadOnlyBeaconKitBlock
		for _, resp := range resps {
			if resp.err != nil || resp.response == nil {
				continue
			}
			out = append(out, resp.response.GetBlock())
		}

		// If no builders returned a block, return an error
		if len(out) == 0 {
			return nil, errors.New("all builders failed to return a block")
		}

		// Return the first block for now
		// TODO: actually sort thru and create an algorithmn to choose the best block.
		return out, nil
	}

	blocks, err := processBuilders(localBuilders)
	if err != nil {
		return nil, err
	}

	// TOOD: Process remote builders
	// processBuilders(remoteBuilders)

	// Return the first block for now
	// TODO: actually sort thru and create an algo to choose the best block.
	return blocks[0], nil
}
