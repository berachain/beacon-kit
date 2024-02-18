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

package engine

import (
	"github.com/itsdevbear/bolaris/types/consensus/version"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// EmptyExecutionPayloadWithVersion returns an empty execution payload for the given version.
func EmptyPayloadAttributesWithVersion(v int) PayloadAttributer {
	switch v {
	case version.Deneb:
		return &enginev1.PayloadAttributesContainer{
			Attributes: &enginev1.PayloadAttributesContainer_V3{},
		}
	case version.Capella:
		return &enginev1.PayloadAttributesContainer{
			Attributes: &enginev1.PayloadAttributesContainer_V2{},
		}
	default:
		return nil
	}
}
