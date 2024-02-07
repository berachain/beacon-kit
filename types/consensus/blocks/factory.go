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

package blocks

import (
	eth "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
)

// NewBeaconBlock creates a beacon block from a protobuf beacon block.
func NewBeaconKitBlock(i interface{}) (interfaces.ReadOnlyBeaconKitBlock, error) {
	switch b := i.(type) {
	case nil:
		return nil, blocks.ErrNilObject
	case *eth.BeaconKitBlockCapella:
		_ = b
		return nil, nil
	// case *eth.BeaconBlockDeneb:
	// 	return initSignedBlockFromProtoDeneb(b)
	default:
		return nil, errors.Wrapf(errUnsupportedBeaconBlock, "unable to create block from type %T", i)
	}
}

// NewBeaconKitBlockBody creates a beacon block body from a protobuf beacon block body.
func NewBeaconKitBlockBody(i interface{}) (interfaces.ReadOnlyBeaconKitBlockBody, error) {
	switch b := i.(type) {
	case nil:
		return nil, blocks.ErrNilObject
	case *eth.BeaconKitBlockBodyCapella:
		_ = b
		return nil, nil
	// case *eth.BeaconBlockBodyDeneb:
	// 	return initBlockBodyFromProtoDeneb(b)
	default:
		return nil, errors.Wrapf(
			errUnsupportedBeaconBlockBody, "unable to create block body from type %T", i)
	}
}
