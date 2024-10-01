// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package server

import (
	"fmt"
	"strings"

	pruningtypes "cosmossdk.io/store/pruning/types"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/server/types"
	"github.com/spf13/cast"
)

// GetPruningOptionsFromFlags parses command flags and returns the correct
// PruningOptions. If a pruning strategy is provided, that will be parsed and
// returned, otherwise, it is assumed custom pruning options are provided.
func GetPruningOptionsFromFlags(
	appOpts types.AppOptions,
) (pruningtypes.PruningOptions, error) {
	strategy := strings.ToLower(cast.ToString(appOpts.Get(FlagPruning)))

	switch strategy {
	case pruningtypes.PruningOptionDefault,
		pruningtypes.PruningOptionNothing,
		pruningtypes.PruningOptionEverything:
		return pruningtypes.NewPruningOptionsFromString(strategy), nil

	case pruningtypes.PruningOptionCustom:
		opts := pruningtypes.NewCustomPruningOptions(
			cast.ToUint64(appOpts.Get(FlagPruningKeepRecent)),
			cast.ToUint64(appOpts.Get(FlagPruningInterval)),
		)

		if err := opts.Validate(); err != nil {
			return opts, fmt.Errorf("invalid custom pruning options: %w", err)
		}

		return opts, nil

	default:
		return pruningtypes.PruningOptions{}, fmt.Errorf(
			"unknown pruning strategy %s",
			strategy,
		)
	}
}
