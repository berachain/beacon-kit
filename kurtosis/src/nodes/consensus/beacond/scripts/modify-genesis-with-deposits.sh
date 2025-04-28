#!/bin/bash
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
# AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
# EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
# TITLE.


# Sets the deposit storage in the the new eth-genesis file in the home directory.
/usr/bin/beacond genesis set-deposit-storage $ETH_GENESIS --beacon-kit.chain-spec $CHAIN_SPEC --home /tmp/config_genesis/.beacond

# Get values directly from the storage fields
DEPOSIT_COUNT=$(jq -r '.alloc["0x4242424242424242424242424242424242424242"].storage["0x0000000000000000000000000000000000000000000000000000000000000000"]' /tmp/config_genesis/.beacond/genesis.json)
DEPOSIT_ROOT=$(jq -r '.alloc["0x4242424242424242424242424242424242424242"].storage["0x0000000000000000000000000000000000000000000000000000000000000001"]' /tmp/config_genesis/.beacond/genesis.json)

/usr/bin/beacond genesis execution-payload /tmp/config_genesis/.beacond/genesis.json --beacon-kit.chain-spec $CHAIN_SPEC --home /tmp/config_genesis/.beacond

# Write each value to separate files for easier parsing
mkdir -p /tmp/values
printf "%s" "$DEPOSIT_COUNT" > /tmp/values/deposit_count.txt
printf "%s" "$DEPOSIT_ROOT" > /tmp/values/deposit_root.txt
