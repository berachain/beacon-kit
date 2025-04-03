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

package ssz

import (
	"fmt"

	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/karalabe/ssz"
)

// Unmarshal is the way we build objects from byte formatted in SSZ encoding.
// This function highlights the common template for SSZ decoding different
// objects.
func Unmarshal[T constraints.SSZUnmarshaler](buf []byte, v T) error {
	if err := ssz.DecodeFromBytes(buf, v); err != nil {
		return fmt.Errorf("failed decoding %T: %w", v, err)
	}

	// Note: ValidateAfterDecodingSSZ may change v even if it returns error
	// (depending on the specific implementations)
	return v.ValidateAfterDecodingSSZ()
}
