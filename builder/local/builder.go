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

package local

import (
	"context"

	"github.com/itsdevbear/bolaris/builder/local/key"
	"github.com/itsdevbear/bolaris/cache"
	"github.com/itsdevbear/bolaris/execution/engine"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// The local Builder is responsible for building beacon blocks and
// handling the execution of beacon chain operations.
type Builder struct {
	service.BaseService

	// TODO: the builder should be decoupled from anything validator related.
	beaconKitValKey key.BeaconKitValKey
	en              engine.Caller
	payloadCache    *cache.PayloadIDCache
}

func NewBuilder(
	base service.BaseService,
	opts ...Option,
) *Builder {
	b := &Builder{
		BaseService: base,
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			panic(err)
		}
	}
	return b
}

func (b *Builder) Start(context.Context) {
}

func (b *Builder) Status() error {
	return nil
}
