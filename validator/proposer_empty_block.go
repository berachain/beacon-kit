// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package validator

import (
	"fmt"

	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

func (s *Service) getEmptyBlock(slot primitives.Slot) (interfaces.BeaconKitBlock, error) {
	var (
		sBlk interfaces.BeaconKitBlock
		err  error
		fork = s.BeaconCfg().ActiveForkVersion(primitives.Epoch(slot))
	)

	switch fork {
	case version.Deneb:
		// TODO: SUPPORT
		panic("ERROR: deneb fork is not yet supported in beacon-kit.")
	case version.Capella:
		// TODO: generalize the beacon kit block building using a factory pattern.
		sBlk, err = consensusv1.NewBeaconKitBlock(slot, nil, version.Capella)
		if err != nil {
			return nil, fmt.Errorf("could not initialize block for proposal: %w", err)
		}
	default:
		err = fmt.Errorf(
			"unknown fork version %d for slot %d", fork, slot)
	}
	return sBlk, err
}
