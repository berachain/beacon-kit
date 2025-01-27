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

package deposit

import (
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/utils/genesis"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/spf13/cobra"
)

// Get the genesis validator root. If the genesis validator root flag is not set, the genesis
// validator root is computed from the genesis file at the last argument (idx: maxArgs - 1).
func getGenesisValidatorRoot(
	cmd *cobra.Command, chainSpec chain.Spec, args []string, maxArgs int,
) (common.Root, error) {
	var genesisValidatorRoot common.Root
	genesisValidatorRootStr, err := cmd.Flags().GetString(useGenesisValidatorRoot)
	if err != nil {
		return common.Root{}, err
	}

	if genesisValidatorRootStr != defaultGenesisValidatorRoot {
		genesisValidatorRoot, err = parser.ConvertGenesisValidatorRoot(genesisValidatorRootStr)
	} else {
		if len(args) != maxArgs {
			return common.Root{}, errors.New(
				"genesis validator root is required if not using the genesis file flag",
			)
		}
		genesisValidatorRoot, err = genesis.ComputeValidatorsRootFromFile(args[maxArgs-1], chainSpec)
	}

	return genesisValidatorRoot, err
}
