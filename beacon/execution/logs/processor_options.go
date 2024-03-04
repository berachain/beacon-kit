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

package logs

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/forkchoicer"
	engineclient "github.com/itsdevbear/bolaris/engine/client"
)

// WithForkChoiceProvider returns an Option for
// setting the fork choice provider for Processor.
func WithForkChoicer(fcs forkchoicer.ForkChoicer) Option[Processor] {
	return func(p *Processor) error {
		p.fcs = fcs
		return nil
	}
}

// WithFinalizedLogsStore returns an Option for
// setting the finalized logs provider for Processor.
func WithFinalizedLogsStore(fls FinalizedLogsStore) Option[Processor] {
	return func(p *Processor) error {
		p.fls = fls
		return nil
	}
}

// WithEngineCaller is an option to set the Engine for the Service.
func WithEngineCaller(ec engineclient.Caller) Option[Processor] {
	return func(p *Processor) error {
		p.engine = ec
		return nil
	}
}

// WithLogFactory returns an Option for setting the
// log factory for Processor.
func WithLogFactory(factory LogFactory) Option[Processor] {
	return func(p *Processor) error {
		p.factory = factory
		registeredSigs := factory.GetRegisteredSignatures()
		p.sigToCache = make(map[ethcommon.Hash]LogCache, len(registeredSigs))
		for _, sig := range registeredSigs {
			p.sigToCache[sig] = &Cache{}
		}
		return nil
	}
}
