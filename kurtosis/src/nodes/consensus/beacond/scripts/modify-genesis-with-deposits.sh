#!/usr/bin/env bash
# SPDX-License-Identifier: BUSL-1.1
#
# Copyright (C) 2025, Berachain Foundation. All rights reserved.
# Use of this software is governed by the Business Source License included
# in the LICENSE file of this repository and at www.mariadb.com/bsl11.
#
# ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
# TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
# VERSIONS OF THE LICENSED WORK.
#
# THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
# LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
# LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
#
# TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
# AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
# EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
# TITLE.


# Sets the deposit storage in the new eth-genesis file in the home directory.
# This creates genesis.json from $ETH_GENESIS
/usr/bin/beacond genesis set-deposit-storage $ETH_GENESIS --beacon-kit.chain-spec $CHAIN_SPEC --home /tmp/config_genesis/.beacond

# The output file is "genesis.json"
ETH_GENESIS_OUTPUT="/tmp/config_genesis/.beacond/genesis.json"

# Generate the execution payload - this populates the deposit contract storage with POL operator keys
/usr/bin/beacond genesis execution-payload $ETH_GENESIS_OUTPUT --beacon-kit.chain-spec $CHAIN_SPEC --home /tmp/config_genesis/.beacond
