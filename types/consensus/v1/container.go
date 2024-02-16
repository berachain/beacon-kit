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

package consensusv1

import "github.com/itsdevbear/bolaris/types/consensus/interfaces"

// BlockContainerFromBlock converts a ReadOnlyBeaconKitBlock interface
// into a BeaconKitBlockContainer pointer.
func BlockContainerFromBlock(
	block interfaces.ReadOnlyBeaconKitBlock,
) *BeaconKitBlockContainer {
	container := &BeaconKitBlockContainer{}
	switch block := block.(type) {
	case *BeaconKitBlockCapella:
		container.Block = &BeaconKitBlockContainer_Capella{
			Capella: block,
		}
	// case *BlindedBeaconKitBlockCapella:
	// 	container.Block = &BeaconKitBlockContainer_CapellaBlinded{
	// 		CapellaBlinded: block,
	// 	}
	default:
		return nil
	}
	return container
}

// BlockFromContainer extracts a ReadOnlyBeaconKitBlock interface
// from a given BeaconKitBlockContainer.
func BlockFromContainer(
	container *BeaconKitBlockContainer,
) interfaces.ReadOnlyBeaconKitBlock {
	switch block := container.GetBlock().(type) {
	case *BeaconKitBlockContainer_Capella:
		return block.Capella
	// case *BeaconKitBlockContainer_CapellaBlinded:
	// 	return block.CapellaBlinded
	default:
		return nil
	}
}
