// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package decoder

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/karalabe/ssz"
)

// SSZUnmarshal is the way we build objects from byte formatted as ssz
// While logically related to constraints package, SSZUnmarshal has its own
// small package to avoid import cycle related to Unused Type
// Also SSZUnmarshal highlight the common template for SSZ decoding different
// objects
func SSZUnmarshal[T SSZUnmarshaler](buf []byte, v T) error {
	switch dest := any(v).(type) {
	case *common.UnusedType:
		// unused types have special formatting for efficiency
		return common.DecodeUnusedType(buf, dest)
	default:
		if err := ssz.DecodeFromBytes(buf, v); err != nil {
			return fmt.Errorf("failed decoding %T: %w", dest, err)
		}
		return v.EnsureSyntaxFromSSZ()
	}
}
