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

package deneb

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine/interfaces"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// WrappedPayloadAttributesV2 wraps the PayloadAttributesV2 from the Prysmatic Labs' Engine API v1.
var _ interfaces.PayloadAttributer = (*WrappedPayloadAttributesV3)(nil)

// WrappedPayloadAttributesV2 is a struct that embeds enginev1.PayloadAttributesV2
// to provide additional functionality required by the PayloadAttributer interface.
type WrappedPayloadAttributesV3 struct {
	enginev1.PayloadAttributesV3
}

// Version returns the consensus version for the Capella upgrade.
func (w *WrappedPayloadAttributesV3) Version() int {
	return version.Deneb
}
