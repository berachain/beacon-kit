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

	types "github.com/itsdevbear/bolaris/builder/types"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus/interfaces"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

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
	resp, err := s.builders.GetBuilder(types.DefaultLocalBuilderName).RequestBestBlock(
		ctx, &types.RequestBestBlockRequest{
			Slot: slot,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.GetBlock(), nil
}
