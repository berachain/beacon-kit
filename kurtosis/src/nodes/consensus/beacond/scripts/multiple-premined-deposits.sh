#!/usr/bin/env bash

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

