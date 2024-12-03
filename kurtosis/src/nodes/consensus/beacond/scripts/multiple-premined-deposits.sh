#!/bin/bash
# SPDX-License-Identifier: BUSL-1.1
#
# Copyright (C) 2024, Berachain Foundation. All rights reserved.
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


/usr/bin/beacond init --chain-id $BEACOND_CHAIN_ID $BEACOND_MONIKER --home /tmp/config0/.beacond --consensus-key-algo $BEACOND_CONSENSUS_KEY_ALGO
/usr/bin/beacond genesis add-premined-deposit --home /tmp/config0/.beacond
cp -r /tmp/config0 /tmp/config_genesis

for ((i=1; i<$NUM_VALS; i++)); do
    BEACOND_HOME=/tmp/config${i}/.beacond
    echo $BEACOND_HOME
    BEACOND_MONIKER=cl-validator-beaconkit-${i}
    /usr/bin/beacond init --chain-id $BEACOND_CHAIN_ID $BEACOND_MONIKER --home $BEACOND_HOME --consensus-key-algo $BEACOND_CONSENSUS_KEY_ALGO
    /usr/bin/beacond genesis add-premined-deposit --home $BEACOND_HOME
    cp -r /tmp/config${i}/.beacond/config/premined-deposits/premined-deposit* /tmp/config_genesis/.beacond/config/premined-deposits/
done

/usr/bin/beacond genesis execution-payload $ETH_GENESIS --home /tmp/config_genesis/.beacond
/usr/bin/beacond genesis collect-premined-deposits --home /tmp/config_genesis/.beacond

