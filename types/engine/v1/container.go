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

package enginev1

import (
	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	"github.com/itsdevbear/bolaris/types/engine/v1/capella"
	"github.com/itsdevbear/bolaris/types/engine/v1/deneb"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// BlockContainerFromBlock converts a ReadOnlyBeaconKitBlock interface
// into a BeaconKitBlockContainer pointer.
func PayloadContainerFromPayload(
	payload interfaces.ExecutionPayload,
) *ExecutionPayloadContainer {
	container := &ExecutionPayloadContainer{}
	switch payload := payload.ToProto().(type) {
	case *enginev1.ExecutionPayloadCapella:
		container.Payload = &ExecutionPayloadContainer_Capella{
			Capella: payload,
		}
	case *enginev1.ExecutionPayloadHeaderCapella:
		container.Payload = &ExecutionPayloadContainer_CapellaHeader{
			CapellaHeader: payload,
		}
	case *enginev1.ExecutionPayloadDeneb:
		container.Payload = &ExecutionPayloadContainer_Deneb{
			Deneb: payload,
		}
	case *enginev1.ExecutionPayloadHeaderDeneb:
		container.Payload = &ExecutionPayloadContainer_DenebHeader{
			DenebHeader: payload,
		}
	}
	return container
}

// BlockFromContainer extracts a ReadOnlyBeaconKitBlock interface
// from a given BeaconKitBlockContainer.
func PayloadFromContainer(
	container *ExecutionPayloadContainer,
	value math.Wei,
) interfaces.ExecutionPayload {
	// TODO: support headers
	switch payload := container.GetPayload().(type) {
	case *ExecutionPayloadContainer_Capella:
		return capella.NewWrappedExecutionPayloadCapella(payload.Capella, value)
	// case *ExecutionPayloadContainer_CapellaHeader:
	// 	// return capella.NewWrappedExecutionPayloadCapella(payload.CapellaHeader, nil)
	case *ExecutionPayloadContainer_Deneb:
		return deneb.NewWrappedExecutionPayloadDeneb(payload.Deneb, value)
	// case *ExecutionPayloadContainer_DenebHeader:
	// 	return deneb.NewWrappedExecutionPayloadDeneb(payload.Capella, nil)
	// case *BeaconKitBlockContainer_CapellaBlinded:
	// 	return block.CapellaBlinded
	default:
		return nil
	}
}
