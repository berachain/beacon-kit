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

	"github.com/ethereum/go-ethereum/common"
	types "github.com/itsdevbear/bolaris/execution/builder/types"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus"
	consensusinterfaces "github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	engineinterfaces "github.com/itsdevbear/bolaris/types/engine/interfaces"
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
) (consensusinterfaces.ReadOnlyBeaconKitBlock, error) {
	// Process all builders, both local and remote
	localBuilders := s.builders.LocalBuilders()
	// TODO: handle remote builders.
	// remoteBuilders := s.builders.RemoteBuilders()

	type concResponse struct {
		builderName string
		response    *types.GetExecutionPayloadResponse
		err         error
	}

	eth1ParentHash, err := s.getParentEth1Hash(ctx)
	if err != nil {
		return nil, err
	}

	processBuilders := func(
		builders []*types.BuilderEntry,
	) ([]engineinterfaces.ExecutionPayload, error) {
		resps := iter.Map[*types.BuilderEntry, *concResponse](
			builders,
			func(b **types.BuilderEntry) *concResponse {
				builder := *b
				dctx, cancel := context.WithTimeout(ctx, DefaultBuilderTimeout)
				defer cancel()
				var resp *types.GetExecutionPayloadResponse
				resp, err = builder.BuilderServiceClient.GetExecutionPayload(
					dctx, &types.GetExecutionPayloadRequest{
						ParentHash: eth1ParentHash.Hex(),
						Slot:       slot,
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
		var out []engineinterfaces.ExecutionPayload
		for _, resp := range resps {
			if resp.err != nil || resp.response == nil {
				continue
			}
			out = append(out, resp.response.GetPayloadContainer())
		}

		// If no builders returned a block, return an error
		if len(out) == 0 {
			return nil, errors.New("all builders failed to return a block")
		}

		// Return the first block for now
		// TODO: actually sort thru and create an algorithmn to choose the best block.
		return out, nil
	}

	payloads, err := processBuilders(localBuilders)
	if err != nil {
		return nil, err
	}

	// Return the first block for now
	// TODO: actually sort thru and create an algo to choose the best block.
	payload := payloads[0]

	// // TODO: SIGN UR RANDAO THINGY HERE OR SOMETHING.
	// _ = k.beaconKitValKey
	// _, err := s.beaconKitValKey.Key.PrivKey.Sign([]byte("hello world"))
	// if err != nil {
	// 	return nil, err
	// }

	// Create a block from the payload
	return consensus.BeaconKitBlockFromState(
		s.BeaconState(ctx),
		payload,
	)
}

// getParentEth1Hash retrieves the parent block hash for the given slot.
//
//nolint:unparam // todo: review this later.
func (s *Service) getParentEth1Hash(ctx context.Context) (common.Hash, error) {
	// The first slot should be proposed with the genesis block as parent.
	st := s.BeaconState(ctx)
	if st.Slot() == 1 {
		return st.GenesisEth1Hash(), nil
	}

	// We always want the parent block to be the last finalized block.
	return st.GetFinalizedEth1BlockHash(), nil
}
